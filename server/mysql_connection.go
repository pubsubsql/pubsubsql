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
	dbConn    *sql.DB
	address   string
	lastError string
}

func newMysqlConnection() *mysqlConnection {
	return &mysqlConnection{
		dbConn:    nil,
		address:   "",
		lastError: "",
	}
}

func (this *mysqlConnection) hasError() bool {
	return "" != this.lastError
}

func (this *mysqlConnection) hasNoError() bool {
	return !this.hasError()
}

func (this *mysqlConnection) getLastError() string {
	return this.lastError
}

func (this *mysqlConnection) isConnected() bool {
	if nil == this.dbConn {
		return false
	} else {
		err := this.dbConn.Ping()
		if nil != err {
			this.lastError = err.Error()
			this.dbConn.Close()
			this.dbConn = nil
			logError(this.lastError)
			return false
		}
		return true
	}
}

func (this *mysqlConnection) isDisconnected() bool {
	return !this.isConnected()
}

func (this *mysqlConnection) disconnect() {
	this.lastError = ""
	if this.isConnected() {
		err := this.dbConn.Close()
		if nil != err {
			this.lastError = err.Error()
			logError(this.lastError)
		}
		this.dbConn = nil
	}
}

func (this *mysqlConnection) connect(address string) bool {
	this.lastError = ""
	if this.isDisconnected() {
		this.address = address

		var err error
		this.dbConn, err = sql.Open("mysql", this.address)
		if nil != err {
			this.lastError = err.Error()
			this.dbConn = nil
			logError(this.lastError)
			return false
		}
	}
	return this.isConnected()
}

func (this *mysqlConnection) findTables() []string {
	this.lastError = ""
	tables := make([]string, 0)
	if this.isDisconnected() {
		this.lastError = "not connected to mysql"
		logError(this.lastError)
		return tables
	}
	rows, err := this.dbConn.Query("show tables")
	if nil != err {
		this.lastError = err.Error()
		logError(this.lastError)
		return tables
	}
	tableName := ""
	for rows.Next() {
		err := rows.Scan(&tableName)
		if nil != err {
			this.lastError = err.Error()
			logError(this.lastError)
			return tables
		}
		tables = append(tables, tableName)
	}
	err = rows.Err()
	if nil != err {
		this.lastError = err.Error()
		logError(this.lastError)
		return tables
	}
	//
	return tables
}

/*
create table t (c int)
create trigger t_t after insert on t for each row insert into log values (1);
*/
func (this *mysqlConnection) subscribe(tableName string) {
	this.lastError = ""
	if this.isDisconnected() {
		this.lastError = "not connected to mysql"
		logError(this.lastError)
		return
	}
	_, err := this.dbConn.Exec("create table t (c int)")
	if nil != err {
		this.lastError = err.Error()
		logError(this.lastError)
		return
	}
	_, err = this.dbConn.Exec("create trigger t_t after insert on t for each row insert into log values (1)")
	if nil != err {
		this.lastError = err.Error()
		logError(this.lastError)
		return
	}
}
