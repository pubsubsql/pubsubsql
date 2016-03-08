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

package server

import (
	"os"
	"testing"
)

func TestMysqlConnection(t *testing.T) {
	conn := newMysqlConnection()
	defer conn.disconnect()

	connString := "pubsubsql:pubsubsql@/pubsubsql"
	if os.Getenv("TRAVIS") == "true" {
		connString = "travis:@/pubsubsql"
	}

	conn.connect(connString)
	if conn.hasError() {
		t.Error("failed to connect:", conn.getLastError())
	}

	if conn.isDisconnected() {
		t.Error("failed to open mysql connection (1):", conn.getLastError())
	}
	if !conn.isConnected() {
		t.Error("failed to open mysql connection (2):", conn.getLastError())
	}

	conn.disconnect()
	if !conn.isDisconnected() {
		t.Error("failed to close mysql connection (1):", conn.getLastError())
	}

	if conn.isConnected() {
		t.Error("failed to close mysql connection (2):", conn.getLastError())
	}
}

func TestMysqlConnectionFindTables(t *testing.T) {
	conn := newMysqlConnection()
	defer conn.disconnect()

	connString := "pubsubsql:pubsubsql@/pubsubsql"
	if os.Getenv("TRAVIS") == "true" {
		connString = "travis:@/pubsubsql"
	}

	conn.connect(connString)
	if conn.isDisconnected() {
		t.Error("failed to open mysql connection:", conn.getLastError())
	}
	t.Log(conn.findTables())
	if conn.hasError() {
		t.Error("failed to find tables:", conn.getLastError())
	}
}
