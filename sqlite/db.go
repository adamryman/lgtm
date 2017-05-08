package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	. "github.com/y0ssar1an/q"
)

func main() {
	fmt.Println("You can do anything!")
	Q("Lets debug some shit")
}

func Open(conn string) (*Database, error) {
	d, err := sql.Open("sqlite3", conn)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot connect to mysql with %s", conn)
	}
	if err := d.Ping(); err != nil {
		return nil, errors.Wrapf(err, "cannot make initial database connection to %s", conn)
	}

	if err := setupDB(d); err != nil {
		return nil, err
	}

	return &Database{d}, nil
}

func setupDB(db *sql.DB) error {
	const users = `CREATE TABLE IF NOT EXISTS users(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				slack_id text,
				oauth_token text);`
	_, err := db.Exec(users)
	if err != nil {
		return err
	}

	const hooks = `CREATE TABLE IF NOT EXISTS hooks(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				hook_id integer,
				user_id integer);`
	_, err = db.Exec(hooks)
	if err != nil {
		return err
	}

	const pullRequests = `CREATE TABLE IF NOT EXISTS pull_requests(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
			 	pull_request_id integer,
				user_id integer,
				timestamp text);`
	_, err = db.Exec(pullRequests)
	if err != nil {
		return err
	}

	return nil
}

type Database struct {
	db *sql.DB
}

func (d *Database) GetUserAuth(slackId string) (string, error) {
	const query = `SELECT * FROM users WHERE slack_id=?`
	resp := d.db.QueryRow(query, slackId)

	var id int
	var token string
	err := resp.Scan(&id, &slackId, &token)
	if err != nil {
		return "", err
	}

	return token, nil
}
