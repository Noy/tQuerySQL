// Created by Noy Hillel
// https://github.com/noy
// You are free to use/modify this software in your own projects
package tsql

import (
	"database/sql"
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

func (c *Client) handleBasicQuerySubmission(queryType, queryStr string) *Client {
	if c.Query == "" {
		handleTSQLError(queryType)
	} else {
		c.Query += queryStr
	}
	return c
}

func (c *Client) From(table string) *Client {
	return c.handleBasicQuerySubmission("FROM", " FROM " + table)
}

func (c *Client) Where(value string) *Client {
	return c.handleBasicQuerySubmission("WHERE", " WHERE " + value)
}

func (c *Client) Equals(value string, quotes bool) *Client {
	if quotes {
		return c.handleBasicQuerySubmission("EQUALS", " = '"+value+"'")
	} else {
		return c.handleBasicQuerySubmission("EQUALS", " = "+value+"")
	}
}

func (c *Client) Like(value string, percentSign bool) *Client {
	if percentSign {
		return c.handleBasicQuerySubmission("LIKE", " LIKE '%"+value+"%'")
	} else {
		return c.handleBasicQuerySubmission("LIKE", " LIKE "+value)
	}
}

func (c *Client) NotLike(value string, percentSign bool) *Client {
	if percentSign {
		return c.handleBasicQuerySubmission("NOT LIKE", " NOT LIKE '%"+value+"%'")
	} else {
		return c.handleBasicQuerySubmission("NOT LIKE", " NOT LIKE '"+value+"'")
	}
}

func (c *Client) NotEqual(value int) *Client {
	return c.handleBasicQuerySubmission("NOT EQUAL (!=)", " != " + String(value))
}

func String(n int) string {
	buf := [11]byte{}
	pos := len(buf)
	i := int64(n)
	signed := i < 0
	if signed {
		i = -i
	}
	for {
		pos--
		buf[pos], i = '0'+byte(i%10), i/10
		if i == 0 {
			if signed {
				pos--
				buf[pos] = '-'
			}
			return string(buf[pos:])
		}
	}
}


func (c *Client) LT(value string, quotes bool) *Client {
	if quotes {
		return c.handleBasicQuerySubmission("< (less than)", " < '" + value + "'")
	} else {
		return c.handleBasicQuerySubmission("< (less than)", " < " + value)
	}
}

func (c *Client) GT(value string, quotes bool) *Client {
	if quotes {
		return c.handleBasicQuerySubmission("> (more than)", " > '" + value + "'")
	} else {
		return c.handleBasicQuerySubmission("> (more than)", " > " + value)
	}
}

func (c *Client) GTE(value string, quotes bool) *Client {
	if quotes {
		return c.handleBasicQuerySubmission(">= (GTE)", " >= '" + value + "'")
	} else {
		return c.handleBasicQuerySubmission(">= (GTE)", " >= " + value)
	}
}

func (c *Client) LTE(value string, quotes bool) *Client {
	if quotes {
		return c.handleBasicQuerySubmission("<= (LTE)", " <= '" + value + "'")
	} else {
		return c.handleBasicQuerySubmission("<= (LTE)", " <= " + value)
	}
}

func (c *Client) And(value string) *Client {
	return c.handleBasicQuerySubmission("AND", " AND " + value)
}

func (c *Client) Or(value string) *Client {
	return c.handleBasicQuerySubmission("OR", " OR " + value)
}

func (c *Client) Union() *Client {
	return c.handleBasicQuerySubmission("UNION", " UNION ")
}

func (c *Client) UnionAll() *Client {
	return c.handleBasicQuerySubmission("UNION ALL", " UNION ALL ")
}

func (c *Client) IsNotNull() *Client {
	return c.handleBasicQuerySubmission("IS NOT NULL", " IS NOT NULL ")
}

func (c *Client) handleVarArgs(queryType string, vals []string) *Client {
	if c.Query == "" {
		handleTSQLError(queryType)
	}
	c.Query += queryType
	for _, v := range vals {
		c.Query += v + ","
	}
	ql := len(c.Query)
	if ql > 0 && c.Query[ql-1] == ',' {
		c.Query = c.Query[:ql-1]
	}
	if queryType == " VALUES(" {
		c.Query += ")"
	}
	return c
}

func (c *Client) GroupBy(vals ...string) *Client {
	return c.handleVarArgs(" GROUP BY ", vals)
}

func (c *Client) Values(vals ...string) *Client {
	return c.handleVarArgs(" VALUES(", vals)
}

func (c *Client) Limit(amount string) *Client {
	return c.handleBasicQuerySubmission("LIMIT",  " LIMIT " + amount)
}

func (c *Client) Offset(number string) *Client {
	return c.handleBasicQuerySubmission("OFFSET",  " OFFSET " + number)
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
	c.Query += "INSERT INTO " + table
	return c
}

func (c *Client) Update(table string) *Client {
	return c.handleBasicQuerySubmission("UPDATE", "UPDATE " + table)
}

func (c *Client) Set(val string) *Client {
	return c.handleBasicQuerySubmission("SET", " SET " + val)
}

func (c *Client) Join(table string) *Client {
	return c.handleBasicQuerySubmission("JOIN", " JOIN " + table)
}

func (c *Client) LeftOuterJoin(table string) *Client {
	return c.handleBasicQuerySubmission("LEFT OUTER JOIN", " LEFT OUTER JOIN " + table)
}

func (c *Client) LeftInnerJoin(table string) *Client {
	return c.handleBasicQuerySubmission("LEFT INNER JOIN", " LEFT INNER JOIN " + table)
}

func (c *Client) RightOuterJoin(table string) *Client {
	return c.handleBasicQuerySubmission("RIGHT OUTER JOIN", " RIGHT OUTER JOIN " + table)
}

func (c *Client) RightInnerJoin(table string) *Client {
	return c.handleBasicQuerySubmission("RIGHT INNER JOIN", " RIGHT INNER JOIN " + table)
}

func (c *Client) On(val string) *Client {
	return c.handleBasicQuerySubmission("ON", " ON " + val)
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
	if quotes {
		return c.handleBasicQuerySubmission("LOWER", " LOWER('" + val + "')")
	} else {
		return c.handleBasicQuerySubmission("LOWER", " LOWER(" + val + ")")
	}
}
