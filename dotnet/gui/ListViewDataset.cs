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

        public void SyncColumns(PubSubSQL.Client client)
        {
            foreach (string col in client.Columns())
            {
                if (!columnOrdinals.ContainsKey(col))
                {
                    int ordinal = columns.Count;
                    columnOrdinals[col] = ordinal;
                    columns.Add(col);
                }
            }
        }

        public void ProcessRow(PubSubSQL.Client client)
        {
            switch (client.Action())
            {
                case "select":
                case "add":
                case "insert":
                    // add row
                    List<string> row = new List<string>(columns.Count);
                    // for select operations columns are always in the same order
                    foreach(string col in columns)
                    {
                        row.Add(client.Value(col));
                    }
                    rows.Add(row);
                    break;
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

        public int ColumnCount
        {
            get { return columns.Count; }
        }

        public string Column(int index)
        {
            if (index < columns.Count)
            {
                return columns[index];
            }
            return string.Empty;
        }
    }
}
