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
	key     map[string]int
	tags    map[string]*tag
}

func (c *column) hasKey() bool {
	return c.key != nil
}

func (c *column) hasTags() bool {
	return c.tags != nil
}

func (c *column) keyContainsValue(k string) bool {
	_, contains := c.key[k]
	return contains
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

// getColumnCount returns number of columns
func (t *table) getColumnCount() int {
	l := len(t.colSlice)
	if l != len(t.colMap) {
		panic("Something bad happened column slice and map do not match")
	}
	return l
}

// addColumn adds column and returns column ordinal
func (t *table) addColumn(name string) *column {
	ordinal := len(t.colSlice)
	col := &column{
		name:    name,
		ordinal: ordinal,
	}
	t.colMap[name] = col
	t.colSlice = append(t.colSlice, col)
	return col
}

// getAddColumn tries to retrieve existing column  or adds it if does not exist
// returns true when new column was added
func (t *table) getAddColumn(name string) (*column, bool) {
	col, columnExists := t.colMap[name]
	if columnExists {
		return col, false
	}
	return t.addColumn(name), true
}

// getColumn retrieves existing column
func (t *table) getColumn(name string) *column {
	col, ok := t.colMap[name]
	if ok {
		return col
	}
	return nil
}

// newRecord creates new record but does not add it to the table
func (t *table) newRecord() (*record, int) {
	l := len(t.records)
	r := newRecord(len(t.colSlice), strconv.Itoa(l))
	return r, l
}

// adNewRecord add newly created record to the table
func (t *table) addNewRecord(r *record) {
	addRecordToSlice(&t.records, r)
}

// addRecordToSlice generic helper function that adds record to the slice and
// automatically expands the slice
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

// sqlKey defines unique index in the table
func (t *table) sqlKey(req *sqlKeyRequest) response {
	// key is already defined for this column
	col := t.getColumn(req.column)
	if col != nil && (col.hasKey() || col.hasTags()) {
		return newErrorResponse("key or tag already defined for column:" + req.column)
	}
	// new column on existing records
	if col == nil && len(t.records) > 0 {
		return newErrorResponse("can not define key for non existant column due to possible duplicates")
	}
	key := make(map[string]int, cap(t.records))
	// new column no records
	if col == nil {
		col, _ = t.getAddColumn(req.column)
		col.key = key
	} else {
		// index all records and check if there are duplicates
		key = make(map[string]int, cap(t.records))
		for idx, rec := range t.records {
			val := rec.getValue(col.ordinal)
			if col.keyContainsValue(val) {
				return newErrorResponse("can not define key due to possible duplicates in existing records")
			}
			key[val] = idx
		}
	}
	col.key = key
	return newOkResponse()
}

func addValueToTags(tags map[string]*tag, val string, idx int) {
	head := tags[val]
	tg := addTag(head, idx)
	if head == nil {
		tags[val] = tg
	}
}

func (t *table) sqlTag(req *sqlTagRequest) response {
	// tag is already defined for this column
	col := t.getColumn(req.column)
	if col != nil && (col.hasKey() || col.hasTags()) {
		return newErrorResponse("key or tag already defined for column:" + req.column)
	}
	col, _ = t.getAddColumn(req.column)
	// tag existing values
	tags := make(map[string]*tag, cap(t.records))
	for idx, rec := range t.records {
		val := rec.getValue(col.ordinal)
		addValueToTags(col.tags, val, idx)
	}
	//
	col.tags = tags
	return newOkResponse()
}

// removeColumns removes columns starting at particular ordinal
func (t *table) removeColumns(ordinal int) {
	if len(t.colSlice) <= ordinal {
		return
	}
	tail := t.colSlice[ordinal:]
	for _, col := range tail {
		delete(t.colMap, col.name)
	}
	t.colSlice = t.colSlice[:ordinal]
}

// sqlInsert proceses sql insert request and returns response
func (t *table) sqlInsert(req *sqlInsertRequest) response {
	rec, id := t.newRecord()
	// validate unique keys constrain
	cols := make([]*column, len(req.colVals))
	originalColLen := len(t.colSlice)
	for idx, colVal := range req.colVals {
		col, _ := t.getAddColumn(colVal.col)
		if col.hasKey() && col.keyContainsValue(colVal.val) {
			//remove created columns
			t.removeColumns(originalColLen)
			return newErrorResponse("insert failed due to duplicate column key:" + colVal.col + " value:" + colVal.val)
		}
		cols[idx] = col
	}
	// ready to insert	
	for idx, colVal := range req.colVals {
		col := cols[idx]
		// update key
		if col.hasKey() {
			col.key[colVal.val] = id
		} else if col.hasTags() {
			addValueToTags(col.tags, colVal.val, id)
		}
		rec.setValue(col.ordinal, colVal.val)
	}
	t.addNewRecord(rec)
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
