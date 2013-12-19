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

package pubsubsql

import "testing"
import "time"

func TestDataServiceRunAndStop(t *testing.T) {
	stoper := NewStoper()
	dataSrv := newDataService(stoper)
	go dataSrv.run()
	if !stoper.Stop(3 * time.Second) {
		t.Errorf("stoper.Stop() expected true but got false")
	}
}

func sqlHelper(sql string, sender *responseSender) *requestItem {
	pc := newTokens()
	lex(sql, pc)
	req := parse(pc).(request)
	return &requestItem{
		req:    req,
		sender: sender,
	}
}

func TestDataService(t *testing.T) {
	stoper := NewStoper()
	dataSrv := newDataService(stoper)
	go dataSrv.run()
	sender := newResponseSenderStub(1)
	// insert
	dataSrv.accept(sqlHelper("insert into stocks (ticker, bid, ask, sector) values (IBM, 123, 124, TECH) ", sender))
	res := sender.recv()
	validateSqlInsertResponseId(t, res, "0")
	// select
	dataSrv.accept(sqlHelper(" select * from stocks ", sender))
	res = sender.recv()
	validateSqlSelect(t, res, 1, 5)
	// key 
	dataSrv.accept(sqlHelper(" key stocks ticker ", sender))
	res = sender.recv()
	validateOkResponse(t, res)
	// tag 
	dataSrv.accept(sqlHelper(" tag stocks sector ", sender))
	res = sender.recv()
	validateOkResponse(t, res)
	// subscribe	
	dataSrv.accept(sqlHelper(" subscribe * from stocks sector = TECH ", sender))
	res = sender.recv()
	validateSqlSubscribeResponse(t, res)
	res = sender.recv() // action add
	// update
	dataSrv.accept(sqlHelper(" update stocks set bid = 140 where ticker = IBM ", sender))
	res = sender.recv() // first is action update
	res = sender.recv()
	validateSqlUpdate(t, res, 1)
	// delete
	dataSrv.accept(sqlHelper(" delete from stocks where ticker = IBM ", sender))
	res = sender.recv() // first is action delete
	res = sender.recv()
	validateSqlDelete(t, res, 1)
	// unsubscribe
	dataSrv.accept(sqlHelper(" unsubscribe from stocks where pubsubid = 1 ", sender))
	res = sender.recv() // first is action delete
	validateSqlUnsubscribe(t, res, 1)

	stoper.Stop(time.Millisecond * 1000)
}
