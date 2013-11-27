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

// Default values
const (
	tableCOLUMNS int = 10
	tableRECORDS     = 5000
)

// this function is purely for testing porposes
func (t *table) getTagedColumnValuesCount(col string, val string) int {
	c := t.getColumn(col)
	if c == nil || !c.isTag() {
		return 0
	}
	i := 0
	for tg := c.tagmap.getTag(val); tg != nil; tg = tg.next {
		i++
	}
	return i
}

// table  
type table struct {
	name         string
	colMap       map[string]*column
	colSlice     []*column
	records      []*record
	tagedColumns []*column
	pubsub       pubSub
}

// table factory 
func newTable(name string) *table {
	t := &table{
		name:         name,
		colMap:       make(map[string]*column),
		colSlice:     make([]*column, 0, tableCOLUMNS),
		records:      make([]*record, 0, tableRECORDS),
		tagedColumns: make([]*column, 0, tableCOLUMNS),
	}
	t.addColumn("id")
	return t
}

// COLUMNS functions

// Returns total number of columns. 
func (t *table) getColumnCount() int {
	l := len(t.colSlice)
	if l != len(t.colMap) {
		panic("Something bad happened column slice and map do not match")
	}
	return l
}

// Adds column and returns column added column.
func (t *table) addColumn(name string) *column {
	col := newColumn(name, len(t.colSlice))
	t.colMap[name] = col
	t.colSlice = append(t.colSlice, col)
	return col
}

// Tries to retrieve existing column or adds it if does not exist.
// Returns true when new column was added.
func (t *table) getAddColumn(name string) (*column, bool) {
	col, columnExists := t.colMap[name]
	if columnExists {
		return col, false
	}
	return t.addColumn(name), true
}

// Retrieves existing column
func (t *table) getColumn(name string) *column {
	col, ok := t.colMap[name]
	if ok {
		return col
	}
	return nil
}

// Deletes columns starting at particular ordinal.
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

// RECORDS functions

// Creates new record but does not add it to the table.
// Returns new record and to be record id
func (t *table) prepareRecord() (*record, int) {
	id := len(t.records)
	r := newRecord(len(t.colSlice), id)
	l := len(t.tagedColumns) + 1
	r.links = make([]link, l)
	return r, id
}

// adNewRecord add newly created record to the table
func (t *table) addNewRecord(r *record) {
	addRecordToSlice(&t.records, r)
}

// addRecordToSlice generic helper function that adds record to the slice and
// automatically expands the slice
func addRecordToSlice(records *[]*record, r *record) {
	//check if records slice needs to grow by third 
	l := len(*records)
	if cap(*records) == len(*records) {
		temp := *records
		*records = make([]*record, l, l+(l/3))
		copy(*records, temp)
	}
	*records = append(*records, r)
}

// Returns record by id
func (t *table) getRecord(id int) *record {
	if len(t.records) > id {
		return t.records[id]
	}
	return nil
}

// Returns total number of records in the table
func (t *table) getRecordCount() int {
	return len(t.records)
}

// Delete record from the table.
func (t *table) deleteRecord(rec *record) {
	// delete record tags
	for _, col := range t.tagedColumns {
		t.deleteTag(rec, col)
	}
	// delete record
	t.records[rec.idx()] = nil
}

// Looks up record by id.
// Returns record slice with max one element.
func (t *table) getRecordById(val string) []*record {
	idx, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		return nil
	}
	records := make([]*record, 1, 1)
	records[0] = t.records[idx]
	return records
}

// Retrieves records based on the supplied filter
func (t *table) getRecordsBySqlFilter(filter sqlFilter) ([]*record, response) {
	var col *column
	if len(filter.col) > 0 {
		col = t.getColumn(filter.col)
		if col == nil {
			return nil, newErrorResponse("invalid column: " + filter.col)
		}
	}
	if col == nil {
		// all
		return t.records, nil
	}
	switch col.typ {
	case columnTypeId:
		return t.getRecordById(filter.val), nil
	case columnTypeKey:
		return t.getRecordsByTag(filter.val, col), nil
	case columnTypeTag:
		return t.getRecordsByTag(filter.val, col), nil
	}
	return nil, newErrorResponse("can not use non indexed column " + filter.col + " as valid filter")
}

// Looks up records by tag. 
func (t *table) getRecordsByTag(val string, col *column) []*record {
	// we need to optimize allocations
	// perhaps its possible to know in advance how manny records
	// will be returned
	records := make([]*record, 0, 100)
	for tg := col.tagmap.getTag(val); tg != nil; tg = tg.next {
		records = append(records, t.records[tg.idx])
		l := len(records)
		if cap(records) == l {
			temp := records
			records = make([]*record, l, l+(l/3))
			copy(records, temp)
		}
	}
	return records
}

// Bind records values, keys and tags. 
func (t *table) bindRecord(cols []*column, colVals []*columnValue, rec *record, id int) {
	for idx, colVal := range colVals {
		col := cols[idx]
		rec.setValue(col.ordinal, colVal.val)
		// update key
		switch col.typ {
		case columnTypeKey:
			t.tagValue(col, id, rec)
		case columnTypeTag:
			t.tagValue(col, id, rec)
		}
	}
}

// Updates record with new values, keys and tags.
func (t *table) updateRecord(cols []*column, colVals []*columnValue, rec *record, id int) {
	for idx, colVal := range colVals {
		col := cols[idx]
		switch col.typ {
		case columnTypeKey:
			// delete previous key
			t.deleteTag(rec, col)
			rec.setValue(col.ordinal, colVal.val)
			t.tagValue(col, id, rec)
		case columnTypeTag:
			// delete previous tag
			t.deleteTag(rec, col)
			rec.setValue(col.ordinal, colVal.val)
			t.tagValue(col, id, rec)
		case columnTypeNormal:
			rec.setValue(col.ordinal, colVal.val)
		}
	}
}

// TAGS helper functions

// Add value to non unique indexed column.
func addValueToTags(col *column, val string, idx int) (*tag, *pubSub) {
	return col.tagmap.addTag(val, idx)
}

// Binds tag, pubsub and record.   
func (t *table) tagValue(col *column, idx int, rec *record) {
	val := rec.getValue(col.ordinal)
	tg, pubsub := addValueToTags(col, val, idx)
	lnk := link{
		tg:     tg,
		pubsub: pubsub,
	}
	if len(rec.links) <= col.tagIndex {
		rec.links = append(rec.links, lnk)
	} else {
		rec.links[col.tagIndex] = lnk
	}
}

// Deletes tag value for a particular record
func (t *table) deleteTag(rec *record, col *column) {
	lnk := &rec.links[col.tagIndex]
	if lnk.tg != nil {
		switch removeTag(lnk.tg) {
		case removeTagLast:
			col.tagmap.removeTag(rec.getValue(col.ordinal))
		case removeTagSlide:
			// we need to retag the slided record	
			slidedRecord := t.records[lnk.tg.idx]
			if slidedRecord != nil {
				slidedRecord.links[col.tagIndex].tg = lnk.tg
			}
		}
	}
	lnk.clear()
}

// INSERT sql statement

// Proceses sql insert request by inserting record in the table.
// On success returns sqlInsertResponse.
func (t *table) sqlInsert(req *sqlInsertRequest) response {
	rec, id := t.prepareRecord()
	// validate unique keys constrain
	cols := make([]*column, len(req.colVals))
	originalColLen := len(t.colSlice)
	for idx, colVal := range req.colVals {
		col, _ := t.getAddColumn(colVal.col)
		if col.isKey() && col.keyContainsValue(colVal.val) {
			//remove created columns
			t.removeColumns(originalColLen)
			return newErrorResponse("insert failed due to duplicate column key:" + colVal.col + " value:" + colVal.val)
		}
		cols[idx] = col
	}
	// ready to insert	
	t.bindRecord(cols, req.colVals, rec, id)
	t.addNewRecord(rec)
	res := sqlInsertResponse{id: rec.idAsString()}
	return &res
}

// SELECT sql statement

// Processes sql select request.
// On success returns sqlSelectResponse.
func (t *table) sqlSelect(req *sqlSelectRequest) response {
	records, errResponse := t.getRecordsBySqlFilter(req.filter)
	if errResponse != nil {
		return errResponse
	}
	res := sqlSelectResponse{
		columns: t.colSlice,
		records: make([]*record, 0, len(records)),
	}
	for _, rec := range records {
		if rec != nil {
			res.copyRecordData(rec)
		}
	}
	return &res
}

// UPDATE sql statement

// Processes sql update request.
// On success returns sqlUpdateResponse.
func (t *table) sqlUpdate(req *sqlUpdateRequest) response {
	records, errResponse := t.getRecordsBySqlFilter(req.filter)
	if errResponse != nil {
		return errResponse
	}
	res := &sqlUpdateResponse{updated: 0}
	var onlyRecord *record
	switch len(records) {
	case 0:
		return res
	case 1:
		onlyRecord = records[0]
	}
	// validate duplicate keys
	cols := make([]*column, len(req.colVals))
	originalColLen := len(t.colSlice)
	for idx, colVal := range req.colVals {
		col, _ := t.getAddColumn(colVal.col)
		if col.isKey() && col.keyContainsValue(colVal.val) {
			if onlyRecord == nil || onlyRecord != t.getRecordsByTag(colVal.val, col)[0] {
				//remove created columns
				t.removeColumns(originalColLen)
				return newErrorResponse("update failed due to duplicate column key:" + colVal.col + " value:" + colVal.val)
			}
		}
		cols[idx] = col
	}
	// all is valid ready to update
	updated := 0
	for _, rec := range records {
		if rec != nil {
			t.updateRecord(cols, req.colVals, rec, int(rec.idx()))
			updated++
		}
	}
	res.updated = updated
	return res
}

// DELETE sql statement

// Processes sql delete request.
// On success returns sqlDeleteResponse.
func (t *table) sqlDelete(req *sqlDeleteRequest) response {
	records, errResponse := t.getRecordsBySqlFilter(req.filter)
	if errResponse != nil {
		return errResponse
	}
	deleted := 0
	for _, rec := range records {
		if rec != nil {
			deleted++
			t.deleteRecord(rec)
		}
	}
	return &sqlDeleteResponse{deleted: deleted}
}

// Key sql statement

// Processes sql key request.
// On success returns sqlOkResponse.
func (t *table) sqlKey(req *sqlKeyRequest) response {
	// key is already defined for this column
	col := t.getColumn(req.column)
	if col != nil && col.isIndexed() {
		return newErrorResponse("key or tag already defined for column:" + req.column)
	}
	// new column on existing records
	if col == nil && len(t.records) > 0 {
		return newErrorResponse("can not define key for non existant column due to possible duplicates")
	}
	// new column no records
	if col != nil {
		unique := make(map[string]int, cap(t.records))
		// check if there are duplicates
		for idx, rec := range t.records {
			val := rec.getValue(col.ordinal)
			if _, contains := unique[val]; contains {
				return newErrorResponse("can not define key due to possible duplicates in existing records")
			}
			unique[val] = idx
		}
	}
	//
	t.tagOrKeyColumn(req.column, columnTypeKey)
	return newOkResponse()
}

// TAG sql statement

func (t *table) tagOrKeyColumn(c string, coltyp columnType) {
	col, _ := t.getAddColumn(c)
	t.tagedColumns = append(t.tagedColumns, col)
	col.makeTags(len(t.tagedColumns))
	col.typ = coltyp
	// tag existing values
	for idx, rec := range t.records {
		t.tagValue(col, idx, rec)
	}
}

// Processes sql tag request.
// On success returns sqlOkResponse.
func (t *table) sqlTag(req *sqlTagRequest) response {
	// tag is already defined for this column
	col := t.getColumn(req.column)
	if col != nil && col.isIndexed() {
		return newErrorResponse("key or tag already defined for column:" + req.column)
	}
	//
	t.tagOrKeyColumn(req.column, columnTypeTag)
	return newOkResponse()
}
