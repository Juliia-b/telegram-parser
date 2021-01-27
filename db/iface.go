package db

type DB interface {
	Close()
	Insert(m *Message) error
	GetAllData() ([]*Message, error)
	GetMessageById(chatID int64, messageID int64) (*Message, error)
	Update(u *UpdateRow) (updateCount int64, err error)
	GetMessagesForATimePeriod(period string) ([]*Message, error)
}
