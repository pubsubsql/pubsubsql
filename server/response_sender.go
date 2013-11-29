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

// TODO design and implement
// responseSender is responsible for channeling reponses to client connection 
type responseSender struct {
	/*
		sender chan response // channel to publish responses to
		quit   chan int      // reader will close quite chan to indicate writers to quit	
		active bool          // indicates if channel is active
	*/
	i int
}

// factory
func newResponseSenderStub() *responseSender {
	return &responseSender{
		i: 0,
	}
}
