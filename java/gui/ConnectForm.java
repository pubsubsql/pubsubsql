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

public class ConnectForm extends JDialog {

	private ConnectPanel panel;
	private boolean ok = false;

	public ConnectForm(JFrame owner) {
		super(owner, "Connect", true);		

		panel = new ConnectPanel();
		add(panel, BorderLayout.CENTER);
		pack();
		setResizable(false);
		
		panel.okButton.addActionListener( new ActionListener() {	
			public void actionPerformed(ActionEvent event) {
				ok = true;
				setVisible(false);
			}
		});

		panel.cancelButton.addActionListener( new ActionListener() {	
			public void actionPerformed(ActionEvent event) {
				ok = false;
				setVisible(false);
			}
		});
	}

	public void setHost(String host) {
		panel.hostText.setText(host);			
	}

	public void setPort(int port) {
		panel.portSpinner.setValue(port);		
	}

	public String getAddress() {
		return String.format("%s:%s", panel.hostText.getText(), panel.portSpinner.getValue()); 
	}	

	public boolean Ok() {
		return ok;
	}

}

