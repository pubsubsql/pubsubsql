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

// requestRouter routs request to appropriate service for processing
type requestRouter struct {
	dataSrv *dataService
}

// requestRouter factory
func newRequestRouter(dataSrv *dataService) *requestRouter {
	return &requestRouter{dataSrv: dataSrv}
}

func (rt *requestRouter) onError(r *requestItem) {
	ereq := r.req.(*errorRequest)
	res := newErrorResponse(ereq.err)
	r.sender.send(res)
}

func (rt *requestRouter) route(r *requestItem) {
	switch r.req.getRequestType() {
	case requestTypeSql:
		rt.dataSrv.accept(r)
	case requestTypeError:
		rt.onError(r)
	default:
		panic("unsuported request type")
	}
}
