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
        private string DEFAULT_ADDRESS = "localhost:7777";
        private int PUBSUB_TIMEOUT = 5;
        private PubSubSQL.Client client = PubSubSQL.Factory.NewClient();
        private bool cancelExecuteFlag = false;
        private string connectedAddress = string.Empty;
        private ListViewDataset dataset = new ListViewDataset();

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
            resultsTabContainer.SelectedTab = statusTab;
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
            if (client.Connect(address))
            {
                setTitle(address);
            }
            setStatus();
        }

        private void disconnect(object sender, EventArgs e)
        {
            setTitle(string.Empty);
            clear();
            client.Disconnect();
            connectedAddress = string.Empty;
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
                processResponse();
            }
            finally
            {
                doneExecuting();
                // we were stoped in the middle
                if (cancelExecuteFlag)
                {
                    // lets make it simple and just reconnect
                    // done in order to ignore possible pubsub event or not fully read result set
                    if (!string.IsNullOrEmpty(connectedAddress))
                    {
                        connect(connectedAddress);
                        clear();
                    }
                }
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
            listView.Columns.Clear();
            listView.VirtualListSize = 0;
            statusText.Text = "";
            rawdataText.Text = "";
        }

        private bool setStatus()
        {
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
            connectedAddress = address;
        }

        private void executing()
        {        
            clear();
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

        private void processResponse()
        {
            setStatus();
            setRawData();
            dataset.Reset();
            if (client.Failed()) return;
            // determine if we just subscribed  
            if (client.PubSubId() != string.Empty && client.Action() == "subscribe")
            {
                waitForPubSubEvent();
                return;
            }
            processResults();
        }

        private void syncColumns()
        {
            for (int i = listView.Columns.Count; i < dataset.ColumnCount; i++)
            {
                string col = dataset.Column(i); 
                listView.Columns.Add(col);
                listView.Columns[i].Width = 100;
            }
        }

        private void waitForPubSubEvent()
        {
            while (client.Ok() && !cancelExecuteFlag)
            {
                Application.DoEvents();
                if (client.WaitForPubSub(PUBSUB_TIMEOUT) && client.Ok())
                {
                    setStatus();
                    setRawData();
                    processResults();
                }
            }
        }

        private void processResults()
        {
            setRawData();
            bool results = false;
            // check if it is result set
            if (client.RecordCount() > 0 && client.ColumnCount() > 0)
            {
                // inside dataset
                dataset.SyncColumns(client);
                syncColumns();
                results = true;
                dataset.AddRowsCapacity(client.RecordCount());
                while (client.NextRecord() && !cancelExecuteFlag)
                {
                    dataset.ProcessRow(client);
                    listView.VirtualListSize = dataset.RowCount;        
                    Application.DoEvents();
                }
            }
            if (client.Failed())
            {
                setStatus();
                setRawData();
                resultsTabContainer.SelectedTab = statusTab;
            }
            else
            {
                if (results)
                {
                    resultsTabContainer.SelectedTab = resultsTab;
                    listView.BeginUpdate();
                    listView.EndUpdate();
                    Application.DoEvents();
                }
                else
                {
                    resultsTabContainer.SelectedTab = statusTab;
                }
            }
        }

        private void listView_RetrieveVirtualItem(object sender, RetrieveVirtualItemEventArgs e)
        {
            List<string> row = dataset.GetRow(e.ItemIndex);
            ListViewItem item = new ListViewItem();
            for (int i = 0; i < listView.Columns.Count; i++)
            {
                string str = string.Empty;
                if (i < row.Count) str = row[i];
                if (i == 0) item.Text = str;
                else item.SubItems.Add(str);
            }
            e.Item = item;
        }

        private void MainForm_FormClosing(object sender, FormClosingEventArgs e)
        {
            cancelExecuteFlag = true;
        }
    }
}
