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

type dataServiceItem struct {
	req       *sqlRequest
	responser *responseSender
}

// responseSender is responsible for channeling reponses to client connection 
type dataService struct {
	requests chan dataServiceItem
	stoper   *Stoper
}

// factory
func newDataService(bufferSize int, stoper *Stoper) *dataService {
	if stoper == nil {
		panic("dataService.stoper can not be nil")
	}
	return &dataService{
		requests: make(chan dataServiceItem, bufferSize),
		stoper:   stoper,
	}
}

// runs dataService event loop
func (d *dataService) run() {
	d.stoper.Enter()
	defer d.stoper.Leave()
	for {
		select {
		case item := <-d.requests:
			onSqlRequest(item.req, item.responser)
		case <-d.stoper.GetChan():
			return
		}
	}
}

// 
func onSqlRequest(r *sqlRequest, responser *responseSender) {

}

func (d *dataService) findTable(table string) bool {
	return false
}

func (d *dataService) createTable(table string, columns []string) {

}

func (d *dataService) sqlInsert(req *sqlInsertRequest, responser *responseSender) {

}

func (d *dataService) sqlSelect(req *sqlSelectRequest, responser *responseSender) {

}
