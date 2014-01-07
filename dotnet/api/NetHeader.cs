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

using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;

/*
--------------------+--------------------
|   message size    |    request id     |
--------------------+--------------------
|      uint32       |      uint32       |
--------------------+--------------------
*/

namespace PubSubSQL
{
    public struct NetHeader
    {
        public static readonly int HEADER_SIZE = sizeof(UInt32) + sizeof(UInt32);
        public UInt32 MessageSize;
        public UInt32 RequestId;

        public NetHeader(UInt32 messageSize, UInt32 requestId)
        {
            MessageSize = messageSize;
            RequestId = requestId;
        }

        public void ReadFrom(byte[] bytes)
        {
            setEndianess(bytes); 
            MessageSize = BitConverter.ToUInt32(bytes, 0);
            RequestId = BitConverter.ToUInt32(bytes, sizeof(UInt32));
        }

        public void WriteTo(byte[] bytes)
        {
            Array.Copy(BitConverter.GetBytes(MessageSize), 0, bytes, 0, sizeof(UInt32));
            Array.Copy(BitConverter.GetBytes(RequestId), 0, bytes, sizeof(UInt32), sizeof(UInt32));
            setEndianess(bytes); 
        }

        public byte[] GetBytes()
        {
            byte[] bytes = new byte[HEADER_SIZE];
            WriteTo(bytes);
            return bytes;
        }

        private void setEndianess(byte[] bytes)
        {
            if (BitConverter.IsLittleEndian)
            {
                Array.Reverse(bytes, 0, sizeof(UInt32));
                Array.Reverse(bytes, sizeof(UInt32), sizeof(UInt32));
            }
        }
    }
}
