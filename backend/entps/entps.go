package entps

import (
	"database/sql"
	"database/sql/driver"

	"modernc.org/sqlite"
)

type sqlite3Driver struct {
	*sqlite.Driver
}

type sqlite3DriverConn interface {
	//nolint:inamedparam
	Exec(string, []driver.Value) (driver.Result, error)
}

//nolint:nonamedreturns
func (d sqlite3Driver) Open(name string) (conn driver.Conn, err error) {
	conn, err = d.Driver.Open(name)
	if err != nil {
		return
	}
	_, err = conn.(sqlite3DriverConn).Exec("PRAGMA foreign_keys = ON; PRAGMA busy_timeout = 10000;", nil)
	if err != nil {
		_ = conn.Close()
	}
	return
}

func init() {
	sql.Register("sqlite3", sqlite3Driver{Driver: &sqlite.Driver{}})
}
