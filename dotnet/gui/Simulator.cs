using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading;

namespace PubSubSQLGUI
{
    class Simulator
    {
        public int Columns = 0;
        public int Rows = 0;
        public string TableName = string.Empty;
        public string Address = string.Empty;
        private PubSubSQL.Client client = PubSubSQL.Factory.NewClient();
        private volatile bool stopFlag = false;
        Thread thread = null;
        List<string> ids = new List<string>();

        private void Run()
        {
            try
            {
                ids.Clear();
                if (!client.Connect(Address)) throw new Exception("Failed to connect");
                // first insert data
                for (int row = 1; row <= Rows && !stopFlag; row++)
                {
                    string insert = generateInsert(row);
                    if (!client.Execute(insert)) throw new Exception("Failed to insert: " + insert);
                    if (!client.NextRow()) throw new Exception("Failed to move to the first record");
                    string id = client.Value("id");
                    if (string.IsNullOrEmpty(id)) throw new Exception("id is empty");
                    ids.Add(id);
                }
                //
                while (!stopFlag)
                {
                    for (int i = 0; i < 100 && !stopFlag; i++)
                    {
                        string update = generateUpdate();
                        if (!client.Execute(update)) throw new Exception(client.Error());
                    }
                    // gui thread can not process that many messages from the server
                    // slow down the updates
                    Thread.Sleep(100);
                }
                client.Disconnect();
            }
            catch (Exception e)
            {
                System.Windows.Forms.MessageBox.Show(e.Message);
            }
            finally
            {
                client.Disconnect();
            }
        }

        public void Reset()
        {
            Columns = 0;
            Rows = 0;
            TableName = string.Empty;
            Address = string.Empty;
            thread = null;
        }

        public void Start()
        {
            Stop();
            stopFlag = false;
            thread = new System.Threading.Thread(Run);
            thread.Start();
        }

        public void Stop()
        {
            stopFlag = true;
            if (thread != null)
            {
                thread.Join();
                thread = null;
            }
        }

        Random rnd = new Random(DateTime.Now.Second);
        private string generateUpdate()
        {
            int idIndex = rnd.Next(0, ids.Count);
            string id = ids[idIndex];
            int col = rnd.Next(1, Columns + 1);
            int value = rnd.Next(1, 1000000);
            return string.Format("update {0} set col{1} = {2} where id = {3}", TableName, col, value, id); 
        }

        private string generateInsert(int row)
        {
            StringBuilder builder = new StringBuilder();
            builder.Append("insert into ");
            builder.Append(TableName);
            // columns
            for (int i = 0; i < Columns; i++)
            {
                if (i == 0) builder.Append(" ( ");
                else builder.Append(" , ");
                builder.Append(string.Format("col{0}", i + 1));
            }
            // values
            builder.Append(") values ");
            for (int i = 0; i < Columns; i++)
            {
                if (i == 0) builder.Append(" ( ");
                else builder.Append(" , ");
                builder.Append(row.ToString());
            }
            builder.Append(")");
            return builder.ToString();
        }
    }
}
