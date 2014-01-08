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
using System.Text;
using System.Runtime.Serialization.Json;
using System.Runtime.Serialization;
using System.IO;
using System.Net.Sockets;

namespace PubSubSQL
{
    public interface Client
    {
        bool Connect(string address);
        void Disconnect();
        bool Ok();
        bool Failed();
        string Error();
        bool Execute(string command);
        string JSON();
        string Action();
        string Id();
        string PubSubId();
        int RecordCount();
        bool NextRecord();
        string Value();
        bool HasColumn();
        List<string> Columns();
        bool WaitForPubSub(Int64 timeout);
    }

    [DataContract]
    class responseData
    {
        [DataMember(Name = "status")]
        public string Status { get; set; }
        [DataMember(Name = "msg")]
        public string Msg { get; set; }
        [DataMember(Name = "action")]
        public string Action { get; set; }
        [DataMember(Name = "id")]
        public string Id { get; set; }
        [DataMember(Name = "pubsubid")]
        public string PubSubId { get; set; }
        [DataMember(Name = "rows")]
        public string Rows { get; set; }
        [DataMember(Name = "fromrow")]
        public string Fromrow { get; set; }
        [DataMember(Name = "torow")]
        public string Torow { get; set; }
        [DataMember(Name = "columns")]
        public string[] columns { get; set; }
        //Data     []map[string]string
        // TOBE DECIDED
    }

    public class Factory
    {
        public static Client NewClient()
        {
            return new client();
        }
    }

    class client : Client
    {
        string host;
        int port;
        NetHelper rw = new NetHelper();
        UInt32 requestId;
        string err;
        byte[] rawjson;
        responseData response;
        int record;

        const int CLIENT_DEFAULT_BUFFER_SIZE = 2048;

        public void testColmpile()
        {
            string str = "{\"status\":\"ok\",\"columns\":[\"col1\",\"col2\"]}"; 
            byte[] byteArray = Encoding.ASCII.GetBytes( str );
            MemoryStream stream = new MemoryStream( byteArray );
            DataContractJsonSerializer jsonSerializer = new DataContractJsonSerializer(typeof(responseData));
            object objResponse = jsonSerializer.ReadObject(stream);
            responseData jsonResponse = objResponse as responseData;
            System.Console.WriteLine(jsonResponse.Status);
            System.Console.WriteLine(jsonResponse.Msg);
            foreach (string s in jsonResponse.columns)
            {
                System.Console.WriteLine(s);
            }
        }

        public bool Connect(string address)
        {
            Disconnect();
            // validate address
            int sep = address.IndexOf(':');
            if (sep < 0)
            {
                setErrorString("invalid network address");
                return false;
            }
            // set host and port
            host = address.Substring(0, sep);
            if (!toPort(ref port, address.Substring(sep + 1))) return false;
            // connect
            try
            {
                Socket socket = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
                socket.Connect(host, port);
                rw.Set(socket, CLIENT_DEFAULT_BUFFER_SIZE); 
                return true;
            }
            catch (Exception e)
            {
                setError(e);
            }
            //
            return false;
        }

        public void Disconnect()
        {
            write("close");
            // write may generate errro so we reset after instead
            reset();
            rw.Close();
        }

        public bool Ok()
        {
            return string.IsNullOrEmpty(err); 
        }

        public bool Failed()
        {
            return !Ok(); 
        }

        public string Error()
        {
            return err;
        }

        public bool Execute(string command)
        {

            return false;
        }

        public string JSON()
        {
            return "";
        }

        public string Action()
        {
            return "";
        }

        public string Id()
        {
            return "";
        }

        public string PubSubId()
        {
            return "";
        }

        public int RecordCount()
        {
            return 0;
        }

        public bool NextRecord()
        {
            return false;
        }

        public string Value()
        {
            return "";
        }

        public bool HasColumn()
        {
            return false;
        }

        public List<string> Columns()
        {
            return null;
        }

        public bool WaitForPubSub(Int64 timeout)
        {
            return false;
        }

        void reset()
        {
            err = string.Empty;
        }

        bool toPort(ref int port, string sport)
        {
            try
            {
                port = Convert.ToInt32(sport, 10);
                return true;
            }
            catch (Exception )
            {
                setErrorString("Invalid port " + sport);
            }
            return false;
        }

        void setErrorString(string err)
        {
            reset();
            this.err = err;
        }

        void setError(Exception e)
        {
            setErrorString(e.Message);
        }

        bool write(string message)
        {
            try
            {
                if (!rw.Valid()) throw new Exception("Not connected");
                requestId++;
                rw.WriteWithHeader(requestId, NetHelper.ToUTF8(message));
            }
            catch (Exception e)
            {
                setError(e);
                return false;
            }
            return true;
        }

        bool readTimeout(int timeout, ref NetHeader header, out byte[] bytes, ref bool timedout)
        {
            timedout = false;
            bytes = null;
            try
            {
                if (!rw.Valid()) throw new Exception("Not connected");
                if (!rw.ReadTimeout(timeout, ref header, out bytes))
                {
                    timedout = true;
                }
                return true;
            }
            catch (Exception e)
            {
                setError(e);
            }
            return false;
        }

        bool read(ref NetHeader header, out byte[] bytes)
        {
            const int MAX_READ_TIMEOUT_MILLISECONDS = 1000 * 60 * 3;
            bool timedout = false;
            bool err = readTimeout(MAX_READ_TIMEOUT_MILLISECONDS, ref header, out bytes, ref timedout);
            if (timedout)
            {
                setErrorString("Read timed out");
            }
            return timedout || err;
        }

    }

}
