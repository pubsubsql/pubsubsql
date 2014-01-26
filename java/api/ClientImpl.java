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

import com.google.gson.Gson;
import java.nio.charset.Charset;
import java.util.*;

class ClientImpl implements Client {

	private final Charset UTF8_CHARSET = Charset.forName("UTF-8");
	private int CONNECTION_TIMEOUT = 500;
	private int requestId = 1;
	private String host;
	private int port;
	private String err;
	private NetHelper rw = new NetHelper();
	private ResponseData response = new ResponseData();
	private String rawjson = null;
	private int record = -1;
	private Hashtable<String, Integer> columns = new Hashtable<String, Integer>();
	private LinkedList<byte[]> backlog = new LinkedList<byte[]>();
	
	
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
		} 
		catch (Exception e) {
			setError(e);
		}	
		return false;
	}

	public void Disconnect() {
		backlog.clear();	
		write("close");
		// write may generate error so we reset after instead
		reset();
		rw.Close();
	}

	public boolean Connected() {
		return rw.Valid();
	}

	public boolean Ok() {
		return IsNullOrEmpty(err);
	}

	public boolean Failed() {
		return !Ok();
	}

	public String Error() {
		return NotNull(err);	
	}

	public boolean Execute(String command) {
		reset();
		boolean ok = write(command);
		NetHeader header = new NetHeader();
		while (ok) {
			reset();
			byte[] bytes = readTimeout(0, header);
			if (Failed()) return false;
			if (bytes == null) {
				setErrorString("Timed out");
				return false;
			}
			if (header.RequestId == requestId) {
				// response we are waiting for
				return unmarshallJSON(bytes);
			} else if (header.RequestId == 0) {
				//backlog
				backlog.add(bytes);
			} else if (header.RequestId < requestId) {
				// we did not read full result set from previous command irnore it
				reset();
			} else {
				invalidRequestIdError();
				return false;
			}
		}
		return false;
	}

	public String JSON() {
		return NotNull(rawjson);
	}

	public String Action() {
		return NotNull(response.action);
	}

	
	public String PubSubId() {
		return NotNull(response.pubsubid);
	}

	public int RowCount() {
		return response.rows;
	}

	public boolean NextRow() {
		while (Ok()) {
			// no result set
			if (response.rows == 0) return false;
			if (response.fromrow == 0 || response.torow == 0) return false;
			// the current record is valid
			record++;
			if (record <= (response.torow - response.fromrow)) return true;
			// we reached the end of the result set?
			if (response.rows == response.torow) {
				record--;
				return false;
			}
			// there is another batch of data
			reset();
			NetHeader header = new NetHeader();
			byte[] bytes = readTimeout(0, header);
			if (Failed()) return false;
			if (bytes == null) {
				setErrorString("Timed out");
				return false;
			}
			if (header.RequestId != requestId) {
				invalidRequestIdError();
				return false;
			}
			unmarshallJSON(bytes);
		}
		return false;
	}

	public String Value(String column) {
		if (response.data == null) return "";
		if (response.data.size() <= record) return "";
		int ordinal = getColumn(column);
		if (ordinal == -1) return "";
		//
		return response.data.get(record).get(ordinal);
	}

	public String ValueByOrdinal(int ordinal) {
		if (response.data == null) return "";
		if (response.data.size() <= record) return "";
		if (ordinal >= response.columns.size()) return "";
		return response.data.get(record).get(ordinal);
	}

	public boolean HasColumn(String column) {
        return getColumn(column) != -1;
	}

	public int ColumnCount() {
		if (response.columns == null) return 0;
		return response.columns.size();
	}

	public boolean WaitForPubSub(int timeout) {
		if (timeout <= 0) return false;		
		reset();
		// process backlog first
		if (backlog.size() > 0) {
			byte[] bytes = backlog.remove();
			return unmarshallJSON(bytes);
		}
		for (;;) {
			NetHeader header = new NetHeader();
			byte[] bytes = readTimeout(timeout, header);
			if (Failed()) return false;
			if (bytes == null) return false; 
			if (header.RequestId == 0) { 
				return unmarshallJSON(bytes);			
			}
			// this is not pubsub message; are we reading abandoned result set?
			// ignore and continue
		}
	}

	// helper functions

	private boolean IsNullOrEmpty(String str) {
		return (str == null || str.length() == 0);
	}

	private String NotNull(String str) {
		if (str == null) return "";
		return str;
	}

	private int toPort(String port) {
		try {
			return Integer.parseInt(port); 	
		} 
		catch (Exception e) {
				
		}
		return 0;
	}

	private void reset() {
		err = "";
		response = new ResponseData();
		rawjson = null;
		record = -1;
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
			rw.WriteWithHeader(requestId, message.getBytes(UTF8_CHARSET));
			return true;

		} 
		catch (Exception e) {
			hardDisconnect();
			setError(e);
		}
		return false;
	}

	private byte[] readTimeout(int timeout, NetHeader header) {
		try {
			if (!rw.Valid()) throw new Exception("Not connected");
			return rw.ReadTimeout(timeout, header);
		}		
		catch (Exception e) {
			hardDisconnect();
			setError(e);
		}
		return null;
	}

	private void setErrorString(String err) {
		reset();
		this.err = err;
	}

	private void setError(Exception e) {
		String err = e.getMessage();
		if (IsNullOrEmpty(err)) err = "Unknown error";
		setErrorString(err);
	}

	private void invalidRequestIdError() {
		setErrorString("Protocol error invalid request id");
	}

	private boolean unmarshallJSON(byte[] bytes) {
		try {
			Gson gson = new Gson();
			rawjson = new String(bytes, UTF8_CHARSET);
			response = gson.fromJson(rawjson, ResponseData.class);
			if (!response.status.equals("ok")) throw new Exception(response.msg); 
			setColumns();
			return true;
		}
		catch (Exception e) {
			setError(e);	
		}
		return false; 
	}

	void setColumns() {
		if (response.columns != null) {
			int index = 0;
			for(String column : response.columns) {
				columns.put(column, index);
				index++;
			}
		}	
	}

	int getColumn(String column) {
		Integer ordinal = columns.get(column);
		if (ordinal == null) return -1;
		return (int)ordinal;
	}
}

