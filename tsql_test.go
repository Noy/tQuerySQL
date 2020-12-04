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

func TestClient_From(t *testing.T) {
	// tests a little bit more than just from but still
	statement := "SELECT * FROM test WHERE test = 'test'"
	client := NewSQLClient(nil, nil)
	tSQLStatement := client.Select("*").From("test").Where("test").Equals("test", true)
	if tSQLStatement.Query == statement {
		t.Log("Tests passed: " + tSQLStatement.Query)
	} else {
		t.Errorf("Something went wrong. TsqlStatement=%v and query string=%v", tSQLStatement.Query, statement)
	}
}

func TestClient_BigQuery(t *testing.T) {
	// Don't even think this query makes sense lol
	statement := "SELECT one,two,three FROM table.test WHERE test != 3 AND test1 = 5 AND test2 >= 50 " +
		"LEFT OUTER JOIN table.test3 ON table.test.test2 = table.test3.test2 GROUP BY table.test ORDER BY table.test.test1 DESC"
	client := NewSQLClient(nil, nil)
	tSQLStatement := client.Select("one,two,three").From("table.test").Where("test").NotEqual(3).And("test1").Equals("5", false).
		And("test2").GTE("50", false).LeftOuterJoin("table.test3").On("table.test.test2 = table.test3.test2").GroupBy("table.test").
		OrderBy("table.test.test1", false)
	if tSQLStatement.Query == statement {
		t.Log("Tests passed: " + tSQLStatement.Query)
	} else {
		t.Errorf("Something went wrong. TsqlStatement=%v and query string=%v", tSQLStatement.Query, statement)
	}
}

func TestClient_GroupBy(t *testing.T) {
	statement := "SELECT * FROM test WHERE s = 1 GROUP BY one,two,three"
	client := NewSQLClient(nil, nil)
	tSQLStatement := client.Select("*").From("test").Where("s").Equals("1", false).GroupBy("one", "two", "three")
	if tSQLStatement.Query == statement {
		t.Log("Tests passed: " + tSQLStatement.Query)
	} else {
		t.Errorf("Something went wrong. TsqlStatement=%v and query string=%v", tSQLStatement.Query, statement)
	}
}

func TestClient_Values(t *testing.T) {
	statement := "INSERT INTO test VALUES('one','two','three')"
	client := NewSQLClient(nil, nil)
	tSQLStatement := client.InsertInto("test").Values("'one'", "'two'", "'three'")
	if tSQLStatement.Query == statement {
		t.Log("Tests passed: " + tSQLStatement.Query)
	} else {
		t.Errorf("Something went wrong. TsqlStatement=%v and query string=%v", tSQLStatement.Query, statement)
	}
}