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

import java.util.*;

public class TableDataset {

	public static class Cell {
		public String Value;
		public long LastUpdated;

		public Cell(String value) {
			Value = value;
			LastUpdated = System.nanoTime();
		}	
	}

	private ArrayList<String> columns = new ArrayList<String>();
	private Hashtable<String, Integer> columnOrdinals = new Hashtable<String, Integer>();
	private ArrayList<ArrayList<Cell>> rows = new ArrayList<ArrayList<Cell>>();
	private Hashtable<String, ArrayList<Cell>> idsToRows = new Hashtable<String, ArrayList<Cell>>();
	private volatile boolean dirtyFlag = false;
	private volatile boolean clearFlag = false;
	
	public boolean ResetDirty() {
		boolean ret = dirtyFlag;
		dirtyFlag = false;
		return ret;
	}

	public boolean ResetClear() {
		boolean ret = clearFlag;
		clearFlag = false;
		return ret;
	}	

	public void Clear() {
		columns.clear();
		columnOrdinals.clear();
		rows.clear();
		idsToRows.clear();
		dirtyFlag = true;
		clearFlag = true;
	}

	public void SyncColumns(pubsubsql.Client client) {
		for(String col : client.Columns()) {
			if (!columnOrdinals.containsKey(col)) {
				clearFlag = true;
				int ordinal = columns.size();
				columnOrdinals.put(col, ordinal);
				columns.add(col);
			}
		}
	}

	public void ProcessRow(pubsubsql.Client client) {
		dirtyFlag = true;
		String id = client.Value("id");
		ArrayList<Cell> row = null;
		switch (client.Action()) {
			case "select":
			case "add":
			case "insert":
				// add row
				row = new ArrayList<Cell>(columns.size());
				// for each selct operations columns are always in the same order
				for (String col : columns) {
					row.add(new Cell(client.Value(col)));	
				}
				rows.add(row);
				if (id.length() > 0) {
					idsToRows.put(id, row);
				}
				break;
			case "update":
				row = idsToRows.get(id);
				if (row != null) {
					for (String col : client.Columns()) {
						Integer ordinal = columnOrdinals.get(col);
						// auto expand row
						for (int i = row.size(); i <= ordinal; i++) {
							row.add(new Cell(""));
						}
						Cell cell = row.get(ordinal);
						cell.Value = client.Value(col);
						cell.LastUpdated = System.nanoTime();
					}
				}
				break;

			case "delete":
			case "remove":
				row = idsToRows.get(id);
				if (row != null) {
					idsToRows.remove(id);
					rows.remove(row);
				}
				break;
		}
	}
	
	public ArrayList<Cell> Row(int rowIndex) {
		if (rowIndex < rows.size()) return rows.get(rowIndex);
		return new ArrayList<Cell>();
	}

	public int RowCount() {
		return rows.size();
	}

	public int ColumnCount() {
		return columns.size();
	}
	
	public String Column(int colIndex) {
		if (colIndex < columns.size()) {
			return columns.get(colIndex);
		}
		return "";
	}
} 
