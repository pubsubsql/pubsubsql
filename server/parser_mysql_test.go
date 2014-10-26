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

// MYSQL CONNECT
func validateMysqlConnect(t *testing.T, a request, y *mysqlConnectRequest) {
	switch a.(type) {
	case *errorRequest:
		e := a.(*errorRequest)
		t.Errorf("parse error: " + e.err)

	case *mysqlConnectRequest:
		x := a.(*mysqlConnectRequest)
		// connectionAddress
		if x.address != y.address {
			t.Errorf("parse error: connectionAddress do not match")
		}

	default:
		t.Errorf("parse error: invalid request type expected mysqlConnectRequest")
	}
}

// MYSQL DISCONNECT
func validateMysqlDisconnect(t *testing.T, req request) {
	switch req.(type) {
	case *errorRequest:
		e := req.(*errorRequest)
		t.Errorf("parse error: " + e.err)

	case *mysqlDisconnectRequest:

	default:
		t.Errorf("parse error: invalid request type expected mysqlDisconnectRequest")
	}
}

// MYSQL STATUS
func validateMysqlStatus(t *testing.T, req request) {
	switch req.(type) {
	case *errorRequest:
		e := req.(*errorRequest)
		t.Errorf("parse error: " + e.err)

	case *mysqlStatusRequest:

	default:
		t.Errorf("parse error: invalid request type expected mysqlStatusRequest")
	}
}

// MYSQL SUBSCRIBE
func validateMysqlSubscribe(t *testing.T, req request) {
	switch req.(type) {
	case *errorRequest:
		e := req.(*errorRequest)
		t.Errorf("parse error: " + e.err)

	case *mysqlSubscribeRequest:

	default:
		t.Errorf("parse error: invalid request type expected mysqlSubscribeRequest")
	}
}

// MYSQL UNSUBSCRIBE
func validateMysqlUnsubscribe(t *testing.T, req request) {
	switch req.(type) {
	case *errorRequest:
		e := req.(*errorRequest)
		t.Errorf("parse error: " + e.err)

	case *mysqlUnsubscribeRequest:

	default:
		t.Errorf("parse error: invalid request type expected mysqlUnsubscribeRequest")
	}
}

// MYSQL TABLES
func validateMysqlTables(t *testing.T, req request) {
	switch req.(type) {
	case *errorRequest:
		e := req.(*errorRequest)
		t.Errorf("parse error: " + e.err)

	case *mysqlTablesRequest:

	default:
		t.Errorf("parse error: invalid request type expected mysqlTablesRequest")
	}
}

func TestParseMysqlConnect(t *testing.T) {
	pc := newTokens()
	lex(" mysql connect xyz123 ", pc)
	x := parse(pc)
	var y mysqlConnectRequest
	y.address = "xyz123"
	validateMysqlConnect(t, x, &y)
}

func TestParseMysqlConnectQuoted(t *testing.T) {
	pc := newTokens()
	lex(" mysql connect xyz123 ", pc)
	x := parse(pc)
	var y mysqlConnectRequest
	y.address = "xyz123"
	validateMysqlConnect(t, x, &y)
}

func TestParseMysqlDisconnect(t *testing.T) {
	pc := newTokens()
	lex(" mysql disconnect ", pc)
	req := parse(pc)
	validateMysqlDisconnect(t, req)
}

func TestParseMysqlStatus(t *testing.T) {
	pc := newTokens()
	lex(" mysql status ", pc)
	req := parse(pc)
	validateMysqlStatus(t, req)
}

func TestParseMysqlSubscribe(t *testing.T) {
	pc := newTokens()
	lex(" mysql subscribe * from stocks ", pc)
	req := parse(pc)
	validateMysqlSubscribe(t, req)
}

func TestParseMysqlUnsubscribe(t *testing.T) {
	pc := newTokens()
	lex(" mysql unsubscribe from stocks ", pc)
	req := parse(pc)
	validateMysqlUnsubscribe(t, req)
}

func TestParseMysqlTables(t *testing.T) {
	pc := newTokens()
	lex(" mysql tables ", pc)
	req := parse(pc)
	validateMysqlTables(t, req)
}
