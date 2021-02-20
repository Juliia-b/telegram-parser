package db

import (
	"database/sql"
	"fmt"
)

// InsertMessage inserts struct Message to the table posts of tg_parser database.
func (t *TablePost) InsertMessage(message *Message) error {
	sqlStatement := fmt.Sprintf(`INSERT INTO %v (message_id, chat_id, chat_title, content , date, views, forwards, replies, link) VALUES (%v, %v, '%v', '%v', %v, %v, %v, %v, '%v');`, t.Name, message.MessageID, message.ChatID, message.ChatTitle, message.Content, message.Date, message.Views, message.Forwards, message.Replies, message.Link)

	_, err := t.Connection.Exec(sqlStatement)
	return err
}

// GetAllMessages returns all rows from the table "post" of tg_parser database.
func (t *TablePost) GetAllMessages() ([]*Message, error) {
	var sqlStatement = fmt.Sprintf(`SELECT * FROM %v;`, t.Name)

	rows, err := t.Connection.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*Message

	for rows.Next() {
		m, err := scanPosts(rows)
		if err != nil {
			return nil, err
		}

		messages = append(messages, m)
	}

	return messages, nil
}

// GetMessage returns only one post with the given chat id and message id.
func (t *TablePost) GetMessage(chatID int64, messageID int64) (*Message, error) {
	var sqlStatement = fmt.Sprintf(`SELECT * FROM %v WHERE chat_id=%v AND message_id=%v ;`, t.Name, chatID, messageID)

	rows, err := t.Connection.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rows.Next()
	m, err := scanPosts(rows)
	if err != nil {
		return nil, err
	}

	return m, err
}

// UpdateMessage updates statistics and content of the message in the table "post" of "tg_parser" database.
func (t *TablePost) UpdateMessage(u *UpdateRow) (updateCount int64, err error) {
	var sqlStatement = fmt.Sprintf(`UPDATE %v SET chat_title = '%v', content = '%v' , views = %v , forwards = %v, replies = %v, date = %v, link = %v WHERE chat_id = %v AND message_id = %v RETURNING message_id;`, t.Name, u.NewChatTitle, u.NewContent, u.NewViews, u.NewForwards, u.NewReplies, u.NewDate, u.NewLink, u.ChatID, u.MessageID)

	result, err := t.Connection.Exec(sqlStatement)
	if err != nil {
		return 0, err
	}
	updateCount, _ = result.RowsAffected()

	return updateCount, err
}

// GetMessageWithPeriod returns messages for the selected time period.
// The list of time intervals is in the structure TimePeriods.
func (t *TablePost) GetMessageWithPeriod(from int64, to int64, limit int) ([]*Message, error) {
	var sqlStatement = fmt.Sprintf(`SELECT * FROM %v WHERE date>=%v AND date<=%v AND views>1 ORDER BY views DESC, forwards DESC, replies DESC  LIMIT %v;`, t.Name, from, to, limit)

	rows, err := t.Connection.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*Message

	for rows.Next() {
		m, err := scanPosts(rows)
		if err != nil {
			return nil, err
		}

		messages = append(messages, m)
	}

	return messages, nil
}

// DeleteMessage deletes message from the table "post" of tg_parser database.
func (t *TablePost) DeleteMessage(message *Message) (deleteCount int64, err error) {
	var sqlStatement = fmt.Sprintf(`DELETE FROM %v WHERE message_id = %v AND chat_id = %v`, t.Name, message.MessageID, message.ChatID)

	result, err := t.Connection.Exec(sqlStatement)
	if err != nil {
		return 0, err
	}

	deleteCount, err = result.RowsAffected()

	return deleteCount, nil
}

/*-----------------------------------HELPERS-----------------------------------------*/

// scanPosts scans rows from table "post" to Message struct.
func scanPosts(rows *sql.Rows) (*Message, error) {
	m := &Message{}

	if err := rows.Scan(&m.MessageID, &m.ChatID, &m.ChatTitle, &m.Content, &m.Date, &m.Views, &m.Forwards, &m.Replies, &m.Link); err != nil {
		return nil, err
	}
	return m, nil
}
