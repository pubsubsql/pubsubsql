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

type mysqlConnectRequest struct {
	cmdRequest
}

type mysqlDisconnectRequest struct {
	cmdRequest
}

type mysqlSubscribeRequest struct {
	cmdRequest
}

type mysqlUnsubscribeRequest struct {
	cmdRequest
}

// mysql connect
func (this *parser) parseMysqlConnect() request {
	tok := this.tokens.Produce()
	if tok.typ != tokenTypeEOF {
		return this.parseError("mysql connect: unexpected extra token")
	}
	return new(mysqlConnectRequest)
}

// mysql disconnect
func (this *parser) parseMysqlDisconnect() request {
	tok := this.tokens.Produce()
	if tok.typ != tokenTypeEOF {
		return this.parseError("mysql disconnect: unexpected extra token")
	}
	return new(mysqlDisconnectRequest)
}

// mysql subscribe
func (this *parser) parseMysqlSubscribe() request {
	tok := this.tokens.Produce()
	if tok.typ != tokenTypeEOF {
		return this.parseError("mysql subscribe: unexpected extra token")
	}
	return new(mysqlSubscribeRequest)
}

// mysql unsubscribe
func (this *parser) parseMysqlUnsubscribe() request {
	tok := this.tokens.Produce()
	if tok.typ != tokenTypeEOF {
		return this.parseError("mysql unsubscribe: unexpected extra token")
	}
	return new(mysqlUnsubscribeRequest)
}

// mysql
func (this *parser) parseSqlMysql() request {
	tok := this.tokens.Produce()
	switch tok.typ {
	case tokenTypeSqlConnect:
		return this.parseMysqlConnect()
	case tokenTypeSqlDisconnect:
		return this.parseMysqlDisconnect()
	case tokenTypeSqlSubscribe:
		return this.parseMysqlSubscribe()
	case tokenTypeSqlUnsubscribe:
		return this.parseMysqlUnsubscribe()
	}
	return this.parseError("invalid mysql request")
}
