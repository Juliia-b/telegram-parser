package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Arman92/go-tdlib"
	_ "github.com/lib/pq"
	"os"
)

//    --------------------------------------------------------------------------------
//                                  STRUCTURES
//    --------------------------------------------------------------------------------

type PostgresClient struct {
	Connection *sql.DB
	DbInfo     *DbInfo
}

type DbInfo struct {
	DbName    string
	TableName string
}

type Message struct {
	MessageID int64  `json:"messageid"`
	ChatID    int64  `json:"chatid"`
	Date      int64  `json:"date"`
	ChatTitle string `json:"chattitle"`
	Content   string `json:"content"`
	Views     int32  `json:"views"`
	Forwards  int32  `json:"forwards"`
	Replies   int32  `json:"replies"`
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

	return &PostgresClient{Connection: db, DbInfo: &DbInfo{"postgres", tableName}}, nil
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

	return updateCount, err
}

// GetMessagesForATimePeriod returns messages for the selected time period.
// The list of time intervals is in the structure TimePeriods.
func (pg *PostgresClient) GetMessagesForATimePeriod(from int64, to int64, limit int) ([]*Message, error) {
	var sqlStatement = fmt.Sprintf(`SELECT * FROM %v WHERE date>=%v AND date<=%v AND views>1 ORDER BY views DESC, forwards DESC, replies DESC  LIMIT %v;`, pg.DbInfo.TableName, from, to, limit)

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

// NewMessage returns a structure compatible with the database schema.
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

// scan scans row data into *Message.
func scan(row *sql.Rows) (*Message, error) {
	m := &Message{}

	if err := row.Scan(&m.MessageID, &m.ChatID, &m.ChatTitle, &m.Content, &m.Date, &m.Views, &m.Forwards, &m.Replies); err != nil {
		return nil, err
	}
	return m, nil
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
