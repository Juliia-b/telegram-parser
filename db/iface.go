package db

type DB interface {
	CloseConnection()
	DBTablePost
}

type DBTablePost interface {
	InsertMessage(message *Message) error
	GetAllMessages() ([]*Message, error)
	GetMessage(chatID int64, messageID int64) (*Message, error)
	GetMessageWithPeriod(from int64, to int64, limit int) ([]*Message, error)
	UpdateMessage(u *UpdateRow) (updateCount int64, err error)
	DeleteMessage(message *Message) (deleteCount int64, err error)
}
