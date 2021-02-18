package db

//type DB interface {
//	CloseConnection()
//
//	// CRUD methods for "post" table
//	InsertMessage(m *Message) error
//	GetAllMessages() ([]*Message, error)
//	GetMessage(chatID int64, messageID int64) (*Message, error)
//	GetMessageWithPeriod(from int64, to int64, limit int) ([]*Message, error)
//	UpdateMessage(u *UpdateRow) (updateCount int64, err error)
//
//	// CRUD methods for "client" table
//	InsertClient(client *Client) error
//	GetAllClients() ([]*Client, error)
//	GetClient(cookie string) (*Client, error)
//	UpdateClient(lastCli *Client, newCookie string) (updateCount int64, err error)
//}

type DB interface {
	CloseConnection()
	DBTablePost
	DBTableClient
}

type DBTablePost interface {
	InsertMessage(m *Message) error
	GetAllMessages() ([]*Message, error)
	GetMessage(chatID int64, messageID int64) (*Message, error)
	GetMessageWithPeriod(from int64, to int64, limit int) ([]*Message, error)
	UpdateMessage(u *UpdateRow) (updateCount int64, err error)
}

type DBTableClient interface {
	InsertClient(client *Client) error
	GetAllClients() ([]*Client, error)
	GetClient(cookie string) (*Client, error)
	UpdateClient(lastCli *Client, newCookie string) (updateCount int64, err error)
}
