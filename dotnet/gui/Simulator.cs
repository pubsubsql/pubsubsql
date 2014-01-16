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
        private bool stopFlag = false;
        Thread thread = null;

        public volatile int TotalPublished = 0;
        public volatile int TotalConsumed = 0;

        private void Run()
        {
            try
            {
                TotalPublished = 0;
                TotalConsumed = 0;
                Thread.Sleep(0);
                if (!client.Connect(Address)) throw new Exception("Failed to connect");
                if (!client.Execute(string.Format("key {0} col1", TableName))) throw new Exception(client.Error());
                // first insert data
                for (int row = 1; row <= Rows && !stopFlag; row++)
                {
                    string insert = generateInsert(row);
                    if (!client.Execute(insert)) throw new Exception("Failed to insert: " + insert);            
                    while (TotalPublished - TotalConsumed > 2000 && !stopFlag)
                    {
                        Thread.Sleep(50);
                    }
                }
                //System.Windows.Forms.MessageBox.Show("INSERTED");
                while (!stopFlag)
                {
                    string update = generateUpdate();
                    if (!client.Execute(update)) throw new Exception(client.Error());
                    TotalPublished++;
                    while (TotalPublished - TotalConsumed > 2000 && !stopFlag)
                    {
                        Thread.Sleep(50);
                    }
                }
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
            int row = rnd.Next(1, Rows + 1);
            int col = rnd.Next(2, Columns + 1);
            int value = rnd.Next(1, 1000000);
            return string.Format("update {0} set col{1} = {2} where col1 = {3}", TableName, col, value, row); 
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
