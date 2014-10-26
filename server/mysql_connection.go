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
http://go-database-sql.org/
//
create database pubsubsql;
create user pubsubsql identified by 'pubsubsql';
GRANT ALL PRIVILEGES ON *.* TO 'pubsubsql'@'%' WITH GRANT OPTION;
create table table_a(id int);
create table table_b(id int);
create table table_c(id int);
show tables
 */
type mysqlConnection struct {
	dbConn *sql.DB
	address string
	lastError error
}

func newMysqlConnection() *mysqlConnection {
	return &mysqlConnection {
		dbConn: nil,
		address: "",
		lastError: nil,
	}
}

func (this *mysqlConnection) hasError() bool {
	return nil != this.lastError
}

func (this *mysqlConnection) hasNoError() bool {
	return ! this.hasError()
}

func (this *mysqlConnection) getLastError() string {
	if (this.hasError()) {
		return this.lastError.Error()
	} else {
		return ""
	}
}

func (this *mysqlConnection) isConnected() bool {
	if nil == this.dbConn {
		return false
	} else {
		this.lastError = this.dbConn.Ping()
		if nil != this.lastError {
			this.dbConn.Close()
			this.dbConn = nil
			return false
		}
		return true
	}
}

func (this *mysqlConnection) isDisconnected() bool {
	return ! this.isConnected()
}

func (this *mysqlConnection) disconnect() {
	this.lastError = nil
	if this.isConnected() {
		this.lastError = this.dbConn.Close()
		this.dbConn = nil
	}
}

func (this *mysqlConnection) connect(address string) bool {
	this.lastError = nil
	if this.isDisconnected() {
		this.address = address
		// "pubsubsql:pubsubsql@/pubsubsql"
		this.dbConn, this.lastError = sql.Open("mysql", this.address)
		if nil != this.lastError {
			this.dbConn = nil
			return false
		}
	}
	return this.isConnected();
}

func (this *mysqlConnection) findTables() []string {
	this.lastError = nil
	tables := make([]string, 0)
	if (this.isDisconnected()) {
		return tables
	}
	rows, err := this.dbConn.Query("show tables")
	this.lastError = err
	if nil != this.lastError {
		logError(this.lastError)
		return tables
	}
	tableName := ""
	for rows.Next() {
		this.lastError = rows.Scan(&tableName)
		if  nil != this.lastError {
			logError(this.lastError)
			return tables
		}
		tables = append(tables, tableName)
	}
	this.lastError = rows.Err()
	if nil != this.lastError {
		logError(this.lastError)
		return tables
	}
	//
	return tables
}
