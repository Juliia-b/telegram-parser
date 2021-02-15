package db

type DB interface {
	Close()
	Insert(m *Message) error
	GetAllData() ([]*Message, error)
	GetMessageById(chatID int64, messageID int64) (*Message, error)
	Update(u *UpdateRow) (updateCount int64, err error)
	GetMessagesForATimePeriod(from int64, to int64) ([]*Message, error)
}
