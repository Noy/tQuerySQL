# tSQL
Very Basic MySQL query handler written in go. the t stands for Toffee (Long time project I've worked on)

### Example

```go
package main

import (
	"github.com/Noy/tsql"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	database, err := sql.Open("mysql", "username:password@tcp(localhost)/database")
	sqlClient := tsql.NewSQLClient(database, err)
	// SELECT * FROM Users WHERE firsname="noy" AND lastname LIKE "%h%"
	query := sqlClient.
		Select("*").
		From("Users").
		Where("firstname").
		Equals("noy", true).
		And("lastname").
		Like("h", true)
	
	rows, err := query.QueryResult(true) // debug option
	if err != nil {
		// handle accordingly
	}
	var firstName, lastName string
	for rows.Next() {
		if err := rows.Scan(&firstName, &lastName); err != nil {
			// handle accordingly
		}
		// do what you will with these results
		log.Println(firstName, lastName)
	}
}

```

### You can also do something like

```go
package main

import (
	"github.com/Noy/tsql"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	database, err := sql.Open("mysql", "username:password@tcp(localhost)/database")
	sqlClient := tsql.NewSQLClient(database, err)
	// SELECT * FROM Users WHERE firsname="noy" AND lastname LIKE "%h%"
	res, err := sqlClient.
		InsertInto("SomeTable").
		// String values will need extra quotes
		Values(`"Noy"`, `"H"`, `1.99`, `76`).
		Execute(true) // debug option
	if err != nil {
		// handle
	}
	// Do what you will with the result, for executions, you can replace the var name with _ if you don't need it..
	// ..and just want to check the error.
	log.Println(res)
}
```

##### It simply joins the strings, so nothing special, but it will panic if you do something like:

```go
client.Select("*").And("example") // will not work
```



