using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Data;
using System.Drawing;
using System.Linq;
using System.Text;
using System.Windows.Forms;

namespace PubSubSQLGUI
{
    public partial class MainForm : Form
    {
        private PubSubSQL.Client client = PubSubSQL.Factory.NewClient();
        private string DEFAULT_ADDRESS = "localhost:7777";

        public MainForm()
        {
            InitializeComponent();
            // set up controls and events
            setControls(newButton, newMenu, new_);
            exitMenu.Click += exit;

            connectLocalButton.ToolTipText = "Connect to " + DEFAULT_ADDRESS;
            setControls(connectLocalButton, connectLocalMenu, connectLocal);
            connectLocalMenu.Text = connectLocalMenu.ToolTipText;

            setControls(connectButton, connectMenu, connect);
            setControls(disconnectButton, disconnectMenu, disconnect);
            setControls(executeButton, executeMenu, execute);
            setControls(cancelButton, cancelMenu, cancelExecute);

            nextPaneMenu.Click += nextPane;
            aboutMenu.Click += about;
        }

        private void setControls(ToolStripButton button, ToolStripMenuItem menu, EventHandler click)
        {
            menu.ToolTipText = button.ToolTipText;
            menu.Click += click;
            button.Click += click;
        }

        // gui events

        private void new_(object sender, EventArgs e)
        {
            this.Close();
        }

        private void exit(object sender, EventArgs e)
        {
            this.Close();
        }

        private void connectLocal(object sender, EventArgs e)
        {
            connect(DEFAULT_ADDRESS); 
        }

        private void connect(object sender, EventArgs e)
        {
            
        }

        private void connect(string address)
        {
            clear();
            client.Connect(address);
            setStatus();
        }

        private void disconnect(object sender, EventArgs e)
        {
            clear();
            client.Disconnect();
        }

        private void execute(object sender, EventArgs e)
        {

        }

        private void cancelExecute(object sender, EventArgs e)
        {

        }

        private void nextPane(object sender, EventArgs e)
        {

        }

        private void about(object sender, EventArgs e)
        {

        }

        // helper functions

        private void clear()
        {
            statusText.Text = "";
        }

        private bool setStatus()
        {
            resultsTabContainer.SelectedTab = statusTab;
            if (client.Ok())
            {
                statusText.ForeColor = Color.Black;
                statusText.Text = "ok";
                return true;
            }
            statusText.ForeColor = Color.Red;
            statusText.Text = "error\r\n" + client.Error();
            return false;
        }
    }
}
