// Created by Noy Hillel
// https://github.com/noy
// You are free to use/modify this software in your own projects
package tsql

import (
	"testing"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func TestNewSQLClient(t *testing.T) {
	database, err := sql.Open("mysql", "username:password@tcp(localhost)/database")
	database.SetMaxIdleConns(0)
	sqlClient := NewSQLClient(database, err)
	query := sqlClient.Select("username, id").From("users").Where("username").Like("test", true)
	rows, err := query.QueryResult(true)

	var id int
	var username, password string
	for rows.Next() {
		if err := rows.Scan(&username, &id); err != nil {
			t.Log(err.Error())
		} else {
			t.Log(id, username, password) // will not print password as we did not include it in our query
		}
	}
}
