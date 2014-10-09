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
func validateMysqlConnect(t *testing.T, req request) {
	switch req.(type) {
	case *errorRequest:
		e := req.(*errorRequest)
		t.Errorf("parse error: " + e.err)

	case *mysqlConnectRequest:

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

func TestParseMysqlConnect(t *testing.T) {
	pc := newTokens()
	lex(" mysql connect xyz123 ", pc)
	req := parse(pc)
	validateMysqlConnect(t, req)
}

func TestParseMysqlDisconnect(t *testing.T) {
	pc := newTokens()
	lex(" mysql disconnect ", pc)
	req := parse(pc)
	validateMysqlDisconnect(t, req)
}

func TestParseMysqlSubscribe(t *testing.T) {
	pc := newTokens()
	lex(" mysql subscribe ", pc)
	req := parse(pc)
	validateMysqlSubscribe(t, req)
}

func TestParseMysqlUnsubscribe(t *testing.T) {
	pc := newTokens()
	lex(" mysql unsubscribe from stocks ", pc)
	req := parse(pc)
	validateMysqlUnsubscribe(t, req)
}
