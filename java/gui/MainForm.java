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

	private JTextArea queryText;
	private JTabbedPane resultsTabContainer;
	private JTextArea statusText;
	private JTextArea jsonText;

	public MainForm() {
		setTitle("Interactive Query");
		setupMenuAndToolBar();		
		// query text
		queryText = new JTextArea();
		queryText.setPreferredSize(new Dimension(100, 100));
		// tabs
		resultsTabContainer = new JTabbedPane();
		statusText = new JTextArea();		
		resultsTabContainer.addTab("Status", statusText);
		jsonText = new JTextArea();
		resultsTabContainer.addTab("JSON Response", jsonText);
		// splitter
		JSplitPane splitPane = new JSplitPane(JSplitPane.VERTICAL_SPLIT, queryText, resultsTabContainer); 
		this.add(splitPane, BorderLayout.CENTER);	
		pack();
	}

	void setupMenuAndToolBar() {
		JMenuBar menuBar = new JMenuBar();
		this.setJMenuBar(menuBar);	
		JToolBar toolBar = new JToolBar();
		this.add(toolBar, BorderLayout.NORTH);
		// File
		JMenu fileMenu = new JMenu("File");
			// New
			JMenuItem newMenu = new JMenuItem(new_);	
			new_.putValue(Action.SHORT_DESCRIPTION, "New PubSubSQL Interactive Query");
			fileMenu.add(newMenu);
			fileMenu.addSeparator();
			toolBar.add(new_);
			toolBar.addSeparator();
			// Exit	
			JMenuItem exitMenu = new JMenuItem(exit);
			defaultTooltips(exit);
			fileMenu.add(exitMenu);
		menuBar.add(fileMenu);
		// Connection
		JMenu connectionMenu = new JMenu("Connection");
			// Connect local
			JMenuItem connectLocalMenu = new JMenuItem(connectLocal);
			defaultTooltips(connectLocal);
			connectionMenu.add(connectLocalMenu);
			toolBar.add(connectLocal);
			// Connect
			JMenuItem connectMenu = new JMenuItem(connect);
			connect.putValue(Action.SHORT_DESCRIPTION, "Connect to remote server");
			connectionMenu.add(connectMenu);
			toolBar.add(connect);
			// Disconnect
			JMenuItem disconnectMenu = new JMenuItem(disconnect);
			defaultTooltips(disconnect);
			connectionMenu.add(disconnectMenu);
			toolBar.add(disconnect);
			toolBar.addSeparator();
		menuBar.add(connectionMenu);	
		// Query
		JMenu queryMenu = new JMenu("Query");
			// Execute 
			JMenuItem executeMenu = new JMenuItem(execute);
			defaultTooltips(execute);
			queryMenu.add(executeMenu);
			toolBar.add(execute);
			// Cancel Executing Query 
			JMenuItem cancelMenu = new JMenuItem(cancelExecute);
			defaultTooltips(cancelExecute);
			queryMenu.add(cancelMenu);
			toolBar.add(cancelExecute);
			// Simulate 
			JMenuItem simulateMenu = new JMenuItem(simulate);
			defaultTooltips(simulate);
			queryMenu.add(simulateMenu);
		menuBar.add(queryMenu);	
		// Help
		JMenu helpMenu = new JMenu("Help");
			// About 
			JMenuItem aboutMenu = new JMenuItem(about);
			defaultTooltips(about);
			helpMenu.add(aboutMenu);
		menuBar.add(helpMenu);	
	}

	// events
	Action new_ = new AbstractAction("New", createImageIcon("images/New.png")) {
		public void actionPerformed(ActionEvent event) {
			
		}
	};
	
	Action exit = new AbstractAction("Exit") {
		public void actionPerformed(ActionEvent event) {
			System.exit(0);
		}
	};

	Action connectLocal = new AbstractAction("Connect to localhost:7777", createImageIcon("images/ConnectLocal.png")) {
		public void actionPerformed(ActionEvent event) {
			System.exit(0);
		}
	};

	Action connect = new AbstractAction("Connect...", createImageIcon("images/Connect.png")) {
		public void actionPerformed(ActionEvent event) {
			System.exit(0);
		}
	};

	Action disconnect = new AbstractAction("Disconnect", createImageIcon("images/Disconnect.png")) {
		public void actionPerformed(ActionEvent event) {
			System.exit(0);
		}
	};

	Action execute = new AbstractAction("Execute", createImageIcon("images/Execute2.png")) {
		public void actionPerformed(ActionEvent event) {
			System.exit(0);
		}
	};

	Action cancelExecute = new AbstractAction("Cancel Executing Query", createImageIcon("images/Stop.png")) {
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
	private void defaultTooltips(Action action) {
		action.putValue(Action.SHORT_DESCRIPTION, action.getValue(Action.NAME)); 
	}

	private ImageIcon createImageIcon(String path) {
		java.net.URL url = getClass().getResource(path);
		if (url == null) return null;
		return new ImageIcon(url);
	}
}
