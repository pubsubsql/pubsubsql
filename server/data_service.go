/* Copyright (C) 2013 CompleteD LLC.
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

type requestItem struct {
	req    request
	sender *responseSender
}

// dataService prer-processes sqlRequests and channels them to approptiate tables for further proccessging
type dataService struct {
	requests   chan *requestItem
	stoper     *Stoper
	tables     map[string]*table
	bufferSize int
}

// dataService factory
func newDataService(bufferSize int, stoper *Stoper) *dataService {
	return &dataService{
		requests:   make(chan *requestItem, bufferSize),
		stoper:     stoper,
		tables:     make(map[string]*table),
		bufferSize: bufferSize,
	}
}

// accepts the request
func (d *dataService) accept(r *requestItem) {
	select {
	case d.requests <- r:
	case <-d.stoper.GetChan():
	}
}

// runs dataService event loop
func (d *dataService) run() {
	s := d.stoper
	s.Enter()
	defer s.Leave()
	for {
		select {
		case item := <-d.requests:
			if s.IsStoping() {
				debug("data service exited isStoping")
				return
			}
			d.onSqlRequest(item)
		case <-s.GetChan():
			debug("data service exited stoped")
			return
		}
	}
}

func (d *dataService) onSqlRequest(item *requestItem) {
	tableName := item.req.getTableName()
	tbl := d.tables[tableName]
	if tbl == nil {
		// auto create table and enter table event loop
		tbl = newTable(tableName)
		d.tables[tableName] = tbl
		tbl.stoper = d.stoper
		tbl.requests = make(chan *requestItem, d.bufferSize)
		go tbl.run()
	}
	// forward sql request to the table
	tbl.requests <- item
}
