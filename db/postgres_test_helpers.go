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
	m1 = &Message{10, 1, 5000, "T1", "chat1-message10", 12, 13, 14}
	m2 = &Message{11, 1, 7000, "T1", "chat1-message11", 40, 2, 1}
	m3 = &Message{10, 5, 3000, "T5", "chat5-message10", 99, 50, 0}
	m4 = &Message{11, 5, 7000, "T5", "chat5-message11", 20, 19, 14}
	m5 = &Message{40, 9, 2000, "T9", "chat9-message40", 50, 20, 10}
	m6 = &Message{41, 9, 3000, "T9", "chat9-message41", 40, 10, 0}
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
	// Message over which the fields will be updated (after the test will be deleted from database)
	m7 = &Message{20, 70, 7000, "OLD", "OLD-cont", 30, 20, 10}

	// Will update data
	u1 = &UpdateRow{20, 70, 7000, "OLD", "OLD-cont", 35, 25, 15}
	u2 = &UpdateRow{20, 70, 8000, "NEW", "NEW-cont", 30, 20, 10}
	u3 = &UpdateRow{20, 70, 9000, "NEW", "NEW-cont", 35, 25, 15}

	// No such chat ID and message ID in database
	u4 = &UpdateRow{1, 1, 7000, "OLD", "OLD-cont", 35, 25, 15}
	u5 = &UpdateRow{2, 2, 8000, "NEW", "NEW-cont", 30, 20, 10}
	u6 = &UpdateRow{3, 3, 9000, "NEW", "OLD-cont", 35, 25, 15}
)

type testUpdate struct {
	InitialMessage      *Message
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

// postgresTestConnection creates a connection to a test database in a docker container
func postgresTestConnection() (*PostgresClient, error) {
	var driverName = "postgres"
	var tableName = "tg_parser"
	var dbName = "postgres"

	var dataSourceName = fmt.Sprintf("postgres://postgres:secret@localhost:%s/%s?sslmode=disable", testPort, "postgres")

	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	var sqlStatement = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %v ( message_id bigint, chat_id bigint, chat_title text, content text, date bigint, views integer, forwards integer, replies integer, PRIMARY KEY(message_id, chat_id) );`, tableName)

	_, err = db.Exec(sqlStatement)
	if err != nil {
		return nil, err
	}

	return &PostgresClient{
		Connection:  db,
		DbInfo:      &DbInfo{dbName, tableName},
		TimePeriods: getTimePeriods(),
	}, nil
}

// postgresTestGetRowsCount returns the number of fields in a test database in a docker container
func postgresTestGetRowsCount(pg *PostgresClient) (int, error) {
	var count int
	var sqlStatement = fmt.Sprintf("SELECT COUNT(*) FROM %v", pg.DbInfo.TableName)

	row := pg.Connection.QueryRow(sqlStatement)
	err := row.Scan(&count)

	return count, err
}

// postgresTestDeleteRow deletes a field in a test database in a docker container
func postgresTestDeleteRow(pg *PostgresClient, chatId int64, messageId int64) error {
	var sqlStatement = fmt.Sprintf("DELETE FROM %v WHERE chat_id = %v AND message_id = %v;", pg.DbInfo.TableName, chatId, messageId)

	_, err := pg.Connection.Exec(sqlStatement)

	return err
}
