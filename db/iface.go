package db

type DB interface {
	CloseConnection()
	DBTablePost
	DBTableClient
	DBTableTop3hour
}

type DBTablePost interface {
	InsertMessage(message *Message) error
	GetAllMessages() ([]*Message, error)
	GetMessage(chatID int64, messageID int64) (*Message, error)
	GetMessageWithPeriod(from int64, to int64, limit int) ([]*Message, error)
	UpdateMessage(u *UpdateRow) (updateCount int64, err error)
	DeleteMessage(message *Message) (deleteCount int64, err error)
}

type DBTableClient interface {
	InsertClient(client *Client) error
	GetAllClients() ([]*Client, error)
	GetClient(cookie string) (*Client, error)
	UpdateClient(lastCli *Client, newCookie string) (updateCount int64, err error)
}

type DBTableTop3hour interface {
	InsertTop3hour(message *Message) error
	GetAllTop3hour() ([]*Message, error)
	DeleteAllTop3hour() (deleteCount int64, err error)
}
