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
        private bool cancelExecuteFlag = false;

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

            setTitle(string.Empty);
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
            if (client.Connect(address)) setTitle(address);
            setStatus();
        }

        private void disconnect(object sender, EventArgs e)
        {
            setTitle(string.Empty);
            clear();
            client.Disconnect();
        }

        private void execute(object sender, EventArgs e)
        {
            try
            {
                executing();
                cancelExecuteFlag = false;
                string command = queryText.Text.Trim();
                if (string.IsNullOrEmpty(command)) return;
                client.Execute(queryText.Text);
                processResults();
            }
            finally
            {
                doneExecuting();
            }
        }

        private void cancelExecute(object sender, EventArgs e)
        {
            cancelExecuteFlag = true;
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
            rawdataText.Text = "";
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

        private void setRawData()
        {
            rawdataText.Text = client.JSON();
        }

        private void setTitle(string address)
        {
            Text = "PubSubSQL Interactive Query " + address;
        }

        private void executing()
        {        
            queryText.Enabled = false;
            executeButton.Enabled = false;
            executeMenu.Enabled = false;
            cancelButton.Enabled = true;
            cancelMenu.Enabled = true;
        }

        private void doneExecuting()
        {
            queryText.Enabled = true;
            executeButton.Enabled = true;
            executeMenu.Enabled = true;
            cancelButton.Enabled = false;
            cancelMenu.Enabled = false;
        }

        private void processResults()
        {
            setStatus();
            setRawData();
        }
    }
}
