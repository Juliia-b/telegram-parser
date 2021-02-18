package db

import (
	"database/sql"
	"fmt"
)

// InsertClient inserts struct Client to the table "client" of tg_parser database.
func (c *TableClient) InsertClient(client *Client) error {
	sqlStatement := fmt.Sprintf(`INSERT INTO %v (cookie) VALUES (%v);`, c.Name, client.Cookie)

	_, err := c.Connection.Exec(sqlStatement)
	return err
}

/*-----------------------------------------------------------------------------------*/

// GetAllClients returns all rows from the table "post" of tg_parser database.
func (c *TableClient) GetAllClients() ([]*Client, error) {
	var sqlStatement = fmt.Sprintf(`SELECT * FROM %v;`, c.Name)

	rows, err := c.Connection.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []*Client

	for rows.Next() {
		client, err := scanClients(rows)
		if err != nil {
			return nil, err
		}

		clients = append(clients, client)
	}

	return clients, nil
}

/*-----------------------------------------------------------------------------------*/

// GetClient returns only one client with the given cookie.
func (c *TableClient) GetClient(cookie string) (*Client, error) {
	var sqlStatement = fmt.Sprintf(`SELECT * FROM %v WHERE cookie=%v ;`, c.Name, cookie)

	rows, err := c.Connection.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rows.Next()
	m, err := scanClients(rows)
	if err != nil {
		return nil, err
	}

	return m, err
}

/*-----------------------------------------------------------------------------------*/

// UpdateClient updates cookie value of a client in the table "client" of "tg_parser" database.
func (c *TableClient) UpdateClient(lastCli *Client, newCookie string) (updateCount int64, err error) {
	var sqlStatement = fmt.Sprintf(`UPDATE %v SET cookie = '%v' WHERE cookie = %v OR id = %v;`, c.Name, newCookie, lastCli.Cookie, lastCli.ID)

	result, err := c.Connection.Exec(sqlStatement)
	updateCount, _ = result.RowsAffected()

	return updateCount, err
}

/*-----------------------------------HELPERS-----------------------------------------*/

// scanClients scans rows from table "client" to Client struct.
func scanClients(rows *sql.Rows) (*Client, error) {
	c := &Client{}

	if err := rows.Scan(&c.ID, &c.Cookie); err != nil {
		return nil, err
	}
	return c, nil
}
