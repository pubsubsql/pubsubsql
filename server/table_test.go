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

func validateSqlSelectResponse(t *testing.T, res response, rows int, cols int) {
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
	r, _ := tbl.newRecord()
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
	r, _ := tbl.newRecord()
	tbl.addNewRecord(r)
	validateTableRecordsCount(t, tbl, 1)
	validateRecordValuesCount(t, r, 4)
	validateRecordValue(t, r, 0, "0")
	r = tbl.getRecord(0)
	validateTableRecordsCount(t, tbl, 1)
	validateRecordValuesCount(t, r, 4)
	validateRecordValue(t, r, 0, "0")
	//	
	r, _ = tbl.newRecord()
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

func TestTableSqlSelect(t *testing.T) {
	tbl := newTable("stocks")
	//
	insertHelper(tbl, " insert into stocks (ticker, bid, ask) values (IBM, 12, 14.5645) ")
	res := selectHelper(tbl, " select * from stocks ")
	validateSqlSelectResponse(t, res, 1, 4)
	//	
	insertHelper(tbl, " insert into stocks (ticker, bid, ask, sector) values (IBM, 12, 14.5645, 'TECH') ")
	res = selectHelper(tbl, " select * from stocks ")
	validateSqlSelectResponse(t, res, 2, 5)
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
	// define key
	res := keyHelper(tbl, "key stocks ticker")
	validateOkResponse(t, res)
	// insert record
	res = insertHelper(tbl, " insert into stocks (ticker, bid, ask) values (IBM, 12, 14.5645) ")
	validateSqlInsertResponseId(t, res, "0")
	// now define key for new column 
	res = keyHelper(tbl, "key stocks sector")
	validateErrorResponse(t, res)
	// should fail due to duplicate key 
	res = insertHelper(tbl, " insert into stocks (ticker, bid, ask) values (IBM, 12, 14.5645) ")
	validateErrorResponse(t, res)
	// now create another record with valid secotor
	res = insertHelper(tbl, " insert into stocks (ticker, sector, bid, ask) values (MSFT, sec1, 12, 14.5645) ")
	validateSqlInsertResponseId(t, res, "1")
	// now sector is now unique empty string for IBM and sec1 for MSFT
	res = keyHelper(tbl, "key stocks sector")
	validateOkResponse(t, res)
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
	res = insertHelper(tbl, " insert into stocks (ticker, bid, ask) values (IBM, 12, 14.5645) ")
	validateSqlInsertResponseId(t, res, "1")
	res = insertHelper(tbl, " insert into stocks (ticker, bid, ask) values (MSFT, 12, 14.5645) ")
	validateSqlInsertResponseId(t, res, "2")
	res = insertHelper(tbl, " insert into stocks (ticker, bid, ask) values (IBM, 12, 14.5645) ")
	validateSqlInsertResponseId(t, res, "3")
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

}
