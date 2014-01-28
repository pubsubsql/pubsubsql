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

import java.awt.*;
import java.awt.event.*;
import javax.swing.*;
import javax.swing.table.*;
import java.util.*;

public class TableView extends JPanel {

	public int FLASH_TIMEOUT;
	private TableDataset dataset;	
	private JTable table;
	private TableModel model;
		
	public TableView(int flashTimeout, TableDataset dataset) {
		this.FLASH_TIMEOUT = flashTimeout;
		this.dataset = dataset;
		setLayout(new BorderLayout());
		model = this.new TableModel();
		table = new JTable(model);
		table.setDefaultRenderer(Object.class, this.new CellRenderer());
		add(new JScrollPane(table));
	}

	public void Update() {
		model.Update();
	}

	// Model
	private class TableModel extends AbstractTableModel {

		private int rows = 0;
		private int cols = 0;

		public void Update() {
			boolean structureChanged = false;
			if (dataset.RowCount() != rows) structureChanged = true;
			if (dataset.ColumnCount() != cols) structureChanged = true;
			if (dataset.ResetClear()) structureChanged = true;
			rows = dataset.RowCount();
			cols = dataset.ColumnCount();
			if (structureChanged) fireTableStructureChanged();
			else fireTableDataChanged();
		}

		// AbstractTableModel
		public int getRowCount() {
			return dataset.RowCount();
		}

		public int getColumnCount() {
			return dataset.ColumnCount();
		}

		public Object getValueAt(int r, int c) {
			ArrayList<TableDataset.Cell> row = dataset.Row(r);				
			if (row.size() <= c) return null;
			return row.get(c);
		}

		@Override
		public String getColumnName(int c) {
			return dataset.Column(c);		
		}
	}

	// Renderer	
	private class CellRenderer extends DefaultTableCellRenderer {
		
		@Override
		public Component getTableCellRendererComponent(JTable table, Object value, boolean isSelected, 
															boolean hasFocus, int r, int c) {
						
			Color backcolor = Color.white;
			if (value == null) {
				 setValue("");
				 setBackground(backcolor);
				 return this;
			}
			TableDataset.Cell cell = (TableDataset.Cell)value;
			if (System.nanoTime() - cell.LastUpdated < FLASH_TIMEOUT) backcolor = Color.pink;
			setBackground(backcolor);
			setValue(cell.Value);
			return this;
		}
	}
}
