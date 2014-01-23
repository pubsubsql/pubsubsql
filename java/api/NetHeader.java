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

import java.nio.ByteBuffer;

/*
--------------------+--------------------
|   message size    |    request id     |
--------------------+--------------------
|      uint32       |      uint32       |
--------------------+--------------------
*/

public class NetHeader {

	// JAVA VM always uses BIG ENDIAN to encode integer types so no conversion is nessasary
    public static final int HEADER_SIZE = 8;
    public int MessageSize;
    public int RequestId;

	public NetHeader(int messageSize, int requestId) {
		MessageSize = messageSize;
		RequestId = requestId;
    }

	public void ReadFrom(byte[] bytes) {
		ByteBuffer buffer = ByteBuffer.wrap(bytes);
		MessageSize = buffer.getInt(); 
		RequestId = buffer.getInt();
	}

	public void WriteTo(byte[] bytes) {
		ByteBuffer buffer = ByteBuffer.wrap(bytes);
		buffer.putInt(MessageSize);
		buffer.putInt(RequestId);
	}

	public byte[] GetBytes() {
		byte[] bytes = new byte[HEADER_SIZE];
		WriteTo(bytes);
		return bytes;
	}

}
