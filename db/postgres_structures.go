package db

import (
	"database/sql"
	"github.com/Arman92/go-tdlib"
)

/*---------------------------------STRUCTURES----------------------------------------*/

// PostgresClient contains connection to "tg_parser" database and all tables methods.
type PostgresClient struct {
	Connection *sql.DB
	DbName     string

	// tables in database
	*TableClient
	*TablePost
}

// Tables contains list of available tables in "tg_parser" database.
type Tables struct {
	Post   *TablePost
	Client *TableClient
}

type TablePost struct {
	*Table
}

type TableClient struct {
	*Table
}

// Table contains a list of available tables for the database.
type Table struct {
	Name       string
	Connection *sql.DB
}

// Message structure compatible with table "post" schema.
// Used to create and read rows.
type Message struct {
	MessageID int64  `json:"messageid"`
	ChatID    int64  `json:"chatid"`
	Date      int64  `json:"date"`
	ChatTitle string `json:"chattitle"`
	Content   string `json:"content"`
	Views     int32  `json:"views"`
	Forwards  int32  `json:"forwards"`
	Replies   int32  `json:"replies"`
	Link      string `json:"link"`
}

// UpdateRow structure compatible with table "post" schema.
// Used to update rows.
type UpdateRow struct {
	MessageID    int64
	ChatID       int64
	NewDate      int64
	NewChatTitle string
	NewContent   string
	NewViews     int32
	NewForwards  int32
	NewReplies   int32
	NewLink      string
}

// Client structure compatible with table "client" schema.
// Used to create, read and update row.
type Client struct {
	ID     int    // auto increment
	Cookie string // the unix time of the first visit of the user to the site is unique as string
}

/*-----------------------------------HELPERS-----------------------------------------*/

// NewMessage returns a structure compatible with the database schema.
func NewMessage(message *tdlib.Message, chatTitle string, link *tdlib.MessageLink) *Message {
	m := &Message{
		MessageID: message.ID,
		ChatID:    message.ChatID,
		ChatTitle: chatTitle,
		Content:   message.Content.(*tdlib.MessageText).Text.Text,
		Date:      int64(message.Date),
		Link:      link.Link,
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

// getPostTableStruct returns the filled structure TablePost.
func getPostTableStruct(dbConn *sql.DB) *TablePost {
	return &TablePost{&Table{
		Name:       "post",
		Connection: dbConn,
	}}
}

// getClientTableStruct returns the filled structure TableClient.
func getClientTableStruct(dbConn *sql.DB) *TableClient {
	return &TableClient{&Table{
		Name:       "client",
		Connection: dbConn,
	}}
}

/*----------------------------------DB STRUCT----------------------------------------*/

/*

DATABASE: tg_parser

TABLES:

1. client {
     id          SERIAL   NOT NULL  - auto increment user id
     cookie      text     NOT NULL  - user cookie

     PRIMARY KEY(id)
     UNIQUE(cookie)
}

2. post {                   in last tg_parser
     message_id  bigint    NOT NULL  -
     chat_id     bigint    NOT NULL  -
     chat_title  text      NOT NULL  -
     content     text      NOT NULL  -
     date        bigint    NOT NULL  -
     views       integer   NOT NULL  -
     forwards    integer   NOT NULL  -
     replies     integer   NOT NULL  -
     link        text      NOT NULL  -

     UNIQUE (message_id, chat_id) - Такое ограничение указывает,
                                    что сочетание значений перечисленных столбцов должно быть уникально во  всей таблице,
                                    тогда как значения каждого столбца по отдельности не должны быть (и обычно не будут)
                                    уникальными.

     // deprecate PRIMARY KEY(message_id, chat_id)
}

*/
