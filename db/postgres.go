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
//                                  STRUCTURES
//    --------------------------------------------------------------------------------

type PostgresClient struct {
	Connection  *sql.DB
	DbInfo      *DbInfo
	TimePeriods *TimePeriods
}

type DbInfo struct {
	DbName    string
	TableName string
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
	Date      int64
	ChatTitle string
	Content   string
	Views     int32
	Forwards  int32
	Replies   int32
}

type UpdateRow struct {
	MessageID    int64
	ChatID       int64
	NewDate      int64
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
	pgInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("PGHOST"), os.Getenv("PGPORT"), os.Getenv("PGUSER"), os.Getenv("PGPASSWORD"), os.Getenv("PGDBNAME"))

	db, err := sql.Open("postgres", pgInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	var tableName = "tg_parser"
	var sqlStatement = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %v ( message_id bigint, chat_id bigint, chat_title text, content text, date bigint, views integer, forwards integer, replies integer, PRIMARY KEY(message_id, chat_id) );`, tableName)

	res, err := db.Exec(sqlStatement)
	if err != nil {
		return nil, err
	}

	// count can be 0 if table already EXISTS, else count 1
	if count, _ := res.RowsAffected(); count > 1 {
		return nil, errors.New("table connection error")
	}

	//TODO check if it is possible to reduce or simplify the structure PostgresClient
	return &PostgresClient{Connection: db, DbInfo: &DbInfo{"postgres", tableName}, TimePeriods: getTimePeriods()}, nil
}

// Close closes the connection to the PostgreSQL
func (pg *PostgresClient) Close() {
	pg.Connection.Close()
}

// Insert inserts data to the table
func (pg *PostgresClient) Insert(m *Message) error {
	sqlStatement := fmt.Sprintf(`INSERT INTO %v (message_id, chat_id, chat_title, content , date, views, forwards, replies) VALUES (%v, %v, '%v', '%v', %v, %v, %v, %v);`, pg.DbInfo.TableName, m.MessageID, m.ChatID, m.ChatTitle, m.Content, m.Date, m.Views, m.Forwards, m.Replies)

	_, err := pg.Connection.Exec(sqlStatement)
	return err
}

// GetAllData returns all table rows
func (pg *PostgresClient) GetAllData() ([]*Message, error) {
	var sqlStatement = fmt.Sprintf(`SELECT * FROM %v;`, pg.DbInfo.TableName)

	rows, err := pg.Connection.Query(sqlStatement)
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
	var sqlStatement = fmt.Sprintf(`SELECT * FROM %v WHERE chat_id=%v AND message_id=%v ;`, pg.DbInfo.TableName, chatID, messageID)

	rows, err := pg.Connection.Query(sqlStatement)
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
	var sqlStatement = fmt.Sprintf(`UPDATE %v SET chat_title = '%v', content = '%v' , views = %v , forwards = %v, replies = %v, date = %v WHERE chat_id = %v AND message_id = %v RETURNING message_id;`, pg.DbInfo.TableName, u.NewChatTitle, u.NewContent, u.NewViews, u.NewForwards, u.NewReplies, u.NewDate, u.ChatID, u.MessageID)

	result, err := pg.Connection.Exec(sqlStatement)
	updateCount, _ = result.RowsAffected()

	//TODO зачем возвращать количество обновленных???
	return updateCount, err
}

// GetMessagesForATimePeriod returns messages for the selected time period.
// The list of time intervals is in the structure TimePeriods in PostgresClient
func (pg *PostgresClient) GetMessagesForATimePeriod(period string) ([]*Message, error) {
	from, to, err := dateCalculation(pg, period)
	if err != nil {
		return nil, err
	}

	var sqlStatement = fmt.Sprintf(`SELECT * FROM %v WHERE date>=%v AND date<=%v;`, pg.DbInfo.TableName, from, to)

	rows, err := pg.Connection.Query(sqlStatement)
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

// NewMessage returns a structure compatible with the database schema
func NewMessage(message *tdlib.Message, chatTitle string) *Message {
	m := &Message{
		MessageID: message.ID,
		ChatID:    message.ChatID,
		ChatTitle: chatTitle,
		Content:   message.Content.(*tdlib.MessageText).Text.Text,
		Date:      int64(message.Date),
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

	// TODO уточнить временные периоды
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
