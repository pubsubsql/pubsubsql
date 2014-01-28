using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;

namespace PubSubSQLGUI
{
    class Cell
    {
        public string Value;
        public long LastUpdated;

        public Cell(string value)
        {
            Value = value;
            LastUpdated = DateTime.Now.Ticks;
        }
    }

    class ListViewDataset
    {
        private List<string> columns = new List<string>();
        private Dictionary<string, int> columnOrdinals = new Dictionary<string, int>();
        private List<List<Cell>> rows = new List<List<Cell>>();
        private Dictionary<string, List<Cell>> idsToRows = new Dictionary<string, List<Cell>>();
        private volatile bool dirtyFlag = false;

        public bool ResetDirty()
        {
            bool ret = dirtyFlag;
            dirtyFlag = false;
            return ret;
        }

        public void Clear()
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
            dirtyFlag = true;
            string id = client.Value("id");
            List<Cell> row = null;
            switch (client.Action())
            {
                case "select":
                case "add":
                case "insert":
                    // add row
                    row = new List<Cell>(columns.Count);
                    // for select operations columns are always in the same order
                    foreach(string col in columns)
                    {
                        row.Add(new Cell(client.Value(col)));
                    }
                    rows.Add(row);
                    if (!string.IsNullOrEmpty(id))
                    {
                        idsToRows[id] = row;                
                    }
                    break;
                case "update":
                    if (idsToRows.TryGetValue(id, out row))
                    {
                        foreach (string col in client.Columns())
                        {
                            int ordinal = columnOrdinals[col];
                            // auto expand row
                            for (int i = row.Count; i <= ordinal; i++)
                            {
                                row.Add(new Cell(string.Empty));
                            }
                            row[ordinal].Value = client.Value(col);
                            row[ordinal].LastUpdated = DateTime.Now.Ticks;
                        }
                    }
                    break;
                case "delete":
                case "remove":
                    if (idsToRows.TryGetValue(id, out row))
                    {
                        idsToRows.Remove(id);
                        rows.Remove(row);
                    }
                    break;
            }
        }

        public List<Cell> GetRow(int rowIndex)
        {
            if (rowIndex < rows.Count)
            {
                return rows[rowIndex];     
            }
            return new List<Cell>();
        }

        public int RowCount
        {
            get {return rows.Count;}
        }

        public int ColumnCount
        {
            get { return columns.Count; }
        }

        public string Column(int colIndex)
        {
            if (colIndex < columns.Count)
            {
                return columns[colIndex];
            }
            return string.Empty;
        }
    }
}
