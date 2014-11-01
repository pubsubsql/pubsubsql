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

import "testing"

// MYSQL CONNECT xyz123
func TestMysqlConnect(t *testing.T) {
	consumer := chanTokenConsumer{channel: make(chan *token)}
	go lex("mysql connect xyz123", &consumer)
	expected := []token{
		{tokenTypeCmdMysql, "mysql"},
		{tokenTypeCmdConnect, "connect"},
		{tokenTypeSqlValue, "xyz123"},
		{tokenTypeEOF, ""}}

	validateTokens(t, expected, consumer.channel)
}

// MYSQL DISCONNECT
func TestMysqlDisconnect(t *testing.T) {
	consumer := chanTokenConsumer{channel: make(chan *token)}
	go lex("mysql disconnect", &consumer)
	expected := []token{
		{tokenTypeCmdMysql, "mysql"},
		{tokenTypeCmdDisconnect, "disconnect"},
		{tokenTypeEOF, ""}}

	validateTokens(t, expected, consumer.channel)
}

// MYSQL STATUS
func TestMysqlStatus(t *testing.T) {
	consumer := chanTokenConsumer{channel: make(chan *token)}
	go lex("mysql status", &consumer)
	expected := []token{
		{tokenTypeCmdMysql, "mysql"},
		{tokenTypeCmdStatus, "status"},
		{tokenTypeEOF, ""}}

	validateTokens(t, expected, consumer.channel)
}

// MYSQL TABLES
func TestMysqlTables(t *testing.T) {
	consumer := chanTokenConsumer{channel: make(chan *token)}
	go lex("mysql tables", &consumer)
	expected := []token{
		{tokenTypeCmdMysql, "mysql"},
		{tokenTypeCmdTables, "tables"},
		{tokenTypeEOF, ""}}

	validateTokens(t, expected, consumer.channel)
}

// MYSQL UNSUBSCRIBE
func TestMysqlUnsubscribe(t *testing.T) {
	consumer := chanTokenConsumer{channel: make(chan *token)}
	go lex("mysql unsubscribe from stocks", &consumer)
	expected := []token{
		{tokenTypeCmdMysql, "mysql"},
		{tokenTypeSqlUnsubscribe, "unsubscribe"},
		{tokenTypeSqlFrom, "from"},
		{tokenTypeSqlTable, "stocks"},
		{tokenTypeEOF, ""}}

	validateTokens(t, expected, consumer.channel)
}

// MYSQL SUBSCRIBE
func TestMysqlSubscribe(t *testing.T) {
	consumer := chanTokenConsumer{channel: make(chan *token)}
	go lex("mysql subscribe * from stocks", &consumer)
	expected := []token{
		{tokenTypeCmdMysql, "mysql"},
		{tokenTypeSqlSubscribe, "subscribe"},
		{tokenTypeSqlStar, "*"},
		{tokenTypeSqlFrom, "from"},
		{tokenTypeSqlTable, "stocks"},
		{tokenTypeEOF, ""}}

	validateTokens(t, expected, consumer.channel)
}
