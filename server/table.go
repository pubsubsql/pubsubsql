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

import "strconv"

const (
	tableCOLUMNS int = 10
	tableRECORDS     = 5000
)

// column  
type column struct {
	name    string
	ordinal int
}

// table  
type table struct {
	name     string
	colMap   map[string]*column
	colSlice []*column
	records  []*record
	// 
}

// table factory 
func newTable(name string) *table {
	t := &table{
		name:     name,
		colMap:   make(map[string]*column),
		colSlice: make([]*column, 0, tableCOLUMNS),
		records:  make([]*record, 0, tableRECORDS),
	}
	t.addColumn("id")
	return t
}

// addColumn adds column and returns column ordinal
func (t *table) addColumn(name string) int {
	ordinal := len(t.colSlice)
	col := &column{
		name:    name,
		ordinal: ordinal,
	}
	t.colMap[name] = col
	t.colSlice = append(t.colSlice, col)
	return ordinal
}

// getAddColumn tries to retrieve existing column  or adds it if does not exist
func (t *table) getAddColumn(name string) int {
	col, ok := t.colMap[name]
	if ok {
		return col.ordinal
	}
	return t.addColumn(name)
}

// getColumn retrieves existing column
func (t *table) getColumn(name string) *column {
	col, ok := t.colMap[name]
	if ok {
		return col
	}
	return nil
}

// getColumnCount returns total number of defined columns in the table
func (t *table) getColumnCount() int {
	return len(t.colSlice)
}

// addRecord adds new record to the table and returns newly added record
func (t *table) addRecord() *record {
	l := len(t.records)
	r := newRecord(len(t.colSlice), strconv.Itoa(l))
	addRecordToSlice(&t.records, r)
	return r
}

// addRecord
func addRecordToSlice(records *[]*record, r *record) {
	//check if records slice needs to grow by 10%
	l := len(*records)
	if cap(*records) == len(*records) {
		temp := *records
		*records = make([]*record, l, l+(l/10))
		copy(*records, temp)
	}
	*records = append(*records, r)
}

// getRecords returns record by id
func (t *table) getRecord(id int) *record {
	if len(t.records) > id {
		return t.records[id]
	}
	return nil
}

// getRecordCount returns total number of records in the table
func (t *table) getRecordCount() int {
	return len(t.records)
}

// sqlInsert proceses sql insert request and returns response
func (t *table) sqlInsert(req *sqlInsertRequest) response {
	rec := t.addRecord()
	for _, colVal := range req.colVals {
		rec.setValue(t.getAddColumn(colVal.col), colVal.val)
	}
	res := sqlInsertResponse{id: rec.getId()}
	return &res
}

// sqlSelect processes sql select request and returns response
func (t *table) sqlSelect(req *sqlSelectRequest) response {
	if req.filter.col != "" {
		return newErrorResponse("filters are not supported ")
	}
	// select * no filter
	var rows int
	rows = len(t.records)
	res := sqlSelectResponse{
		columns: t.colSlice,
		records: make([]*record, 0, rows),
	}
	for _, source := range t.records {
		res.copyRecordData(source)
	}
	return &res
}
