using System;
using System.Collections.Generic;
using System.Text;
using System.Runtime.Serialization.Json;
using System.Runtime.Serialization;
using System.IO;

namespace PubSubSQL
{
    public interface Client
    {
        bool Connect(string address);
        void Disconnect();
        bool Ok();
        bool Failed();
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

    class client : Client
    {
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
            return false;
        }

        public void Disconnect()
        {

        }

        public bool Ok()
        {
            return false;
        }

        public bool Failed()
        {
            return true;
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
    }

}
