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

	String host;
	int port;
	String err;
	NetHelper rw = new NetHelper();
	int CONNECTION_TIMEOUT = 500;
	int requestId = 1;

	public boolean Connect(String address) {
		Disconnect();
		// validate address
		int sep = address.indexOf(':');	
		if (sep < 0) {
			setErrorString("Invalid network address");
			return false;
		}
		// set host and port
		host = address.substring(0, sep);	
		int portIndex = sep + 1;
		if (portIndex >= address.length()) {
			setErrorString("Port is not provided");
			return false;
		}	
		int port = toPort(address.substring(portIndex));
		if (port == 0) {
			setErrorString("Invalid port");
			return false;
		}
		//
		try {
			java.net.Socket socket = new java.net.Socket();
			socket.connect(new java.net.InetSocketAddress(host, port), CONNECTION_TIMEOUT);
			rw.Set(socket);
			return rw.Valid();
		} catch (Exception e) {
			setError(e);
		}	
		return false;
	}

	public void Disconnect() {
		//backlog.Clear();	
		write("close");
		// write may generate error so we reset after instead
		reset();
		rw.Close();
	}

	public boolean Connected() {
		return false;
	}

	public boolean Ok() {
		return IsNullOrEmpty(err);
	}

	public boolean Failed() {
		return !Ok();
	}

	public String Error() {
		return err;	
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

	// helper functions

	private boolean IsNullOrEmpty(String str) {
		return (str == null || str.length() == 0);
	}

	private int toPort(String port) {
		try {
			return Integer.parseInt(port); 	
		} catch (Exception e) {
				
		}
		return 0;
	}

	private void reset() {
		err = "";
	}

	private void hardDisconnect() {
		//backlog.Clear();
		rw.Close();
		reset();
	}

	private boolean write(String message) {
		try {
			if (!rw.Valid()) throw new Exception("Not connected");
			requestId++;
			rw.WriteWithHeader(requestId, message.getBytes("UTF-8"));
			return true;

		} catch (Exception e) {
			hardDisconnect();
			setError(e);
		}
		return false;
	}

	private void setErrorString(String err) {
		reset();
		this.err = err;
	}

	private void setError(Exception err) {
		setErrorString(err.getMessage());
	}
}
