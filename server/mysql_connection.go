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
type mysqlConnection struct {
	dbConn *sql.DB
	lastError error
}

func newMysqlConnection() *mysqlConnection {
	return &mysqlConnection {
		dbConn: nil,
		lastError: nil,
	}
}

func (this *mysqlConnection) getLastError() error {
	return this.lastError
}

func (this *mysqlConnection) isConnected() bool {
	return nil != this.dbConn
}

func (this *mysqlConnection) isDisconnected() bool {
	return ! this.isConnected()
}

func (this *mysqlConnection) disconnect() {
	if this.isConnected() {
		this.dbConn.Close()
		this.dbConn = nil
	}
}

func (this *mysqlConnection) connect() {
	this.dbConn, this.lastError = sql.Open("mysql", "pubsubsql:pubsubsql@/pubsubsql")
	if nil != this.lastError {
		this.dbConn = nil
	}
	if this.isConnected() {
		this.lastError = this.dbConn.Ping()
		if nil != this.lastError {
			this.dbConn = nil
		}
	}
}
