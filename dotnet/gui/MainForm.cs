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
        private int PUBSUB_TIMEOUT = 5; // in milliseconds
        private int FLASH_TIMER_INTERVAL = 150;
        private long FLASH_TIMEOUT = 300 * 10000; // ticks converted to milliseconds
        private PubSubSQL.Client client = PubSubSQL.Factory.NewClient();
        private bool cancelExecuteFlag = false;
        private string connectedAddress = string.Empty;
        private ListViewDataset dataset = new ListViewDataset();
        private Timer flashTimer = new Timer();
        private Simulator simulator = new Simulator();

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
            simulateMenu.Click += simulate;
            aboutMenu.Click += about;
            updateConnectedAddress(string.Empty);
            resultsTabContainer.SelectedTab = statusTab;
            enableDisableControls();

            flashTimer.Interval = FLASH_TIMER_INTERVAL;
            flashTimer.Enabled = true;
            flashTimer.Tick += tick;
        }

        private void enableDisableControls()
        {
            bool connected = client.Connected();
            connectLocalButton.Enabled = !connected;
            connectLocalMenu.Enabled = !connected;
            connectButton.Enabled = !connected;
            connectMenu.Enabled = !connected;
            disconnectButton.Enabled = connected;
            disconnectMenu.Enabled = connected;
            executeButton.Enabled = connected;
            executeMenu.Enabled = connected;
            simulateMenu.Enabled = executeMenu.Enabled;
            cancelButton.Enabled = false;
            cancelMenu.Enabled = false;
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
            System.Diagnostics.Process.Start(Application.ExecutablePath, connectedAddress); 
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
            ConnectForm dlg = new ConnectForm();
            if (dlg.ShowDialog(this) == System.Windows.Forms.DialogResult.OK)
            {
                connect(dlg.Address);
            }
        }

        private void connect(string address)
        {
            clearResults();
            if (client.Connect(address))
            {
                updateConnectedAddress(address);
            }
            setStatus();
            enableDisableControls();
        }

        private void disconnect(object sender, EventArgs e)
        {
            simulator.Stop();
            cancelExecuteFlag = true;
            updateConnectedAddress(string.Empty);
            client.Disconnect();
            enableDisableControls();
            clearResults();
        }

        private void execute(object sender, EventArgs e)
        {
            try
            {
                executing();
                string command = queryText.Text.Trim();
                if (string.IsNullOrEmpty(command)) return;
                client.Execute(command);
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
                        clearResults();
                    }
                }
            }
        }

        private void cancelExecute(object sender, EventArgs e)
        {
            simulator.Stop();
            cancelExecuteFlag = true;
        }

        private void simulate(object sender, EventArgs e)
        {
            SimulatorForm dlg = new SimulatorForm();
            if (dlg.ShowDialog(this) == DialogResult.OK)
            {
                simulator.Stop();
                simulator.Address = connectedAddress;
                simulator.Columns = dlg.Columns;
                simulator.Rows = dlg.Rows;
                simulator.TableName = dlg.TableName;
                simulator.Start();
                queryText.Text = "subscribe * from " + dlg.TableName;
                execute(sender, e);
            }
        }

        private void about(object sender, EventArgs e)
        {
            AboutForm dlg = new AboutForm();
            dlg.ShowDialog(this);
        }

        long flashTicks = DateTime.Now.Ticks;
        private void tick(object sender, EventArgs e)
        {
            if (dataset.ResetDirty())
            {
                setStatus();
                setJSON();
                listView.VirtualListSize = dataset.RowCount;
                resultsTabContainer.SelectedTab = resultsTab;
                listView.BeginUpdate();
                listView.EndUpdate();
                flashTicks = DateTime.Now.Ticks;
            }
            else if ((DateTime.Now.Ticks - flashTicks) < (FLASH_TIMEOUT * 2) )
            {
                // make sure that we clear the background color
                listView.VirtualListSize = dataset.RowCount;
                listView.BeginUpdate();
                listView.EndUpdate();
            }
        }

        // helper functions

        private void clearResults()
        {
            listView.Columns.Clear();
            listView.VirtualListSize = 0;
            statusText.Text = "";
            jsonText.Text = "";
        }

        private void setStatus()
        {
            if (client.Ok())
            {
                statusText.ForeColor = Color.Black;
                statusText.Text = "ok";
                return;
            }
            statusText.ForeColor = Color.Red;
            statusText.Text = "error\r\n" + client.Error();
            enableDisableControls();
            return;
        }

        private void setJSON()
        {
            jsonText.Text = client.JSON();
        }
        
        private void updateConnectedAddress(string address)
        {
            Text = "PubSubSQL Interactive Query " + address;
            connectedAddress = address;
        }

        private void executing()
        {        
            clearResults();
            cancelExecuteFlag = false;
            queryText.Enabled = false;
            executeButton.Enabled = false;
            executeMenu.Enabled = false;
            simulateMenu.Enabled = false;
            cancelButton.Enabled = true;
            cancelMenu.Enabled = true;
        }

        private void doneExecuting()
        {
            queryText.Enabled = true;
            enableDisableControls();
        }

        private void processResponse()
        {
            dataset.Clear();
            // determine if we just subscribed  
            if (client.PubSubId() != string.Empty && client.Action() == "subscribe")
            {
                setStatus();
                setJSON();
                // enter event loop
                waitForPubSubEvent();
                return;
            }
            // check if it is result set
            if (client.RowCount() > 0 && client.ColumnCount() > 0)
            {
                updateDataset(); 
                resultsTabContainer.SelectedTab = resultsTab;
            }
            //            
            if (client.Failed())resultsTabContainer.SelectedTab = statusTab;
            setStatus();
            setJSON();
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
            while (!cancelExecuteFlag)
            {
                bool timedout = !client.WaitForPubSub(PUBSUB_TIMEOUT);
                if (client.Failed()) break;
                if (!timedout) updateDataset();
                Application.DoEvents();
            }
            if (client.Failed()) setStatus();
        }

        private void updateDataset()
        {
            if (!(client.RowCount() > 0 && client.ColumnCount() > 0)) return;
            // inside dataset
            dataset.SyncColumns(client);
            syncColumns();
            dataset.AddRowsCapacity(client.RowCount());
            while (client.NextRow() && !cancelExecuteFlag)
            {
                dataset.ProcessRow(client);
            }
        }

        private void listView_RetrieveVirtualItem(object sender, RetrieveVirtualItemEventArgs e)
        {
            List<Cell> row = dataset.GetRow(e.ItemIndex);
            ListViewItem item = new ListViewItem();
            item.UseItemStyleForSubItems = false;
            for (int i = 0; i < listView.Columns.Count; i++)
            {
                string str = string.Empty;
                Color cellBackColor = Color.White;
                if (i < row.Count)
                {
                    Cell cell = row[i];
                    str = cell.Value;
                    if ((DateTime.Now.Ticks - cell.LastUpdated) < FLASH_TIMEOUT)
                    {
                        cellBackColor = Color.HotPink;
                    }
                }
                if (i == 0)
                {
                    item.Text = str;
                }
                else
                {
                    item.SubItems.Add(str).BackColor = cellBackColor;
                }
            }
            e.Item = item;
        }

        private void MainForm_FormClosing(object sender, FormClosingEventArgs e)
        {
            cancelExecuteFlag = true;
            simulator.Stop();
            client.Disconnect();
        }

        private void MainForm_Load(object sender, EventArgs e)
        {
            string[] args = Environment.GetCommandLineArgs();
            if (args.Length > 1)
            {
                connect(args[1]);
            }
        }

    }
}
