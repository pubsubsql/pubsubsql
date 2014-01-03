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
	"testing"
	"fmt"
)

func TestNetworkHeader(t *testing.T) {
	header1 := NetworkHeader {
		MessageSize: 32567,
		RequestId: 9875235,
	}	
	var header2 NetworkHeader
	bytes := make([]byte, 100, 100)
	//
	header1.WriteTo(bytes)
	header2.ReadFrom(bytes)
	//
	if header1 != header2 {
		t.Error("NetworkHeader data does not match")
	}	
	
	fmt.Println(header1.String())
}
