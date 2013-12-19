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
	pubsub       pubsub
	//
	subscriptions mapSubscriptionByConnection
	subid         uint64
	//
	requests chan *requestItem
	stoper   *Stoper
}

// table factory
func newTable(name string) *table {
	t := &table{
		name:          name,
		colMap:        make(map[string]*column),
		colSlice:      make([]*column, 0, config.TABLE_COLUMNS_CAPACITY),
		records:       make([]*record, 0, config.TABLE_RECORDS_CAPACITY),
		tagedColumns:  make([]*column, 0, config.TABLE_COLUMNS_CAPACITY),
		subscriptions: make(mapSubscriptionByConnection),
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
	t.records[rec.id()] = nil
}

// Looks up record by id.
// Returns record slice with max one element.
func (t *table) getRecordById(val string) []*record {
	idx, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		return nil
	}
	if idx < 0 || int64(len(t.records)) <= idx {
		return nil
	}
	records := make([]*record, 1, 1)
	records[0] = t.records[idx]
	return records
}

// Validates sql filter
// Returns errorResponse on error
func (t *table) validateSqlFilter(filter sqlFilter) (response, *column) {
	var col *column
	if len(filter.col) > 0 {
		col = t.getColumn(filter.col)
		if col == nil {
			return newErrorResponse("invalid column: " + filter.col), nil
		}
	}
	if col != nil && col.typ == columnTypeNormal {
		return newErrorResponse("can not use non indexed column " + filter.col + " as valid filter"), nil
	}
	return nil, col
}

// Retrieves records based by column value
func (t *table) getRecordsByValue(val string, col *column) []*record {
	if col == nil {
		// all
		return t.records
	}
	switch col.typ {
	case columnTypeKey:
		return t.getRecordsByTag(val, col)
	case columnTypeTag:
		return t.getRecordsByTag(val, col)
	case columnTypeId:
		return t.getRecordById(val)
	}
	return nil
}

// Retrieves records based on the supplied filter
func (t *table) getRecordsBySqlFilter(filter sqlFilter) ([]*record, response) {
	e, col := t.validateSqlFilter(filter)
	if e != nil {
		return nil, e
	}
	return t.getRecordsByValue(filter.val, col), nil
}

// Looks up records by tag.
func (t *table) getRecordsByTag(val string, col *column) []*record {
	// we need to optimize allocations
	// perhaps its possible to know in advance how manny records
	// will be returned
	records := make([]*record, 0, config.TABLE_GET_RECORDS_BY_TAG_CAPACITY)
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

type pubsubRA struct {
	removed []*pubsub
	added   map[*pubsub]int
}

func newPubsubRA() *pubsubRA {
	return &pubsubRA{
		removed: make([]*pubsub, 0, 3),
		added:   make(map[*pubsub]int),
	}
}

func getIfHasData(ra *pubsubRA) *pubsubRA {
	if ra != nil && (len(ra.removed) > 0 || len(ra.added) > 0) {
		return ra
	}
	return nil
}

func hasWhatToRemove(ra *pubsubRA) bool {
	return ra != nil && len(ra.removed) > 0
}

func hasWhatToAdd(ra *pubsubRA) bool {
	return ra != nil && len(ra.added) > 0
}

func (ra *pubsubRA) toBeRemoved(pubsub *pubsub) {
	if pubsub != nil {
		ra.removed = append(ra.removed, pubsub)
	}
}

func (ra *pubsubRA) toBeAdded(pubsub *pubsub) {
	if pubsub != nil {
		ra.added[pubsub] = 1
	}
}

func (t *table) updateRecordKeyTag(col *column, val string, rec *record, id int, ra **pubsubRA) {
	r := t.deleteTag(rec, col)
	rec.setValue(col.ordinal, val)
	a := t.tagValue(col, id, rec)
	// updated with the same value ignore this case
	if r == a {
		return
	}
	if *ra == nil {
		*ra = newPubsubRA()
	}
	ra.toBeRemoved(r)
	ra.toBeAdded(a)
}

// Updates record with new values, keys and tags.
func (t *table) updateRecord(cols []*column, colVals []*columnValue, rec *record, id int) *pubsubRA {
	var ra *pubsubRA
	for idx, colVal := range colVals {
		col := cols[idx]
		switch col.typ {
		case columnTypeKey:
			t.updateRecordKeyTag(col, colVal.val, rec, id, &ra)
		case columnTypeTag:
			t.updateRecordKeyTag(col, colVal.val, rec, id, &ra)
		case columnTypeNormal:
			rec.setValue(col.ordinal, colVal.val)
		}
	}
	return getIfHasData(ra)
}

// TAGS helper functions

// Add value to non unique indexed column.
func addValueToTags(col *column, val string, idx int) (*tag, *pubsub) {
	return col.tagmap.addTag(val, idx)
}

// Binds tag, pubsub and record.
func (t *table) tagValue(col *column, idx int, rec *record) *pubsub {
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
	return pubsub
}

// Deletes tag value for a particular record
func (t *table) deleteTag(rec *record, col *column) *pubsub {
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
	ret := lnk.pubsub
	lnk.clear()
	return ret
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
	t.onInsert(rec)
	return &res
}

// SELECT sql statement

func (t *table) copyRecordsToSqlSelectResponse(r *sqlSelectResponse, records []*record) {
	r.columns = t.colSlice
	r.records = make([]*record, 0, len(records))
	for _, rec := range records {
		if rec != nil {
			r.copyRecordData(rec)
		}
	}
}

func (t *table) copyRecordToSqlSelectResponse(r *sqlSelectResponse, rec *record) {
	r.columns = t.colSlice
	r.records = make([]*record, 0, 1)
	r.copyRecordData(rec)
}

// Processes sql select request.
// On success returns sqlSelectResponse.
func (t *table) sqlSelect(req *sqlSelectRequest) response {
	records, errResponse := t.getRecordsBySqlFilter(req.filter)
	if errResponse != nil {
		return errResponse
	}
	var r sqlSelectResponse
	t.copyRecordsToSqlSelectResponse(&r, records)
	return &r
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
	cols := make([]*column, len(req.colVals)+1)
	originalColLen := len(t.colSlice)
	cols[0] = t.colSlice[0]
	for idx, colVal := range req.colVals {
		col, _ := t.getAddColumn(colVal.col)
		if col.isKey() && col.keyContainsValue(colVal.val) {
			if onlyRecord == nil || onlyRecord != t.getRecordsByTag(colVal.val, col)[0] {
				//remove created columns
				t.removeColumns(originalColLen)
				return newErrorResponse("update failed due to duplicate column key:" + colVal.col + " value:" + colVal.val)
			}
		}
		cols[idx+1] = col
	}
	// all is valid ready to update
	updated := 0
	for _, rec := range records {
		if rec != nil {
			updated++
			ra := t.updateRecord(cols[1:], req.colVals, rec, int(rec.id()))
			if hasWhatToRemove(ra) {
				t.onRemove(ra.removed, rec)
			}
			var added *map[*pubsub]int
			if hasWhatToAdd(ra) {
				added = &ra.added
				t.onAdd(ra.added, rec)
			}
			t.onUpdate(cols, rec, added)
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
			t.onDelete(rec)
			t.deleteRecord(rec)
			rec.free()
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

// SUBSCRIBE sql statement

func (t *table) newSubscription(sender *responseSender) *subscription {
	t.subid++
	sub := newSubscription(sender, t.subid)
	t.subscriptions.add(sender.connectionId, sub)
	return sub
}

func (t *table) subscribeToTable(sender *responseSender) (*subscription, []*record) {
	sub := t.newSubscription(sender)
	t.pubsub.add(sub)
	sender.send(newSubscribeResponse(sub))
	return sub, t.records
}

func (t *table) subscribeToKeyOrTag(col *column, val string, sender *responseSender) (*subscription, []*record) {
	sub := t.newSubscription(sender)
	records := t.getRecordsByTag(val, col)
	col.tagmap.getAddTagItem(val).pubsub.add(sub)
	sender.send(newSubscribeResponse(sub))
	return sub, records
}

func (t *table) subscribeToId(id string, sender *responseSender) (*subscription, []*record) {
	records := t.getRecordById(id)
	if len(records) > 0 {
		sub := t.newSubscription(sender)
		records[0].addSubscription(sub)
		sender.send(newSubscribeResponse(sub))
		return sub, records
	}
	sender.send(newErrorResponse("id: " + id + " does not exist"))
	return nil, nil
}

func (t *table) subscribe(col *column, val string, sender *responseSender) (*subscription, []*record) {
	if col == nil {
		return t.subscribeToTable(sender)
	}
	switch col.typ {
	case columnTypeKey:
		return t.subscribeToKeyOrTag(col, val, sender)
	case columnTypeTag:
		return t.subscribeToKeyOrTag(col, val, sender)
	case columnTypeId:
		return t.subscribeToId(val, sender)
	}
	sender.send(newErrorResponse("Unexpected logical error"))
	return nil, nil
}

// Processes sql subscribe request.
// Does not return anything, responses are send directly to response sender.
func (t *table) sqlSubscribe(req *sqlSubscribeRequest) {
	// validate
	e, col := t.validateSqlFilter(req.filter)
	if e != nil {
		req.sender.send(e)
		return
	}
	// subscribe
	sub, records := t.subscribe(col, req.filter.val, req.sender)
	if sub != nil && len(records) > 0 {
		// publish initial action add
		t.publishActionAdd(sub, records)
	}
}

// PUBSUB helpers
type publishAction func(tbl *table, sub *subscription, rec *record) bool

func (t *table) visitSubscriptions(rec *record, publishActionFunc publishAction) {
	f := func(sub *subscription) bool {
		return publishActionFunc(t, sub, rec)
	}
	t.pubsub.visit(f)
	for _, lnk := range rec.links {
		if lnk.pubsub != nil {
			lnk.pubsub.visit(f)
		}
	}
}

func (t *table) publishActionAdd(sub *subscription, records []*record) bool {
	r := new(sqlActionAddResponse)
	r.pubsubid = sub.id
	t.copyRecordsToSqlSelectResponse(&r.sqlSelectResponse, records)
	return sub.sender.send(r)
}

func publishActionInsert(t *table, sub *subscription, rec *record) bool {
	r := new(sqlActionInsertResponse)
	r.pubsubid = sub.id
	t.copyRecordToSqlSelectResponse(&r.sqlSelectResponse, rec)
	return sub.sender.send(r)
}

func publishActionDelete(t *table, sub *subscription, rec *record) bool {
	r := &sqlActionDeleteResponse{
		id:       rec.idAsString(),
		pubsubid: sub.id,
	}
	return sub.sender.send(r)
}

func (t *table) onInsert(rec *record) {
	t.visitSubscriptions(rec, publishActionInsert)
}

func (t *table) onDelete(rec *record) {
	t.visitSubscriptions(rec, publishActionDelete)
}

func (t *table) onRemove(pubsubs []*pubsub, rec *record) {
	f := func(sub *subscription) bool {
		r := &sqlActionRemoveResponse{
			id:       rec.idAsString(),
			pubsubid: sub.id,
		}
		return sub.sender.send(r)
	}
	t.pubsub.visit(f)
	for _, pubsub := range pubsubs {
		pubsub.visit(f)
	}
}

func (t *table) onAdd(added map[*pubsub]int, rec *record) {
	f := func(sub *subscription) bool {
		r := new(sqlActionAddResponse)
		r.pubsubid = sub.id
		t.copyRecordToSqlSelectResponse(&r.sqlSelectResponse, rec)
		return sub.sender.send(r)
	}
	for pubsub, _ := range added {
		pubsub.visit(f)
	}
}

func (t *table) onUpdate(cols []*column, rec *record, added *map[*pubsub]int) {
	f := func(sub *subscription) bool {
		r := newSqlActionUpdateResponse(sub.id, cols, rec)
		return sub.sender.send(r)
	}
	t.pubsub.visit(f)
	for _, lnk := range rec.links {
		if lnk.pubsub != nil {
			// ignore updates for record that was just added
			if added != nil && (*added)[lnk.pubsub] != 0 {
				continue
			}
			lnk.pubsub.visit(f)
		}
	}
}

// UNSUBSCRIBE

// Processes sql unsubscribe request.
func (t *table) sqlUnsubscribe(req *sqlUnsubscribeRequest) response {
	// validate
	if len(req.filter.col) > 0 && req.filter.col != "pubsubid" {
		return newErrorResponse("Invalid filter expected pubsubid but got " + req.filter.col)
	}
	// unsubscribe by pubsubid for a given connection
	res := new(sqlUnsubscribeResponse)
	val := req.filter.val
	if len(val) > 0 {
		pubsubid, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return newErrorResponse("Failed to unsubscribe, pubsubid " + val + " is not valid")
		}
		if t.subscriptions.deactivate(req.connectionId, pubsubid) {
			res.unsubscribed = 1
		}
	} else {
		// unsubscribe all subscriptions for a given connection
		res.unsubscribed = t.subscriptions.deactivateAll(req.connectionId)
	}
	return res
}

// run

func (t *table) run() {
	//
	s := t.stoper
	s.Join()
	defer s.Leave()
	for {
		select {
		case item := <-t.requests:
			if s.Stoped() {
				debug("table exited isStoping")
				return
			}
			t.onSqlRequest(item.req, item.sender)
		case <-s.GetChan():
			debug("table exited stoped")
			return
		}
	}
}

func (t *table) onSqlRequest(r request, sender *responseSender) {
	switch r.(type) {
	case *sqlInsertRequest:
		t.onSqlInsert(r.(*sqlInsertRequest), sender)

	case *sqlSelectRequest:
		t.onSqlSelect(r.(*sqlSelectRequest), sender)

	case *sqlUpdateRequest:
		t.onSqlUpdate(r.(*sqlUpdateRequest), sender)

	case *sqlDeleteRequest:
		t.onSqlDelete(r.(*sqlDeleteRequest), sender)

	case *sqlSubscribeRequest:
		t.onSqlSubscribe(r.(*sqlSubscribeRequest), sender)

	case *sqlUnsubscribeRequest:
		t.onSqlUnsubscribe(r.(*sqlUnsubscribeRequest), sender)

	case *sqlKeyRequest:
		t.onSqlKey(r.(*sqlKeyRequest), sender)

	case *sqlTagRequest:
		t.onSqlTag(r.(*sqlTagRequest), sender)
	}
}

func (t *table) onSqlInsert(req *sqlInsertRequest, sender *responseSender) {
	sender.send(t.sqlInsert(req))
}

func (t *table) onSqlSelect(req *sqlSelectRequest, sender *responseSender) {
	sender.send(t.sqlSelect(req))
}

func (t *table) onSqlUpdate(req *sqlUpdateRequest, sender *responseSender) {
	sender.send(t.sqlUpdate(req))
}

func (t *table) onSqlDelete(req *sqlDeleteRequest, sender *responseSender) {
	sender.send(t.sqlDelete(req))
}

func (t *table) onSqlSubscribe(req *sqlSubscribeRequest, sender *responseSender) {
	req.sender = sender
	t.sqlSubscribe(req)
}

func (t *table) onSqlUnsubscribe(req *sqlUnsubscribeRequest, sender *responseSender) {
	req.connectionId = sender.connectionId
	sender.send(t.sqlUnsubscribe(req))
}

func (t *table) onSqlKey(req *sqlKeyRequest, sender *responseSender) {
	sender.send(t.sqlKey(req))
}

func (t *table) onSqlTag(req *sqlTagRequest, sender *responseSender) {
	sender.send(t.sqlTag(req))
}
