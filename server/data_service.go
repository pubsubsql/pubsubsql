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
	req       request
	responser *responseSender
}

// dataService prer-processes sqlRequests and channels them to approptiate tables for further proccessging 
type dataService struct {
	requests chan dataServiceItem
	stoper   *Stoper
}

// dataService factory
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
			d.onSqlRequest(item.req, item.responser)
		case <-d.stoper.GetChan():
			return
		}
	}
}

func (d *dataService) onSqlRequest(r request, responser *responseSender) {
	switch r.(type) {
	case *sqlInsertRequest:
		d.sqlInsert(r.(*sqlInsertRequest), responser)

	case *sqlSelectRequest:
		d.sqlSelect(r.(*sqlSelectRequest), responser)

	case *sqlUpdateRequest:
		d.sqlUpdate(r.(*sqlUpdateRequest), responser)

	case *sqlDeleteRequest:
		d.sqlDelete(r.(*sqlDeleteRequest), responser)

	case *sqlSubscribeRequest:
		d.sqlSubscribe(r.(*sqlSubscribeRequest), responser)

	case *sqlUnsubscribeRequest:
		d.sqlUnsubscribe(r.(*sqlUnsubscribeRequest), responser)

	case *sqlKeyRequest:
		d.sqlKey(r.(*sqlKeyRequest), responser)

	case *sqlTagRequest:
		d.sqlTag(r.(*sqlTagRequest), responser)

	default:
		panic("Unsuported sql request")
	}
}

func (d *dataService) sqlInsert(req *sqlInsertRequest, responser *responseSender) {

}

func (d *dataService) sqlSelect(req *sqlSelectRequest, responser *responseSender) {

}

func (d *dataService) sqlUpdate(req *sqlUpdateRequest, responser *responseSender) {

}

func (d *dataService) sqlDelete(req *sqlDeleteRequest, responser *responseSender) {

}

func (d *dataService) sqlSubscribe(req *sqlSubscribeRequest, responser *responseSender) {

}

func (d *dataService) sqlUnsubscribe(req *sqlUnsubscribeRequest, responser *responseSender) {

}

func (d *dataService) sqlKey(req *sqlKeyRequest, responser *responseSender) {

}

func (d *dataService) sqlTag(req *sqlTagRequest, responser *responseSender) {

}
