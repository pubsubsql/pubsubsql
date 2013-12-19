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
}

// dataService factory
func newDataService(stoper *Stoper) *dataService {
	return &dataService{
		requests:   make(chan *requestItem, config.CHAN_DATASERVICE_REQUESTS_BUFFER_SIZE),
		stoper:     stoper,
		tables:     make(map[string]*table),
	}
}

// accepts the request
func (this *dataService) accept(item *requestItem) {
	select {
	case this.requests <- item:
	case <-this.stoper.GetChan():
	}
}

// runs dataService event loop
func (this *dataService) run() {
	this.stoper.Join()
	defer this.stoper.Leave()
	for {
		select {
		case item := <-this.requests:
			if this.stoper.Stoped() {
				debug("data service exited due to stop event")
				return
			}
			this.onSqlRequest(item)
		case <-this.stoper.GetChan():
			debug("data service exited due to stop event")
			return
		}
	}
}

func (this *dataService) onSqlRequest(item *requestItem) {
	tableName := item.req.getTableName()
	tbl := this.tables[tableName]
	if tbl == nil {
		// auto create table and go run table event loop
		tbl = newTable(tableName)
		this.tables[tableName] = tbl
		tbl.stoper = this.stoper
		tbl.requests = make(chan *requestItem, config.CHAN_TABLE_REQUESTS_BUFFER_SIZE)
		go tbl.run()
	}
	// forward sql request to the table
	tbl.requests <- item
}

