// Created by Noy Hillel
// https://github.com/noy
// You are free to use/modify this software in your own projects
package tsql

import (
	"database/sql"
	utils "github.com/TrafficLabel/Go-Utilities"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type Client struct {
	Database *sql.DB
	Error error

	Query string
}

func NewSQLClient(db *sql.DB, err error) *Client {
	return &Client{db, err, ""}
}

func (c *Client) Select(value string) *Client {
	c.Query += "SELECT " + value
	return c
}

func (c *Client) From(table string) *Client {
	if c.Query == "" {
		handleTSQLError("FROM")
	}
	c.Query += " FROM " + table
	return c
}

func (c *Client) Where(value string) *Client {
	if c.Query == "" {
		handleTSQLError("WHERE")
	}
	c.Query += " WHERE " + value
	return c
}

func (c *Client) Equals(value string, quotes bool) *Client {
	if c.Query == "" {
		handleTSQLError("EQUALS")
	}
	if quotes {
		c.Query += " = '"+value+"'"
	} else {
		c.Query += " = "+value+""
	}
	return c
}

func (c *Client) Like(value string, percentSign bool) *Client {
	if c.Query == "" {
		handleTSQLError("LIKE")
	}
	if percentSign {
		c.Query += " LIKE '%"+value+"%'"
	} else {
		c.Query += " LIKE "+value+""
	}
	return c
}

func (c *Client) NotLike(value string, percentSign bool) *Client {
	if c.Query == "" {
		handleTSQLError("NOT LIKE")
	}
	if percentSign {
		c.Query += " NOT LIKE '%"+value+"%'"
	} else {
		c.Query += " NOT LIKE '"+value+"'"
	}
	return c
}

func (c *Client) NotEqual(value int) *Client {
	if c.Query == "" {
		handleTSQLError("NOT EQUAL")
	}
	c.Query += " != " + utils.String(value)
	return c
}

func (c *Client) LT(value string, quotes bool) *Client {
	if c.Query == "" {
		handleTSQLError("<")
	}
	if quotes {
		c.Query += " < '" + value + "'"
	} else {
		c.Query += " < " + value
	}
	return c
}

func (c *Client) GT(value string, quotes bool) *Client {
	if c.Query == "" {
		handleTSQLError(">")
	}
	if quotes {
		c.Query += " > '" + value + "'"
	} else {
		c.Query += " > " + value
	}
	return c
}

func (c *Client) GTE(value string, quotes bool) *Client {
	if c.Query == "" {
		handleTSQLError(">=")
	}
	if quotes {
		c.Query += " >= '" + value + "'"
	} else {
		c.Query += " >= " + value
	}
	return c
}

func (c *Client) LTE(value string, quotes bool) *Client {
	if c.Query == "" {
		handleTSQLError("<=")
	}
	if quotes {
		c.Query += " <= '" + value + "'"
	} else {
		c.Query += " <= " + value
	}
	return c
}

func (c *Client) And(value string) *Client {
	if c.Query == "" {
		handleTSQLError("AND")
	}
	c.Query += " AND " + value
	return c
}

func (c *Client) Or(value string) *Client {
	if c.Query == "" {
		handleTSQLError("OR")
	}
	c.Query += " OR " + value
	return c
}


func (c *Client) Union() *Client {
	if c.Query == "" {
		handleTSQLError("UNION")
	}
	c.Query += " UNION "
	return c
}

func (c *Client) UnionAll() *Client {
	if c.Query == "" {
		handleTSQLError("UNION ALL")
	}
	c.Query += " UNION ALL "
	return c
}

func (c *Client) IsNotNull() *Client {
	if c.Query == "" {
		handleTSQLError("IS NOT NULL")
	}
	c.Query += " IS NOT NULL "
	return c
}

func (c *Client) GroupBy(vals ...interface{}) *Client {
	if c.Query == "" {
		handleTSQLError("GROUP BY")
	}
	c.Query += " GROUP BY "
	for _, v := range vals {
		c.Query += v.(string) + ","
	}
	ql := len(c.Query)
	if ql > 0 && c.Query[ql-1] == ',' {
		c.Query = c.Query[:ql-1]
	}
	return c
}

func (c *Client) Limit(amount string) *Client {
	if c.Query == "" {
		handleTSQLError("LIMIT")
	}
	c.Query += " LIMIT " + amount
	return c
}

func (c *Client) Offset(number string) *Client {
	if c.Query == "" {
		handleTSQLError("OFFSET")
	}
	c.Query += " OFFSET " + number
	return c
}

func (c *Client) OrderBy(value string, asc bool) *Client {
	if c.Query == "" {
		handleTSQLError("ORDER BY")
	}
	c.Query += " ORDER BY " + value
	if asc {
		c.Query += " ASC"
	} else {
		c.Query += " DESC"
	}
	return c
}

func (c *Client) InsertInto(table string) *Client {
	if c.Query != "" {
		handleTSQLError("INSERT")
	}
	c.Query += "INSERT INTO " + table
	return c
}

func (c *Client) Update(table string) *Client {
	if c.Query != "" {
		handleTSQLError("UPDATE")
	}
	c.Query += "UPDATE " + table
	return c
}

func (c *Client) Set(val string) *Client {
	if c.Query == "" {
		handleTSQLError("SET")
	}
	c.Query += " SET " + val
	return c
}

func (c *Client) Values(values ...string) *Client {
	if c.Query == "" {
		handleTSQLError("VALUES")
	}
	c.Query += " VALUES("
	for _, s := range values {
		c.Query += s + ","
	}
	ql := len(c.Query)
	if ql > 0 && c.Query[ql-1] == ',' {
		c.Query = c.Query[:ql-1]
	}
	c.Query += ")"
	return c
}

func handleTSQLError(queryType string) {
	panic("You have an error in your sql syntax. '"+queryType+"' should not be here.")
}

func (c *Client) Execute(debug bool) (sql.Result, error) {
	r, err := c.Database.Exec(c.Query)
	if debug {
		log.Println(c.Query)
		if err != nil {
			log.Println(err.Error())
		}
	}
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (c *Client) QueryResult(debug bool) (*sql.Rows, error) {
	query, err := c.Database.Query(c.Query)
	if debug {
		log.Println(c.Query)
		if err != nil {
			log.Println(err.Error())
		}
	}
	if err != nil {
		return nil, err
	}
	return query, nil
}

// Functions

func (c *Client) Lower(val string, quotes bool) *Client {
	if c.Query == "" {
		handleTSQLError("LOWER")
	}
	if quotes {
		c.Query += " LOWER('" + val + "')"
	} else {
		c.Query += " LOWER(" + val + ")"
	}
	return c
}