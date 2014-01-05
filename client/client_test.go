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

import  (
	"testing"
	"encoding/json"
)

var address = "localhost:7777"

func TestConnectDisconnect(t *testing.T) {
	client := NewClient()
	if !client.Connect(address) {
		t.Error("Connect failed.", client.ErrorString())
		return
	}
	if client.Failed() {
		t.Error("unexpected error")
	}		
	client.Disconnect()
}

type unmarshall struct {
	Status string //`json:"status"`
	Data []map[string]string
}

func TestUnmarshall(t *testing.T) {
	client := NewClient()
	client.Connect(address)

	if !client.Execute("insert into stocks (ticker, bid) values (IBM, 140.45)") {
		t.Error("failed to execute status command", client.ErrorString())
	}
	if !client.Execute("select * from stocks") {
		t.Error("failed to execute status command", client.ErrorString())
	}
	var status unmarshall
	message := client.JSON()
	if err := json.Unmarshal([]byte(message), &status); err != nil {
		t.Error(err.Error())
	}
	if status.Status != "ok" {
		t.Error("expected status ok")
	}
	if status.Data[0]["ticker"] != "IBM" {
		t.Error("expected ticker IBM")
	}
	
	client.Disconnect()
}


