using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Net.Sockets;

namespace PubSubSQL
{
    public class NetHelper
    {
        Socket socket;
        byte[] headerBytes;

        public void Set(Socket socket, int bufferSize)
        {
            this.socket = socket;
            this.headerBytes = new byte[NetHeader.HEADER_SIZE];
        }

        public void Close()
        {
            if (socket != null)
            {
                try
                {
                    socket.Shutdown(SocketShutdown.Both);
                    socket.Close();
                    socket = null;
                }
                catch (Exception )
                {
                    // ignore
                }
            }
        }

        public bool Valid()
        {
            return this.socket != null;
        }

        public void Write(byte[] bytes)
        {
            int written = 0;
            while (bytes.Length > written)
            {
                written += socket.Send(bytes, written, bytes.Length - written, SocketFlags.None);
            }
        }

        public void WriteWithHeader(UInt32 requestId, byte[] bytes)
        {
            NetHeader header = new NetHeader((UInt32)bytes.Length, requestId);
            Write(header.GetBytes());
            Write(bytes);
        }

        public bool ReadTimeout(int timeout, ref NetHeader header, out byte[] bytes)
        {
            bytes = null;
            try
            {
                socket.ReceiveTimeout = timeout;
                Read(ref header, out bytes);
                return true;
            }
            catch (SocketException e)
            {
                if (e.SocketErrorCode == SocketError.TimedOut)
                {
                    return false;
                }
                else
                {
                    throw e;
                }
            }
        }

        public void Read(ref NetHeader header, out byte[] bytes)
        {
            int read = socket.Receive(headerBytes, 0, NetHeader.HEADER_SIZE, SocketFlags.None);
            if (read < NetHeader.HEADER_SIZE)
            {
                throw new Exception("Failed to read header.");
            }
            header.ReadFrom(headerBytes);
            bytes = new byte[header.MessageSize];
            // read the rest of the message
            read = 0;
            while (header.MessageSize > read)
            {
                read += socket.Receive(bytes, read, (int)header.MessageSize - read, SocketFlags.None);
            }
        }

        public static byte[] ToUTF8(string str)
        {
            return System.Text.Encoding.UTF8.GetBytes(str);
        }

        public static string FromUTF8(byte[] bytes)
        {
            return System.Text.Encoding.UTF8.GetString(bytes);
        }
    }
}
