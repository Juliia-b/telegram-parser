package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Arman92/go-tdlib"
	_ "github.com/lib/pq"
	"os"
	"time"
)

//    --------------------------------------------------------------------------------
//                                    STRUCTS
//    --------------------------------------------------------------------------------

type PostgresClient struct {
	Connection  *sql.DB
	DbInfo      *DbInfo
	SchemaInfo  *SchemaInfo
	TimePeriods *TimePeriods
}

type DbInfo struct {
	DbName    string
	TableName string
}

type SchemaInfo struct {
	MessageID string
	ChatID    string
	ChatTitle string
	Content   string
	Date      string
	Views     string
	Forwards  string
	Replies   string
}

type TimePeriods struct {
	Today              string
	Yesterday          string
	DayBeforeYesterday string
	ThisWeek           string
	LastWeek           string
	ThisMonth          string
}

type Message struct {
	MessageID int64
	ChatID    int64
	ChatTitle string
	Content   string
	Date      int32
	Views     int32
	Forwards  int32
	Replies   int32
}

type UpdateRow struct {
	MessageId    int64
	ChatId       int64
	NewChatTitle string
	NewContent   string
	NewViews     int32
	NewForwards  int32
	NewReplies   int32
}

//    --------------------------------------------------------------------------------
//                                     METHODS
//    --------------------------------------------------------------------------------

// ConnectToPostgres opens a connection to PostgreSQL
func ConnectToPostgres() (*PostgresClient, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", "localhost", 5432, "postgres", os.Getenv("POSTGRESPASSWORD"), "test")

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	res, err := db.Exec(`CREATE TABLE IF NOT EXISTS tg_parser ( message_id bigint, chat_id bigint, chat_title text, content text, date bigint, views integer, forwards integer, replies integer, PRIMARY KEY(message_id, chat_id) );`)
	if err != nil {
		return nil, err
	}

	// count can be 0 if table already EXISTS, else count 1
	if count, _ := res.RowsAffected(); count > 1 {
		return nil, errors.New("table creation error")
	}

	return &PostgresClient{Connection: db, DbInfo: &DbInfo{"postgres", "tg_parser"}, SchemaInfo: getSchemaInfo(), TimePeriods: getTimePeriods()}, nil
}

// Close closes the connection to the PostgreSQL
func (pg *PostgresClient) Close() {
	pg.Connection.Close()
}

// Insert inserts data to the table
func (pg *PostgresClient) Insert(m *Message) error {
	sqlStatement := `INSERT INTO tg_parser (message_id, chat_id, chat_title, content , date, views, forwards, replies) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`

	_, err := pg.Connection.Exec(sqlStatement, m.MessageID, m.ChatID, m.ChatTitle, m.Content, m.Date, m.Views, m.Forwards, m.Replies)
	if err != nil {
		return err
	}

	return nil
}

// GetAllData returns all table rows
func (pg *PostgresClient) GetAllData() ([]*Message, error) {
	rows, err := pg.Connection.Query(`SELECT * FROM tg_parser;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*Message

	for rows.Next() {
		m, err := scan(rows)
		if err != nil {
			return nil, err
		}

		messages = append(messages, m)
	}

	return messages, nil
}

// GetMessageById returns only one row with the given chat id and message id
func (pg *PostgresClient) GetMessageById(chatID int64, messageID int64) (*Message, error) {
	selectString := fmt.Sprintf(`SELECT * FROM tg_parser WHERE chat_id=%v AND message_id=%v ;`, chatID, messageID)

	rows, err := pg.Connection.Query(selectString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rows.Next()
	m, err := scan(rows)
	if err != nil {
		return nil, err
	}

	return m, err
}

// Update updates statistics and content of the message.
func (pg *PostgresClient) Update(u *UpdateRow) (updateCount int64, err error) {
	//UPDATE table_name SET column1 = value1, column2 = value2, ... WHERE condition;
	updateString := fmt.Sprintf(`UPDATE %v SET chat_title = '%v', content = '%v' , views = %v , forwards = %v, replies = %v WHERE chat_id = %v AND message_id = %v RETURNING message_id;`, pg.DbInfo.TableName, u.NewChatTitle, u.NewContent, u.NewViews, u.NewForwards, u.NewReplies, u.ChatId, u.MessageId)

	result, err := pg.Connection.Exec(updateString)
	updateCount, _ = result.RowsAffected()

	return updateCount, err
}

// TODO изменить => период (последние сутки) = время от 12.01 AM по time.Now()
// GetMessagesForATimePeriod returns messages for the selected time period.
// The list of time intervals is in the structure TimePeriods in PostgresClient
func (pg *PostgresClient) GetMessagesForATimePeriod(period string) ([]*Message, error) {
	from, to, err := dateCalculation(pg, period)
	if err != nil {
		return nil, err
	}

	selectString := fmt.Sprintf(`SELECT * FROM %v WHERE date>=%v AND date<=%v ;`, pg.DbInfo.TableName, from, to)

	rows, err := pg.Connection.Query(selectString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*Message

	for rows.Next() {
		m, err := scan(rows)
		if err != nil {
			return nil, err
		}

		messages = append(messages, m)
	}

	return messages, nil
}

//    --------------------------------------------------------------------------------
//                                     HELPERS
//    --------------------------------------------------------------------------------

// getSchemaInfo returns the names of fields in the database schema
func getSchemaInfo() *SchemaInfo {
	return &SchemaInfo{
		MessageID: "message_id",
		ChatID:    "chat_id",
		ChatTitle: "chat_title",
		Content:   "content",
		Date:      "date",
		Views:     "views",
		Forwards:  "forwards",
		Replies:   "replies",
	}
}

func getTimePeriods() *TimePeriods {
	return &TimePeriods{
		Today:              "today",
		Yesterday:          "yesterday",
		DayBeforeYesterday: "daybeforeyesterday",
		ThisWeek:           "thisweek",
		LastWeek:           "lastweek",
		ThisMonth:          "thismonth",
	}
}

// NewMessage returns a structure compatible with the database schema
func NewMessage(message *tdlib.Message, chat *tdlib.Chat) *Message {
	m := &Message{
		MessageID: message.ID,
		ChatID:    message.ChatID,
		ChatTitle: chat.Title,
		Content:   message.Content.(*tdlib.MessageText).Text.Text,
		Date:      message.Date,
	}

	if message.InteractionInfo != nil {
		m.Views = message.InteractionInfo.ViewCount
		m.Forwards = message.InteractionInfo.ForwardCount
		if message.InteractionInfo.ReplyInfo != nil {
			m.Replies = message.InteractionInfo.ReplyInfo.ReplyCount
		}
	}

	return m
}

// scan scans row data into *Message
func scan(row *sql.Rows) (*Message, error) {
	m := &Message{}

	if err := row.Scan(&m.MessageID, &m.ChatID, &m.ChatTitle, &m.Content, &m.Date, &m.Views, &m.Forwards, &m.Replies); err != nil {
		return nil, err
	}
	return m, nil
}

func dateCalculation(pg *PostgresClient, period string) (from int64, to int64, err error) {
	p := pg.TimePeriods

	from = int64(0)
	to = int64(0)
	err = nil

	switch period {
	case p.Today:
		from = time.Now().AddDate(0, 0, -1).Unix()
		to = time.Now().Unix()
	case p.Yesterday:
		from = time.Now().AddDate(0, 0, -2).Unix()
		to = time.Now().AddDate(0, 0, -1).Unix()
	case p.DayBeforeYesterday:
		from = time.Now().AddDate(0, 0, -3).Unix()
		to = time.Now().AddDate(0, 0, -2).Unix()
	case p.ThisWeek:
		from = time.Now().AddDate(0, 0, -7).Unix()
		to = time.Now().Unix()
	case p.LastWeek:
		from = time.Now().AddDate(0, 0, -14).Unix()
		to = time.Now().AddDate(0, 0, -7).Unix()
	case p.ThisMonth:
		from = time.Now().AddDate(0, -1, 0).Unix()
		to = time.Now().Unix()
	default:
		err = errors.New(fmt.Sprintf("unknown time period %v", period))
	}

	return from, to, err
}

//    --------------------------------------------------------------------------------
//                                     DB STRUCT
//    --------------------------------------------------------------------------------

/*
CREATE TABLE tg_parser (
     message_id bigint,
     chat_id bigint,
     chat_title text,
     content text,
     date bigint,
     views integer,
     forwards integer,
     replies integer
);
*/
