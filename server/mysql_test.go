/* Copyright (C) 2013 CompleteDB LLC.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with PubSubSQL.  If not, see <http://www.gnu.org/licenses/>.
 */

package server

import (
	"testing"
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

/*
https://code.google.com/p/go-wiki/wiki/SQLDrivers
https://github.com/go-sql-driver/mysql
//
create database pubsubsql;
create user pubsubsql identified by 'pubsubsql';
GRANT ALL PRIVILEGES ON *.* TO 'pubsubsql'@'%' WITH GRANT OPTION;
 */
func TestMysqlConnection(t *testing.T) {
	conn, err := sql.Open("mysql", "pubsubsql:pubsubsql@/pubsubsql")
	if nil != err {
		fmt.Println(err)
		t.Error("failed to open mysql connection:", err)
		return
	}
	defer conn.Close()
}
