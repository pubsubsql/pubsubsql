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

public class MainForm extends JFrame {

	public MainForm() {
		setTitle("Interactive Query");
		setupMenuBar();		
	}

	void setupMenuBar() {
		JMenuBar menuBar = new JMenuBar();
		this.setJMenuBar(menuBar);	
		// File
		JMenu fileMenu = new JMenu("File");
			// New
			JMenuItem newMenu = new JMenuItem(new_);	
			newMenu.setToolTipText("New PubSubSQL Interactive Query");
			fileMenu.add(newMenu);
			// Exit	
			JMenuItem exitMenu = new JMenuItem(exit);
			defaultTooltips(exitMenu);
			fileMenu.add(exitMenu);
		menuBar.add(fileMenu);
		// Connection
		JMenu connectionMenu = new JMenu("Connection");
			// Connect local
			JMenuItem connectLocalMenu = new JMenuItem(connectLocal);
			connectLocalMenu.setToolTipText(connectLocalMenu.getText());
			connectionMenu.add(connectLocalMenu);
			// Connect
			JMenuItem connectMenu = new JMenuItem(connect);
			connectMenu.setToolTipText("Connect to remote server");
			connectionMenu.add(connectMenu);
			// Disconnect
			JMenuItem disconnectMenu = new JMenuItem(disconnect);
			defaultTooltips(disconnectMenu);
			connectionMenu.add(disconnectMenu);
		menuBar.add(connectionMenu);	
		// Query
		JMenu queryMenu = new JMenu("Query");
			// Execute 
			JMenuItem executeMenu = new JMenuItem(execute);
			defaultTooltips(executeMenu);
			queryMenu.add(executeMenu);
			// Cancel Executing Query 
			JMenuItem cancelMenu = new JMenuItem(cancelExecute);
			defaultTooltips(cancelMenu);
			queryMenu.add(cancelMenu);
			// Simulate 
			JMenuItem simulateMenu = new JMenuItem(simulate);
			defaultTooltips(simulateMenu);
			queryMenu.add(simulateMenu);
		menuBar.add(queryMenu);	
		// Help
		JMenu helpMenu = new JMenu("Help");
			// About 
			JMenuItem aboutMenu = new JMenuItem(about);
			defaultTooltips(aboutMenu);
			helpMenu.add(aboutMenu);
		menuBar.add(helpMenu);	
	
	}

	// events
	Action new_ = new AbstractAction("New") {
		public void actionPerformed(ActionEvent event) {
			
		}
	};
	
	Action exit = new AbstractAction("Exit") {
		public void actionPerformed(ActionEvent event) {
			System.exit(0);
		}
	};

	Action connectLocal = new AbstractAction("Connect to localhost:7777") {
		public void actionPerformed(ActionEvent event) {
			System.exit(0);
		}
	};

	Action connect = new AbstractAction("Connect...") {
		public void actionPerformed(ActionEvent event) {
			System.exit(0);
		}
	};

	Action disconnect = new AbstractAction("Disconnect") {
		public void actionPerformed(ActionEvent event) {
			System.exit(0);
		}
	};

	Action execute = new AbstractAction("Execute") {
		public void actionPerformed(ActionEvent event) {
			System.exit(0);
		}
	};

	Action cancelExecute = new AbstractAction("Cancel Executing Query") {
		public void actionPerformed(ActionEvent event) {
			System.exit(0);
		}
	};

	Action simulate = new AbstractAction("Simulate") {
		public void actionPerformed(ActionEvent event) {
			System.exit(0);
		}
	};

	Action about = new AbstractAction("About") {
		public void actionPerformed(ActionEvent event) {
			System.exit(0);
		}
	};
	
	// helpers 

	private void defaultTooltips(JMenuItem menu) {
		menu.setToolTipText(menu.getText());
	}
}
