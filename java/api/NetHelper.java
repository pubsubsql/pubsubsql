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

public class NetHelper {

	java.net.Socket socket;
	byte[] headerBytes = new byte[NetHeader.HEADER_SIZE];

	public void Set(java.net.Socket socket) {
		this.socket = socket;
	}

	public boolean Valid() {
		return this.socket != null && this.socket.isConnected();            
	}	

	public void Close() {
		if (socket == null) return;
        try {
			socket.shutdownOutput();
			socket.close();
			socket = null;
		}
		catch (Exception e) {
			// ignore
		}
	}

	public void WriteWithHeader(int requestId, byte[] bytes) throws java.io.IOException {
		NetHeader header = new NetHeader(bytes.length, requestId);
		java.io.OutputStream stream = socket.getOutputStream();
		stream.write(header.GetBytes());
		stream.write(bytes);
		stream.flush();
	}
	
	public byte[] ReadTimeout(int timeout, NetHeader header)throws java.io.IOException {
		try {
			socket.setSoTimeout(timeout);
			return Read(header);			
		} 
		catch (java.net.SocketTimeoutException te) {
			// ignore	
		}
		return null;
	}

	public byte[] Read(NetHeader header) throws java.io.IOException {
		java.io.InputStream stream = socket.getInputStream();
		int read = stream.read(headerBytes);
		if (read < NetHeader.HEADER_SIZE) throw new Exception("Failed to read header");
		header.ReadFrom(headerBytes);	
		bytes = new byte[NetHeader.HEADER_SIZE];
		read = 0;
		while (header.MessageSize > read) {
			read +=  stream.read(bytes, read, header.MessageSize - read);
		}
		return bytes;
	}
}

