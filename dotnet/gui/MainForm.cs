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
        private bool useFlashColor = false;
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
            setTitle(string.Empty);
            resultsTabContainer.SelectedTab = statusTab;
            enableDisableControls();

            flashTimer.Interval = FLASH_TIMER_INTERVAL;
            flashTimer.Enabled = false;
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
            enableDisableControls();
        }

        private void disconnect(object sender, EventArgs e)
        {
            setTitle(string.Empty);
            clear();
            client.Disconnect();
            connectedAddress = string.Empty;
            enableDisableControls();
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
            simulator.Stop();
            cancelExecuteFlag = true;
        }

        private void simulate(object sender, EventArgs e)
        {
            SimulatorForm dlg = new SimulatorForm();
            if (dlg.ShowDialog() == DialogResult.OK)
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

        }

        private void tick(object sender, EventArgs e)
        {
            listView.BeginUpdate();
            listView.EndUpdate();
            if (!useFlashColor) flashTimer.Enabled = false;
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
            enableDisableControls();
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
            bool connected = client.Connected();
            queryText.Enabled = true;
            executeButton.Enabled = connected;
            executeMenu.Enabled = connected;
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
                try
                {
                    useFlashColor = true;
                    flashTimer.Enabled = true;
                    waitForPubSubEvent();
                }
                finally
                {
                    useFlashColor = false;
                }
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
            if (client.Failed())
            {
                setStatus();
            }
        }

        private int updateDataset()
        {
            int times = 0;
            // inside dataset
            dataset.SyncColumns(client);
            syncColumns();
            dataset.AddRowsCapacity(client.RecordCount());
            while (client.NextRecord() && !cancelExecuteFlag)
            {
                times++;
                dataset.ProcessRow(client);
            }
            return times;
        }

        private void processResults()
        {
            setRawData();
            bool results = false;
            // check if it is result set
            if (client.RecordCount() > 0 && client.ColumnCount() > 0)
            {
                results = true;
                updateDataset(); 
                listView.VirtualListSize = dataset.RowCount;        
                Application.DoEvents();
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
                    if (useFlashColor && (DateTime.Now.Ticks - cell.LastUpdated) < FLASH_TIMEOUT)
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
        }
    }
}
