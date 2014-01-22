/* Copyright (C) 2014 CompleteDB LLC.
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

package pubsubsql; 

class client implements Client {

	public boolean Connect(String address) {
		return false;
	}

	public void Disconnect() {

	}

	public boolean Connected() {
		return false;
	}

	public boolean Ok() {
		return false;
	}

	public boolean Failed() {
		return true;
	}

	public String Error() {
		return "";	
	}

	public boolean Execute(String command) {
		return false;
	}

	public String JSON() {
		return "";
	}

	public String Action() {
		return "";
	}

	public int RecordCount() {
		return 0;
	}

	public boolean NextRecord() {
		return false;
	}

	public String Value(String column) {
		return "";
	}

	public boolean HasColumn(String column) {
		return false;
	}

	// Columns();
	public int ColumnCount() {
		return 0;
	}

	public String Column(int index) {
		return "";
	}

	public boolean WaitForPubSub(int timeout) {
		return false;
	}

}
