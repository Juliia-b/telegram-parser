package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"os"
)

// ConnectToPostgres opens a connection to PostgreSQL.
func ConnectToPostgres() *PostgresClient {
	var dbName = "tg_parser"

	pgInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("PGHOST"), os.Getenv("PGPORT"), os.Getenv("PGUSER"), os.Getenv("PGPASSWORD"), dbName)

	dbConn, err := sql.Open("postgres", pgInfo)
	if err != nil {
		logrus.Fatal(err)
	}

	err = dbConn.Ping()
	if err != nil {
		logrus.Fatal(err)
	}

	err = createTables(dbConn)
	if err != nil {
		logrus.Fatal(err)
	}

	return &PostgresClient{
		Connection:    dbConn,
		DbName:        dbName,
		TableClient:   getClientTableStruct(dbConn),
		TablePost:     getPostTableStruct(dbConn),
		TableTop3Hour: getTop3HourTableStruct(dbConn),
	}
}

// Close closes the connection to the PostgreSQL.
func (pg *PostgresClient) CloseConnection() {
	pg.Connection.Close()
}

/*-----------------------------------HELPERS-----------------------------------------*/

// createTables creates all the necessary tables in the database, if they have not been created yet.
// Table names: "client" , "post".
func createTables(db *sql.DB) (err error) {
	if err = createPostTable(db); err != nil {
		return err
	}

	if err = createClientTable(db); err != nil {
		return err
	}

	return nil
}

// createPostTable creates table "post" in database.
func createPostTable(db *sql.DB) (err error) {
	sqlStatement := `CREATE TABLE IF NOT EXISTS post ( message_id bigint NOT NULL , chat_id bigint NOT NULL , chat_title text NOT NULL , content text NOT NULL , date bigint NOT NULL , views integer NOT NULL , forwards integer NOT NULL , replies integer NOT NULL , link text NOT NULL , UNIQUE(message_id, chat_id) );`

	_, err = db.Exec(sqlStatement)
	if err != nil {
		return err
	}

	return nil
}

// createClientTable creates table "client" in database.
func createClientTable(db *sql.DB) (err error) {
	sqlStatement := `CREATE TABLE IF NOT EXISTS client ( id SERIAL, cookie text  NOT NULL , PRIMARY KEY(id) );`

	_, err = db.Exec(sqlStatement)
	if err != nil {
		return err
	}

	return nil
}
