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

type responseStatusType int8

const (
	responseStatusOk  responseStatusType = iota // ok.
	responseStatusErr                           // error.
)

// response 
type response interface {
	getResponseStatus() responseStatusType
}

// statusResponse 
type statusResponse struct {
	response
	status responseStatusType
	msg    string
}

func (r *statusResponse) getResponsStatus() responseStatusType {
	return r.status
}

// TODO String stub will optimize later
func (r *statusResponse) String() string {
	if r.status == responseStatusOk {
		return `{"status":"ok"}`
	}
	return `{"status":"err" "msg":"` + r.msg + `"}`
}
