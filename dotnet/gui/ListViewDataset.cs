using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;

namespace PubSubSQLGUI
{
    class ListViewDataset
    {
        private List<string> columns = new List<string>();
        private Dictionary<string, int> columnOrdinals = new Dictionary<string, int>();
        private List<List<string>> rows = new List<List<string>>();
        private Dictionary<string, List<string>> idsToRows = new Dictionary<string, List<string>>();

        public void Reset()
        {
            columns.Clear();
            columnOrdinals.Clear();
            rows.Clear(); 
            idsToRows.Clear();
        }

        public void AddRowsCapacity(int capacity)
        {
            int remainingCapacity = rows.Capacity - rows.Count;
            if (remainingCapacity < capacity)
            {
                rows.Capacity += capacity - remainingCapacity;
            }
        }

        public void ProcessRow(PubSubSQL.Client client)
        {
            if (client.Action() == "select")
            {
                // lazy add columns
                if (columns.Count == 0)
                {
                    columns = new List<string>(client.ColumnCount());
                    foreach (string col in client.Columns())
                    {
                        columns.Add(col);
                    }
                }
                // add row
                List<string> row = new List<string>(columns.Count);
                foreach(string col in columns)
                {
                    row.Add(client.Value(col));
                }
                rows.Add(row);
            }
            else if (client.Action() == "add" || client.Action() == "insert")
            {
                
            }
            else if (client.Action() == "update")
            {
    
            }
            else if (client.Action() == "delete")
            {

            }
        }

        public List<string> GetRow(int rowIndex)
        {
            if (rowIndex < rows.Count)
            {
                return rows[rowIndex];     
            }
            return new List<string>();
        }

        public int RowCount
        {
            get {return rows.Count;}
        }

        int ColumnCount
        {
            get { return columns.Count; }
        }

        List<string> Columns
        {
            get { return columns; }
        }



    }
}
