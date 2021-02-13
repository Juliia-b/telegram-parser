package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

var testPort string

//    --------------------------------------------------------------------------------
//                                  STRUCTURES
//    --------------------------------------------------------------------------------

//    ---------------------------------INSERT-----------------------------------------
var (
	// Insert
	m1 = &Message{10, 1, "T1", "chat1-message10", 5000, 12, 13, 14}
	m2 = &Message{11, 1, "T1", "chat1-message11", 7000, 40, 2, 1}
	m3 = &Message{10, 5, "T5", "chat5-message10", 3000, 99, 50, 0}
	m4 = &Message{11, 5, "T5", "chat5-message11", 7000, 20, 19, 14}
	m5 = &Message{40, 9, "T9", "chat9-message40", 2000, 50, 20, 10}
	m6 = &Message{41, 9, "T9", "chat9-message41", 3000, 40, 10, 0}
)

var testInsert = []*Message{m1, m2, m3, m4, m5, m6}

//    --------------------------------GETMESSAGEBYID----------------------------------
type testGetMessageById struct {
	MessageID       int64
	ChatID          int64
	ExpectedMessage *Message
	HasError        bool
	ErrorText       string
}

var testGetMessageByIds = []*testGetMessageById{
	// Messages with such IDs exist in the database. No error will be returned
	{10, 1, m1, false, "<nil>"},
	{10, 5, m3, false, "<nil>"},
	{40, 9, m5, false, "<nil>"},

	// Messages with such IDs do not exist in the database. An error will be returned
	{10, 20, nil, true, "some error"},
	{10, 40, nil, true, "some error"},
	{17, 9, nil, true, "some error"},

	//// Erroneous data. Should cause an error in the test
	//{10, 1, m6, true, "nil"},
	//{10, 5, nil, false, "nil"},
	//{10, 20, nil, false, "not nil"},
	//{10, 40, m6, true, "not nil"},
}

//    ------------------------------------UPDATE--------------------------------------
var (
	// To update values (after the test will be deleted from database)
	m7 = &Message{20, 70, "OLD", "OLD-cont", 7000, 30, 20, 10}

	// Will update data
	u1 = &UpdateRow{20, 70, "OLD", "OLD-cont", 35, 25, 15}
	u2 = &UpdateRow{20, 70, "NEW", "NEW-cont", 30, 20, 10}
	u3 = &UpdateRow{20, 70, "NEW", "NEW-cont", 35, 25, 15}

	// No such chat ID and message ID in database
	u4 = &UpdateRow{1, 1, "OLD", "OLD-cont", 35, 25, 15}
	u5 = &UpdateRow{2, 2, "NEW", "NEW-cont", 30, 20, 10}
	u6 = &UpdateRow{3, 3, "NEW", "OLD-cont", 35, 25, 15}
)

type testUpdate struct {
	InitMessage         *Message
	UpdateRow           *UpdateRow
	ErrorText           string
	ExpectedUpdateCount int64
}

var testUpdates = []*testUpdate{
	// The Update function will execute correctly. No error will be returned
	{m7, u1, "<nil>", 1},
	{m7, u2, "<nil>", 1},
	{m7, u3, "<nil>", 1},

	// There are no such rows in the database. An error will be returned
	{m7, u4, "<nil>", 0},
	{m7, u5, "<nil>", 0},
	{m7, u6, "<nil>", 0},

	// Erroneous data. Should cause an error in the test
	//{m7, u4, "<nil>", 1},
	//{m7, u1, "<nil>", 0},
}

//    ---------------------------GETMESSAGESFORATIMEPERIOD----------------------------
// TODO add test

//    --------------------------------------------------------------------------------
//                                     HELPERS
//    --------------------------------------------------------------------------------

func postgresTestConnection() (*PostgresClient, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://postgres:secret@localhost:%s/%s?sslmode=disable", testPort, "postgres"))
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tg_parser ( message_id bigint, chat_id bigint, chat_title text, content text, date bigint, views integer, forwards integer, replies integer, PRIMARY KEY(message_id, chat_id) );`)
	if err != nil {
		return nil, err
	}

	return &PostgresClient{Connection: db, DbInfo: &DbInfo{"postgres", "tg_parser"}, SchemaInfo: getSchemaInfo(), TimePeriods: getTimePeriods()}, nil
}

func postgresTestGetRowsCount(pg *PostgresClient) (int, error) {
	count := 0

	row := pg.Connection.QueryRow("SELECT COUNT(*) FROM tg_parser")
	err := row.Scan(&count)

	return count, err
}

func postgresTestDeleteRow(pg *PostgresClient, chatId int64, messageId int64) error {
	_, err := pg.Connection.Exec(fmt.Sprintf("DELETE FROM %v WHERE chat_id = %v AND message_id = %v;", pg.DbInfo.TableName, chatId, messageId))

	return err
}
