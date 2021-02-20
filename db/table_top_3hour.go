package db

import "fmt"

// InsertTop3hour inserts struct Message to the table "top_3_hours" of tg_parser database.
func (t *TableTop3Hour) InsertTop3hour(message *Message) error {
	sqlStatement := fmt.Sprintf(`INSERT INTO %v (message_id, chat_id, chat_title, content , date, views, forwards, replies, link) VALUES (%v, %v, '%v', '%v', %v, %v, %v, %v, '%v');`, t.Name, message.MessageID, message.ChatID, message.ChatTitle, message.Content, message.Date, message.Views, message.Forwards, message.Replies, message.Link)

	_, err := t.Connection.Exec(sqlStatement)
	return err
}

// GetAllTop3hour returns all rows from the table "top_3_hours" of tg_parser database.
func (t *TableTop3Hour) GetAllTop3hour() ([]*Message, error) {
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

// DeleteAllTop3hour deletes all message from the table "top_3_hours" of tg_parser database.
func (t *TableTop3Hour) DeleteAllTop3hour() (deleteCount int64, err error) {
	var sqlStatement = fmt.Sprintf(`DELETE FROM %v`, t.Name)

	result, err := t.Connection.Exec(sqlStatement)
	if err != nil {
		return 0, err
	}

	deleteCount, err = result.RowsAffected()

	return deleteCount, nil
}
