package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

// ConnectToPostgres opens a connection to PostgreSQL.
func ConnectToPostgres() (*PostgresClient, error) {
	var dbName = "tg_parser"

	pgInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("PGHOST"), os.Getenv("PGPORT"), os.Getenv("PGUSER"), os.Getenv("PGPASSWORD"), dbName)

	dbConn, err := sql.Open("postgres", pgInfo)
	if err != nil {
		return nil, err
	}

	err = dbConn.Ping()
	if err != nil {
		return nil, err
	}

	err = createTables(dbConn)
	if err != nil {
		return nil, err
	}

	return &PostgresClient{
		Connection:  dbConn,
		DbName:      dbName,
		TableClient: getClientTableStruct(dbConn),
		TablePost:   getPostTableStruct(dbConn),
	}, nil
}

// Close closes the connection to the PostgreSQL.
func (pg *PostgresClient) CloseConnection() {
	pg.Connection.Close()
}

/*-----------------------------------HELPERS-----------------------------------------*/

// createTables creates all the necessary tables in the database, if they have not been created yet.
// Table names: "client" , "post".
func createTables(db *sql.DB) (err error) {
	var sqlStatement string

	// creating table "post"
	sqlStatement = `CREATE TABLE IF NOT EXISTS post ( message_id bigint, chat_id bigint, chat_title text, content text, date bigint, views integer, forwards integer, replies integer, UNIQUE(message_id, chat_id) );`

	_, err = db.Exec(sqlStatement)
	if err != nil {
		return err
	}

	// creating table "client"
	sqlStatement = `CREATE TABLE IF NOT EXISTS client ( id SERIAL, cookie text, PRIMARY KEY(id) );`

	_, err = db.Exec(sqlStatement)
	if err != nil {
		return err
	}

	return nil
}
