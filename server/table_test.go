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
import "strconv"

func validateTableRecordsCount(t *testing.T, tbl *table, expected int) {
	val := tbl.getRecordCount()
	if val != expected {
		t.Errorf("records count do not match expected:%d but got:%d", expected, val)
	}
}

func validateSqlInsertResponseId(t *testing.T, res response, expected string) {
	switch res.(type) {
	case *sqlInsertResponse:
		x := res.(*sqlInsertResponse)
		if x.id != expected {
			t.Errorf("table insert error: expected id:%s but got:%s", expected, x.id)
		}
	default:
		t.Errorf("table insert error: invalid response type expected sqlInsertResponse")
	}
}

func validateOkResponse(t *testing.T, res response) {
	switch res.(type) {
	case *okResponse:

	default:
		t.Errorf("invalid response type expected okResponse")
	}
}

func validateErrorResponse(t *testing.T, res response) {
	switch res.(type) {
	case *errorResponse:

	default:
		t.Errorf("invalid response type expected errorResponse")
	}
}

func TestTable1(t *testing.T) {
	tbl := newTable("table1")
	tbl.getAddColumn("col1")
	r, _ := tbl.prepareRecord()
	tbl.addNewRecord(r)
	validateTableRecordsCount(t, tbl, 1)
	validateRecordValuesCount(t, r, 2)
	validateRecordValue(t, r, 0, "0")
	//
	r = tbl.getRecord(0)
	validateTableRecordsCount(t, tbl, 1)
	validateRecordValuesCount(t, r, 2)
	validateRecordValue(t, r, 0, "0")
}

func TestTable2(t *testing.T) {
	tbl := newTable("table1")
	tbl.getAddColumn("col1")
	tbl.getAddColumn("col2")
	tbl.getAddColumn("col3")
	tbl.getAddColumn("col3")
	col1 := tbl.getColumn("col1").ordinal
	col2 := tbl.getColumn("col2").ordinal
	col3 := tbl.getColumn("col3").ordinal
	//
	r, _ := tbl.prepareRecord()
	tbl.addNewRecord(r)
	validateTableRecordsCount(t, tbl, 1)
	validateRecordValuesCount(t, r, 4)
	validateRecordValue(t, r, 0, "0")
	r = tbl.getRecord(0)
	validateTableRecordsCount(t, tbl, 1)
	validateRecordValuesCount(t, r, 4)
	validateRecordValue(t, r, 0, "0")
	//	
	r, _ = tbl.prepareRecord()
	tbl.addNewRecord(r)
	validateTableRecordsCount(t, tbl, 2)
	validateRecordValuesCount(t, r, 4)
	validateRecordValue(t, r, 0, "1")
	validateRecordValue(t, r, 1, "")
	validateRecordValue(t, r, 2, "")
	validateRecordValue(t, r, 3, "")
	r.setValue(col1, "val1")
	r.setValue(col2, "val2")
	r.setValue(col3, "val3")
	validateRecordValue(t, r, 0, "1")
	validateRecordValue(t, r, 1, "val1")
	validateRecordValue(t, r, 2, "val2")
	validateRecordValue(t, r, 3, "val3")
	r = tbl.getRecord(1)
	validateTableRecordsCount(t, tbl, 2)
	validateRecordValuesCount(t, r, 4)
	validateRecordValue(t, r, 0, "1")
	validateRecordValue(t, r, 1, "val1")
	validateRecordValue(t, r, 2, "val2")
	validateRecordValue(t, r, 3, "val3")
}

// INSERT

func insertHelper(t *table, sqlInsert string) response {
	pc := newTokens()
	lex(sqlInsert, pc)
	req := parse(pc).(*sqlInsertRequest)
	return t.sqlInsert(req)
}

func TestTableSqlInsert(t *testing.T) {
	tbl := newTable("stocks")
	//
	res := insertHelper(tbl, " insert into stocks (ticker, bid, ask) values (IBM, 12, 14.5645)")
	validateSqlInsertResponseId(t, res, "0")
	t.Log(res.String())
	//
	res = insertHelper(tbl, " insert into stocks (ticker, bid, ask) values (MSFT, 37, 38) ")
	validateSqlInsertResponseId(t, res, "1")
	t.Log(res.String())
}

func BenchmarkTableSqlInser(b *testing.B) {
	tbl := newTable("stocks")
	for i := 0; i < b.N; i++ {
		insertHelper(tbl, " insert into stocks (ticker, bid, ask) values (IBM, 12, 14.5645) ")
	}
}

// SELECT

func selectHelper(t *table, sqlSelect string) response {
	pc := newTokens()
	lex(sqlSelect, pc)
	req := parse(pc).(*sqlSelectRequest)
	return t.sqlSelect(req)
}

func validateSqlSelect(t *testing.T, res response, rows int, cols int) {
	switch res.(type) {
	case *sqlSelectResponse:
		x := res.(*sqlSelectResponse)
		if len(x.columns) != cols {
			t.Errorf("table select error: expected column count:%d but got:%d", cols, len(x.columns))
		}
		if len(x.records) != rows {
			t.Errorf("table select error: expected rows count:%d but got:%d", rows, len(x.records))
		}
	default:
		t.Errorf("table select error: invalid response type expected sqlSelectResponse")
	}
}

func TestTableSqlSelect1(t *testing.T) {
	tbl := newTable("stocks")
	//
	insertHelper(tbl, " insert into stocks (ticker, bid, ask) values (IBM, 12, 14.5645) ")

	res := selectHelper(tbl, " select * from stocks ")
	validateSqlSelect(t, res, 1, 4)

	res = selectHelper(tbl, " select * from stocks where id = 0")
	validateSqlSelect(t, res, 1, 4)
	//	
	insertHelper(tbl, " insert into stocks (ticker, bid, ask, sector) values (IBM, 12, 14.5645, 'TECH') ")

	res = selectHelper(tbl, " select * from stocks ")
	validateSqlSelect(t, res, 2, 5)

	res = selectHelper(tbl, " select * from stocks where id = 1")
	validateSqlSelect(t, res, 1, 5)
}

// UPDATE

func updateHelper(t *table, sqlUpdate string) response {
	pc := newTokens()
	lex(sqlUpdate, pc)
	req := parse(pc).(*sqlUpdateRequest)
	return t.sqlUpdate(req)
}

func validateSqlUpdate(t *testing.T, res response, expected int) {
	switch res.(type) {
	case *sqlUpdateResponse:
		x := res.(*sqlUpdateResponse)
		if x.updated != expected {
			t.Errorf("table update error: expected update %d but got %d", expected, x.updated)
		}
	case *errorResponse:
		x := res.(*errorResponse)
		t.Errorf(x.msg)
	default:
		t.Errorf("table delete error: invalid response type expected sqlDeleteResponse")
	}
}

func TestTableSqlUpdate(t *testing.T) {
	tbl := newTable("stocks")
	// 1 record
	res := insertHelper(tbl, " insert into stocks (ticker, bid, ask, sector) values (IBM, 12, 14.5645, sec1) ")
	validateSqlInsertResponseId(t, res, "0")
	res = updateHelper(tbl, " update stocks set ticker = 'IBM', bid = 12, ask = 456.34")
	validateSqlUpdate(t, res, 1)
	// 3 records
	res = insertHelper(tbl, " insert into stocks (ticker, bid, ask, sector) values (MSFT, 12, 14.5645, sec1) ")
	validateSqlInsertResponseId(t, res, "1")
	res = insertHelper(tbl, " insert into stocks (ticker, bid, ask, sector) values (ORCL, 12, 14.5645, sec2) ")
	validateSqlInsertResponseId(t, res, "2")
	res = insertHelper(tbl, " insert into stocks (ticker, bid, ask, sector) values (C, 12, 14.5645, sec2) ")
	validateSqlInsertResponseId(t, res, "3")
	//
	res = updateHelper(tbl, " update stocks set bid = 12 ")
	validateSqlUpdate(t, res, 4)
	// create key for ticker
	res = keyHelper(tbl, "key stocks ticker")
	validateOkResponse(t, res)
	// update by key
	res = updateHelper(tbl, " update stocks set bid = 13 where ticker = IBM ")
	validateSqlUpdate(t, res, 1)
	res = updateHelper(tbl, " update stocks set bid = 13 where ticker = C ")
	validateSqlUpdate(t, res, 1)
	// update key by key
	res = updateHelper(tbl, " update stocks set ticker = 'JPM'  where ticker = IBM ")
	validateSqlUpdate(t, res, 1)
	res = selectHelper(tbl, " select * from stocks where ticker = JPM ")
	validateSqlSelect(t, res, 1, 5)
	//res = selectHelper(tbl, " select * from stocks where ticker = IBM ")
	//validateSqlSelect(t, res, 0, 5)
	// create tag for sector 
	res = tagHelper(tbl, "tag stocks sector")
	validateOkResponse(t, res)
	// update by sector
	res = selectHelper(tbl, " select * from stocks where sector = sec1 ")
	validateSqlSelect(t, res, 2, 5)
	res = updateHelper(tbl, " update stocks set bid = 13 where sector = sec1 ")
	validateSqlUpdate(t, res, 2)
	res = selectHelper(tbl, " select * from stocks where sector = sec1 ")
	validateSqlSelect(t, res, 2, 5)
	// update sector by sector
	res = updateHelper(tbl, " update stocks set sector = sec3 where sector = sec1 ")
	validateSqlUpdate(t, res, 2)
	res = selectHelper(tbl, " select * from stocks where sector = sec1 ")
	validateSqlSelect(t, res, 0, 5)
	res = selectHelper(tbl, " select * from stocks where sector = sec3 ")
	validateSqlSelect(t, res, 2, 5)

}

// DELETE

func deleteHelper(t *table, sqlDelete string) response {
	pc := newTokens()
	lex(sqlDelete, pc)
	req := parse(pc).(*sqlDeleteRequest)
	return t.sqlDelete(req)
}

func validateSqlDelete(t *testing.T, res response, expected int) {
	switch res.(type) {
	case *sqlDeleteResponse:
		x := res.(*sqlDeleteResponse)
		if x.deleted != expected {
			t.Errorf("table delete error: expected deleted %d but got %d", expected, x.deleted)
		}
	case *errorResponse:
		x := res.(*errorResponse)
		t.Errorf(x.msg)
	default:
		t.Errorf("table delete error: invalid response type expected sqlDeleteResponse")
	}
}

func TestTableSqlDelete(t *testing.T) {
	tbl := newTable("stocks")
	// 1 record
	res := insertHelper(tbl, " insert into stocks (ticker, bid, ask) values (IBM, 12, 14.5645) ")
	validateSqlInsertResponseId(t, res, "0")
	res = deleteHelper(tbl, " delete from stocks ")
	validateSqlDelete(t, res, 1)
	res = selectHelper(tbl, " select * from stocks ")
	validateSqlSelect(t, res, 0, 4)
	// 3 records
	res = insertHelper(tbl, " insert into stocks (ticker, bid, ask) values (IBM, 12, 14.5645) ")
	validateSqlInsertResponseId(t, res, "1")
	res = insertHelper(tbl, " insert into stocks (ticker, bid, ask) values (IBM, 12, 14.5645) ")
	validateSqlInsertResponseId(t, res, "2")
	res = insertHelper(tbl, " insert into stocks (ticker, bid, ask) values (IBM, 12, 14.5645) ")
	validateSqlInsertResponseId(t, res, "3")
	res = deleteHelper(tbl, " delete from stocks ")
	validateSqlDelete(t, res, 3)
	res = selectHelper(tbl, " select * from stocks ")
	validateSqlSelect(t, res, 0, 4)
}

// KEY

func keyHelper(t *table, sqlKey string) response {
	pc := newTokens()
	lex(sqlKey, pc)
	req := parse(pc).(*sqlKeyRequest)
	return t.sqlKey(req)
}

func TestTableSqlKey(t *testing.T) {
	tbl := newTable("stocks")
	// define key ticker
	res := keyHelper(tbl, "key stocks ticker")
	validateOkResponse(t, res)
	// insert record 
	res = insertHelper(tbl, " insert into stocks (ticker, bid, ask) values (IBM, 12, 14.5645) ")
	validateSqlInsertResponseId(t, res, "0")
	res = selectHelper(tbl, " select * from stocks where ticker = IBM")
	validateSqlSelect(t, res, 1, 4)
	// now define key for new column 
	res = keyHelper(tbl, "key stocks sector")
	validateErrorResponse(t, res)
	// should fail due to duplicate key 
	res = insertHelper(tbl, " insert into stocks (ticker, bid, ask) values (IBM, 12, 14.5645) ")
	validateErrorResponse(t, res)
	// now create another record with valid sector
	res = insertHelper(tbl, " insert into stocks (ticker, sector, bid, ask) values (MSFT, sec1, 12, 14.5645) ")
	validateSqlInsertResponseId(t, res, "1")
	res = selectHelper(tbl, " select * from stocks where ticker = IBM")
	validateSqlSelect(t, res, 1, 5)
	res = selectHelper(tbl, " select * from stocks where ticker = MSFT")
	validateSqlSelect(t, res, 1, 5)
	// test update duplicate key
	res = updateHelper(tbl, " update stocks set ticker = 'MSFT' where ticker = IBM")
	validateErrorResponse(t, res)
	// now sector is now unique empty string for IBM and sec1 for MSFT
	res = keyHelper(tbl, "key stocks sector")
	validateOkResponse(t, res)
	res = selectHelper(tbl, " select * from stocks where sector = ''")
	validateSqlSelect(t, res, 1, 5)
	res = selectHelper(tbl, " select * from stocks where sector = sec1")
	validateSqlSelect(t, res, 1, 5)
	// try to define existing key
	res = keyHelper(tbl, "key stocks ticker")
	validateErrorResponse(t, res)
	res = keyHelper(tbl, "key stocks sector")
	validateErrorResponse(t, res)
	// try to insert with duplicate key
	res = insertHelper(tbl, " insert into stocks (ticker, sector, bid, ask) values (ORCL, sec1, 12, 14.5645) ")
	validateErrorResponse(t, res)
	// try to insert with duplicate key and new column which should not be created
	l := tbl.getColumnCount()
	res = insertHelper(tbl, " insert into stocks (col1, col2, ticker, sector, bid, ask) values (col1, col2, ORCL, sec1, 12, 14.5645) ")
	validateErrorResponse(t, res)
	if l != tbl.getColumnCount() {
		t.Errorf("insert failed after duplicate keys rollback failed")
	}

	res = selectHelper(tbl, " select * from stocks ")
	validateSqlSelect(t, res, 2, 5)
	// delete by key
	res = deleteHelper(tbl, " delete from stocks where ticker = 'IBM'")
	validateSqlDelete(t, res, 1)
	res = selectHelper(tbl, " select * from stocks where ticker = IBM")
	validateSqlSelect(t, res, 0, 5)
	res = selectHelper(tbl, " select * from stocks ")
	validateSqlSelect(t, res, 1, 5)

	res = deleteHelper(tbl, " delete from stocks where ticker = NA")
	res = selectHelper(tbl, " select * from stocks ")
	validateSqlSelect(t, res, 1, 5)
	// delete by sec 
	res = selectHelper(tbl, " select * from stocks where ticker = MSFT")
	validateSqlSelect(t, res, 1, 5)
	res = deleteHelper(tbl, " delete from stocks where sector = 'sec1'")
	validateSqlDelete(t, res, 1)
	res = selectHelper(tbl, " select * from stocks where ticker = MSFT")
	validateSqlSelect(t, res, 0, 5)
	res = selectHelper(tbl, " select * from stocks ")
	validateSqlSelect(t, res, 0, 5)
}

// TAG

func tagHelper(t *table, sqlTag string) response {
	pc := newTokens()
	lex(sqlTag, pc)
	req := parse(pc).(*sqlTagRequest)
	return t.sqlTag(req)
}

func TestTableSqlTag(t *testing.T) {
	tbl := newTable("stocks")
	// tag ticker
	res := tagHelper(tbl, "tag stocks ticker")
	validateOkResponse(t, res)
	// insert records
	res = insertHelper(tbl, " insert into stocks (ticker, bid, ask) values (IBM, 12, 14.5645) ")
	validateSqlInsertResponseId(t, res, "0")
	res = selectHelper(tbl, " select * from stocks where ticker = IBM")
	validateSqlSelect(t, res, 1, 4)

	res = insertHelper(tbl, " insert into stocks (ticker, bid, ask) values (IBM, 12, 14.5645) ")
	validateSqlInsertResponseId(t, res, "1")
	res = selectHelper(tbl, " select * from stocks where ticker = IBM")
	validateSqlSelect(t, res, 2, 4)

	res = insertHelper(tbl, " insert into stocks (ticker, bid, ask) values (MSFT, 12, 14.5645) ")
	validateSqlInsertResponseId(t, res, "2")
	res = selectHelper(tbl, " select * from stocks where ticker = MSFT")
	validateSqlSelect(t, res, 1, 4)

	res = insertHelper(tbl, " insert into stocks (ticker, bid, ask) values (IBM, 12, 14.5645) ")
	validateSqlInsertResponseId(t, res, "3")
	res = selectHelper(tbl, " select * from stocks where ticker = IBM")
	validateSqlSelect(t, res, 3, 4)

	if tbl.getTagedColumnValuesCount("ticker", "IBM") != 3 {
		t.Errorf("invalid taged column values")
	}
	if tbl.getTagedColumnValuesCount("ticker", "MSFT") != 1 {
		t.Errorf("invalid taged column values")
	}
	if 4 != tbl.getColumnCount() {
		t.Errorf("tag failed: expected 4 columns but got %d", tbl.getColumnCount())
	}
	// tag sector
	res = tagHelper(tbl, "tag stocks sector")
	validateOkResponse(t, res)
	if 5 != tbl.getColumnCount() {
		t.Errorf("tag failed: expected 5 columns but got %d", tbl.getColumnCount())
	}
	if tbl.getTagedColumnValuesCount("sector", "") != 4 {
		t.Errorf("invalid taged column values")
	}
	//	
	res = insertHelper(tbl, " insert into stocks (ticker, sector, bid, ask) values (IBM, 'TECH', 12, 14.5645) ")
	validateSqlInsertResponseId(t, res, "4")
	if tbl.getTagedColumnValuesCount("sector", "") != 4 {
		t.Errorf("invalid taged column values")
	}
	if tbl.getTagedColumnValuesCount("sector", "TECH") != 1 {
		t.Errorf("invalid taged column values")
	}
	//
	res = deleteHelper(tbl, " delete from stocks where ticker = 'IBM'")
	validateSqlDelete(t, res, 4)
	res = selectHelper(tbl, " select * from stocks where ticker = IBM")
	validateSqlSelect(t, res, 0, 5)
	res = selectHelper(tbl, " select * from stocks where ticker = MSFT")
	validateSqlSelect(t, res, 1, 5)
	if tbl.getTagedColumnValuesCount("sector", "TECH") != 0 {
		t.Errorf("invalid taged column values")
	}
	if tbl.getTagedColumnValuesCount("sector", "") != 1 {
		t.Errorf("invalid taged column values")
	}
	res = deleteHelper(tbl, " delete from stocks where sector = ''")
	validateSqlDelete(t, res, 1)
	if tbl.getTagedColumnValuesCount("sector", "") != 0 {
		t.Errorf("invalid taged column values")
	}
	//	
	res = selectHelper(tbl, " select * from stocks ")
	validateSqlSelect(t, res, 0, 5)
}

// SUBSCRIBE

func subscribeHelper(t *table, sqlSubscribe string) (response, *responseSender) {
	sender := newResponseSenderStub(0)
	pc := newTokens()
	lex(sqlSubscribe, pc)
	req := parse(pc).(*sqlSubscribeRequest)
	req.sender = sender
	t.sqlSubscribe(req)
	return sender.tryRecv(), sender
}

func validateSqlSubscribeResponse(t *testing.T, res response) *sqlSubscribeResponse {
	if res == nil {
		t.Errorf("table subscribe error: invalid response nil, expected sqlSubscribeResponse")
	}
	switch res.(type) {
	case *sqlSubscribeResponse:
		x := res.(*sqlSubscribeResponse)
		return x
	case *errorResponse:
		x := res.(*errorResponse)
		t.Errorf(x.msg)
	default:
		t.Errorf("table subscribe error: invalid response type expected sqlSubscribeResponse")
	}
	return nil
}

func validateSqlActionAddResponse(t *testing.T, sender *responseSender, pubsubid uint64, records int) {
	res := sender.tryRecv()
	if res == nil {
		t.Errorf("table subscribe error: invalid response nil, expected sqlActionAddResponse")
	}
	switch res.(type) {
	case *sqlActionAddResponse:
		x := res.(*sqlActionAddResponse)
		if x.pubsubid != pubsubid {
			t.Errorf("invalid sqlActionAddResponse pubsubid expected:%d but got:%d", pubsubid, x.pubsubid)
		}
		l := len(x.sqlSelectResponse.records)
		if l != records {
			t.Errorf("invalid sqlActionAddResponse records expected:%d but got:%d", records, l)
		}
	case *errorResponse:
		x := res.(*errorResponse)
		t.Errorf(x.msg)
	default:
		t.Errorf("table subscribe error: invalid response type expected sqlActionAddResponse")
	}
}

func validateNoResponse(t *testing.T, sender *responseSender) {
	res := sender.tryRecv()
	if res != nil {
		t.Errorf("table subscribe error: invalid response, expected nil")
	}
}

func TestTableSqlSubscribe1(t *testing.T) {
	tbl := newTable("stocks")
	// key ticker
	res := keyHelper(tbl, "key stocks ticker")
	validateOkResponse(t, res)
	// tag sector
	res = tagHelper(tbl, "tag stocks sector")
	validateOkResponse(t, res)
	// insert records
	res = insertHelper(tbl, " insert into stocks (ticker, bid, ask, sector) values (IBM, 12, 14.56, TECH) ")
	validateSqlInsertResponseId(t, res, "0")
	// SUBSCRIBE
	// subscribe to table
	var sender *responseSender
	res, sender = subscribeHelper(tbl, "subscribe * from stocks ")
	sub := validateSqlSubscribeResponse(t, res)
	validateSqlActionAddResponse(t, sender, sub.pubsubid, 1)
	// subscribe to existing key
	res, sender = subscribeHelper(tbl, "subscribe * from stocks where ticker = IBM")
	sub = validateSqlSubscribeResponse(t, res)
	validateSqlActionAddResponse(t, sender, sub.pubsubid, 1)
	// subscribe to existing tag
	res, sender = subscribeHelper(tbl, "subscribe * from stocks where sector = TECH")
	sub = validateSqlSubscribeResponse(t, res)
	validateSqlActionAddResponse(t, sender, sub.pubsubid, 1)
	// subscribe to id		
	res, sender = subscribeHelper(tbl, "subscribe * from stocks where id = 0")
	sub = validateSqlSubscribeResponse(t, res)
	validateSqlActionAddResponse(t, sender, sub.pubsubid, 1)
	// subscribe to non existing valid key
	res, sender = subscribeHelper(tbl, "subscribe * from stocks where ticker = MSFT")
	validateSqlSubscribeResponse(t, res)
	validateNoResponse(t, sender)
	// subscribe to non existing valid tag
	res, sender = subscribeHelper(tbl, "subscribe * from stocks where sector = FIN")
	validateSqlSubscribeResponse(t, res)
	validateNoResponse(t, sender)
	// subscribe to non existing invalid key/tag
	res, sender = subscribeHelper(tbl, "subscribe * from stocks where invalidkey = somevalue")
	validateErrorResponse(t, res)
	validateNoResponse(t, sender)
	// subscribe to non existing id
	res, sender = subscribeHelper(tbl, "subscribe * from stocks where id = 1")
	validateErrorResponse(t, res)
	validateNoResponse(t, sender)
}

// UNSUBSCRIBE

func unsubscribeHelper(t *table, sqlUnsubscribe string, connectionId uint64) response {
	pc := newTokens()
	lex(sqlUnsubscribe, pc)
	req := parse(pc).(*sqlUnsubscribeRequest)
	req.connectionId = connectionId
	return t.sqlUnsubscribe(req)
}

func validateSqlUnsubscribeResponse(t *testing.T, res response, unsubscribed int) {
	switch res.(type) {
	case *sqlUnsubscribeResponse:
		x := res.(*sqlUnsubscribeResponse)
		if x.unsubscribed != unsubscribed {
			t.Errorf("invalid sqlUnsubscribeResponse unsubscribed expected:%d but got:%d", unsubscribed, x.unsubscribed)
		}
	case *errorResponse:
		x := res.(*errorResponse)
		t.Errorf(x.msg)
	default:
		t.Errorf("table unsubscribe error: invalid response type expected sqlUnsubscribeResponse")
	}
}

func TestTableSqlUnSubscribe1(t *testing.T) {
	tbl := newTable("stocks")
	// key ticker
	res := keyHelper(tbl, "key stocks ticker")
	validateOkResponse(t, res)
	// tag sector
	res = tagHelper(tbl, "tag stocks sector")
	validateOkResponse(t, res)
	// insert records
	res = insertHelper(tbl, " insert into stocks (ticker, bid, ask, sector) values (IBM, 12, 14.56, TECH) ")
	validateSqlInsertResponseId(t, res, "0")
	// SUBSCRIBE
	// subscribe to table
	var sender *responseSender
	res, sender = subscribeHelper(tbl, "subscribe * from stocks ")
	sub := validateSqlSubscribeResponse(t, res)
	connectionId := sender.connectionId
	pubsubid := strconv.FormatUint(sub.pubsubid, 10)
	// subscribe to existing key
	res, sender = subscribeHelper(tbl, "subscribe * from stocks where ticker = IBM")
	// subscribe to existing tag
	res, sender = subscribeHelper(tbl, "subscribe * from stocks where sector = TECH")
	// subscribe to id		
	res, sender = subscribeHelper(tbl, "subscribe * from stocks where id = 0")
	// subscribe to non existing valid key
	res, sender = subscribeHelper(tbl, "subscribe * from stocks where ticker = MSFT")
	// subscribe to non existing valid tag
	res, sender = subscribeHelper(tbl, "subscribe * from stocks where sector = FIN")

	// unsubscribe	
	res = unsubscribeHelper(tbl, "unsubscribe from stocks where pubsubid = "+pubsubid, connectionId)
	validateSqlUnsubscribeResponse(t, res, 1)
	res = unsubscribeHelper(tbl, "unsubscribe from stocks ", connectionId)
	validateSqlUnsubscribeResponse(t, res, 5)
}
