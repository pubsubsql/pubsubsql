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

import (
	"encoding/binary"
	"encoding/json"
)

/*
--------------------+--------------------
|   message size    |    request id     |
--------------------+--------------------
|      uint32       |      uint32       |
--------------------+--------------------
*/

type NetworkHeader struct {
	MessageSize uint32
	RequestId uint32
}

var HEADER_SIZE = 8
var EMPTY_HEADER = make([]byte, HEADER_SIZE, HEADER_SIZE)

func NewNetworkHeader(messageSize uint32, requestId uint32) *NetworkHeader {
	return &NetworkHeader {
		MessageSize: messageSize,
		RequestId: requestId,
	}
}

func (this *NetworkHeader) ReadFrom(bytes []byte) {
	this.MessageSize = binary.LittleEndian.Uint32(bytes)
	this.RequestId = binary.LittleEndian.Uint32(bytes[4:])
}

func (this *NetworkHeader) WriteTo(bytes []byte) {
	binary.LittleEndian.PutUint32(bytes, this.MessageSize)
	binary.LittleEndian.PutUint32(bytes[4:], this.RequestId)
}

func (this *NetworkHeader) GetBytes() []byte {
	bytes := make([]byte, HEADER_SIZE, HEADER_SIZE)
	this.WriteTo(bytes);
	return bytes;
}

func (this *NetworkHeader) String() string {
	bytes, err := json.Marshal(this)	
	return string(bytes)
} 

