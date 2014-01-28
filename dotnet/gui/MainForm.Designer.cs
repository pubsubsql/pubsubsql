namespace PubSubSQLGUI
{
    partial class MainForm
    {
        /// <summary>
        /// Required designer variable.
        /// </summary>
        private System.ComponentModel.IContainer components = null;

        /// <summary>
        /// Clean up any resources being used.
        /// </summary>
        /// <param name="disposing">true if managed resources should be disposed; otherwise, false.</param>
        protected override void Dispose(bool disposing)
        {
            if (disposing && (components != null))
            {
                components.Dispose();
            }
            base.Dispose(disposing);
        }

        #region Windows Form Designer generated code

        /// <summary>
        /// Required method for Designer support - do not modify
        /// the contents of this method with the code editor.
        /// </summary>
        private void InitializeComponent()
        {
            System.ComponentModel.ComponentResourceManager resources = new System.ComponentModel.ComponentResourceManager(typeof(MainForm));
            this.menuStrip = new System.Windows.Forms.MenuStrip();
            this.fileMenu = new System.Windows.Forms.ToolStripMenuItem();
            this.newMenu = new System.Windows.Forms.ToolStripMenuItem();
            this.toolStripSeparator1 = new System.Windows.Forms.ToolStripSeparator();
            this.exitMenu = new System.Windows.Forms.ToolStripMenuItem();
            this.connectionMenu = new System.Windows.Forms.ToolStripMenuItem();
            this.connectLocalMenu = new System.Windows.Forms.ToolStripMenuItem();
            this.connectMenu = new System.Windows.Forms.ToolStripMenuItem();
            this.disconnectMenu = new System.Windows.Forms.ToolStripMenuItem();
            this.queryMenu = new System.Windows.Forms.ToolStripMenuItem();
            this.customToolStripMenuItem = new System.Windows.Forms.ToolStripSeparator();
            this.executeMenu = new System.Windows.Forms.ToolStripMenuItem();
            this.cancelMenu = new System.Windows.Forms.ToolStripMenuItem();
            this.simulateMenu = new System.Windows.Forms.ToolStripMenuItem();
            this.helpMenu = new System.Windows.Forms.ToolStripMenuItem();
            this.aboutMenu = new System.Windows.Forms.ToolStripMenuItem();
            this.toolStrip = new System.Windows.Forms.ToolStrip();
            this.newButton = new System.Windows.Forms.ToolStripButton();
            this.separator1 = new System.Windows.Forms.ToolStripSeparator();
            this.connectLocalButton = new System.Windows.Forms.ToolStripButton();
            this.connectButton = new System.Windows.Forms.ToolStripButton();
            this.disconnectButton = new System.Windows.Forms.ToolStripButton();
            this.separator2 = new System.Windows.Forms.ToolStripSeparator();
            this.executeButton = new System.Windows.Forms.ToolStripButton();
            this.cancelButton = new System.Windows.Forms.ToolStripButton();
            this.splitContainer = new System.Windows.Forms.SplitContainer();
            this.queryText = new System.Windows.Forms.TextBox();
            this.resultsTabContainer = new System.Windows.Forms.TabControl();
            this.resultsTab = new System.Windows.Forms.TabPage();
            this.statusTab = new System.Windows.Forms.TabPage();
            this.statusText = new System.Windows.Forms.TextBox();
            this.rawDataTab = new System.Windows.Forms.TabPage();
            this.jsonText = new System.Windows.Forms.TextBox();
            this.listView = new PubSubSQLGUI.ListViewDoubleBuffered();
            this.menuStrip.SuspendLayout();
            this.toolStrip.SuspendLayout();
            ((System.ComponentModel.ISupportInitialize)(this.splitContainer)).BeginInit();
            this.splitContainer.Panel1.SuspendLayout();
            this.splitContainer.Panel2.SuspendLayout();
            this.splitContainer.SuspendLayout();
            this.resultsTabContainer.SuspendLayout();
            this.resultsTab.SuspendLayout();
            this.statusTab.SuspendLayout();
            this.rawDataTab.SuspendLayout();
            this.SuspendLayout();
            // 
            // menuStrip
            // 
            this.menuStrip.Items.AddRange(new System.Windows.Forms.ToolStripItem[] {
            this.fileMenu,
            this.connectionMenu,
            this.queryMenu,
            this.helpMenu});
            this.menuStrip.Location = new System.Drawing.Point(0, 0);
            this.menuStrip.Name = "menuStrip";
            this.menuStrip.Size = new System.Drawing.Size(778, 24);
            this.menuStrip.TabIndex = 0;
            this.menuStrip.Text = "menuStrip";
            // 
            // fileMenu
            // 
            this.fileMenu.DropDownItems.AddRange(new System.Windows.Forms.ToolStripItem[] {
            this.newMenu,
            this.toolStripSeparator1,
            this.exitMenu});
            this.fileMenu.Name = "fileMenu";
            this.fileMenu.Size = new System.Drawing.Size(37, 20);
            this.fileMenu.Text = "File";
            // 
            // newMenu
            // 
            this.newMenu.Image = ((System.Drawing.Image)(resources.GetObject("newMenu.Image")));
            this.newMenu.ImageTransparentColor = System.Drawing.Color.Magenta;
            this.newMenu.Name = "newMenu";
            this.newMenu.Size = new System.Drawing.Size(152, 22);
            this.newMenu.Text = "New";
            // 
            // toolStripSeparator1
            // 
            this.toolStripSeparator1.Name = "toolStripSeparator1";
            this.toolStripSeparator1.Size = new System.Drawing.Size(149, 6);
            // 
            // exitMenu
            // 
            this.exitMenu.Name = "exitMenu";
            this.exitMenu.Size = new System.Drawing.Size(152, 22);
            this.exitMenu.Text = "Exit";
            // 
            // connectionMenu
            // 
            this.connectionMenu.DropDownItems.AddRange(new System.Windows.Forms.ToolStripItem[] {
            this.connectLocalMenu,
            this.connectMenu,
            this.disconnectMenu});
            this.connectionMenu.Name = "connectionMenu";
            this.connectionMenu.Size = new System.Drawing.Size(81, 20);
            this.connectionMenu.Text = "Connection";
            // 
            // connectLocalMenu
            // 
            this.connectLocalMenu.Image = ((System.Drawing.Image)(resources.GetObject("connectLocalMenu.Image")));
            this.connectLocalMenu.ImageTransparentColor = System.Drawing.Color.Magenta;
            this.connectLocalMenu.Name = "connectLocalMenu";
            this.connectLocalMenu.Size = new System.Drawing.Size(211, 22);
            this.connectLocalMenu.Text = "Connect to localhost:7777";
            // 
            // connectMenu
            // 
            this.connectMenu.Image = global::PubSubSQLGUI.Properties.Resources.Connect;
            this.connectMenu.ImageTransparentColor = System.Drawing.Color.Magenta;
            this.connectMenu.Name = "connectMenu";
            this.connectMenu.Size = new System.Drawing.Size(211, 22);
            this.connectMenu.Text = "Connect...";
            // 
            // disconnectMenu
            // 
            this.disconnectMenu.Image = ((System.Drawing.Image)(resources.GetObject("disconnectMenu.Image")));
            this.disconnectMenu.ImageTransparentColor = System.Drawing.Color.Magenta;
            this.disconnectMenu.Name = "disconnectMenu";
            this.disconnectMenu.Size = new System.Drawing.Size(211, 22);
            this.disconnectMenu.Text = "Disconnect";
            // 
            // queryMenu
            // 
            this.queryMenu.DropDownItems.AddRange(new System.Windows.Forms.ToolStripItem[] {
            this.customToolStripMenuItem,
            this.executeMenu,
            this.cancelMenu,
            this.simulateMenu});
            this.queryMenu.Name = "queryMenu";
            this.queryMenu.Size = new System.Drawing.Size(51, 20);
            this.queryMenu.Text = "Query";
            // 
            // customToolStripMenuItem
            // 
            this.customToolStripMenuItem.Name = "customToolStripMenuItem";
            this.customToolStripMenuItem.Size = new System.Drawing.Size(196, 6);
            // 
            // executeMenu
            // 
            this.executeMenu.Image = ((System.Drawing.Image)(resources.GetObject("executeMenu.Image")));
            this.executeMenu.ImageTransparentColor = System.Drawing.Color.Magenta;
            this.executeMenu.Name = "executeMenu";
            this.executeMenu.Size = new System.Drawing.Size(199, 22);
            this.executeMenu.Text = "Execute";
            // 
            // cancelMenu
            // 
            this.cancelMenu.Image = ((System.Drawing.Image)(resources.GetObject("cancelMenu.Image")));
            this.cancelMenu.ImageTransparentColor = System.Drawing.Color.Magenta;
            this.cancelMenu.Name = "cancelMenu";
            this.cancelMenu.Size = new System.Drawing.Size(199, 22);
            this.cancelMenu.Text = "Cancel Executing Query";
            // 
            // simulateMenu
            // 
            this.simulateMenu.Name = "simulateMenu";
            this.simulateMenu.Size = new System.Drawing.Size(199, 22);
            this.simulateMenu.Text = "Simulate";
            this.simulateMenu.ToolTipText = "Simulate";
            // 
            // helpMenu
            // 
            this.helpMenu.DropDownItems.AddRange(new System.Windows.Forms.ToolStripItem[] {
            this.aboutMenu});
            this.helpMenu.Name = "helpMenu";
            this.helpMenu.Size = new System.Drawing.Size(44, 20);
            this.helpMenu.Text = "Help";
            // 
            // aboutMenu
            // 
            this.aboutMenu.Name = "aboutMenu";
            this.aboutMenu.Size = new System.Drawing.Size(107, 22);
            this.aboutMenu.Text = "About";
            // 
            // toolStrip
            // 
            this.toolStrip.Items.AddRange(new System.Windows.Forms.ToolStripItem[] {
            this.newButton,
            this.separator1,
            this.connectLocalButton,
            this.connectButton,
            this.disconnectButton,
            this.separator2,
            this.executeButton,
            this.cancelButton});
            this.toolStrip.Location = new System.Drawing.Point(0, 24);
            this.toolStrip.Name = "toolStrip";
            this.toolStrip.Size = new System.Drawing.Size(778, 25);
            this.toolStrip.TabIndex = 1;
            this.toolStrip.Text = "toolStrip1";
            // 
            // newButton
            // 
            this.newButton.DisplayStyle = System.Windows.Forms.ToolStripItemDisplayStyle.Image;
            this.newButton.Image = ((System.Drawing.Image)(resources.GetObject("newButton.Image")));
            this.newButton.ImageTransparentColor = System.Drawing.Color.Magenta;
            this.newButton.Name = "newButton";
            this.newButton.Size = new System.Drawing.Size(23, 22);
            this.newButton.ToolTipText = "New PubSubSQL Interactive Query";
            // 
            // separator1
            // 
            this.separator1.Name = "separator1";
            this.separator1.Size = new System.Drawing.Size(6, 25);
            // 
            // connectLocalButton
            // 
            this.connectLocalButton.DisplayStyle = System.Windows.Forms.ToolStripItemDisplayStyle.Image;
            this.connectLocalButton.Image = ((System.Drawing.Image)(resources.GetObject("connectLocalButton.Image")));
            this.connectLocalButton.ImageTransparentColor = System.Drawing.Color.Magenta;
            this.connectLocalButton.Name = "connectLocalButton";
            this.connectLocalButton.Size = new System.Drawing.Size(23, 22);
            this.connectLocalButton.Text = "Connect to localhost:7777";
            // 
            // connectButton
            // 
            this.connectButton.DisplayStyle = System.Windows.Forms.ToolStripItemDisplayStyle.Image;
            this.connectButton.Image = ((System.Drawing.Image)(resources.GetObject("connectButton.Image")));
            this.connectButton.ImageTransparentColor = System.Drawing.Color.Magenta;
            this.connectButton.Name = "connectButton";
            this.connectButton.Size = new System.Drawing.Size(23, 22);
            this.connectButton.Text = "Connect to remote server";
            // 
            // disconnectButton
            // 
            this.disconnectButton.DisplayStyle = System.Windows.Forms.ToolStripItemDisplayStyle.Image;
            this.disconnectButton.Image = ((System.Drawing.Image)(resources.GetObject("disconnectButton.Image")));
            this.disconnectButton.ImageTransparentColor = System.Drawing.Color.Magenta;
            this.disconnectButton.Name = "disconnectButton";
            this.disconnectButton.Size = new System.Drawing.Size(23, 22);
            this.disconnectButton.Text = "Disconnect";
            // 
            // separator2
            // 
            this.separator2.Name = "separator2";
            this.separator2.Size = new System.Drawing.Size(6, 25);
            // 
            // executeButton
            // 
            this.executeButton.Image = ((System.Drawing.Image)(resources.GetObject("executeButton.Image")));
            this.executeButton.ImageTransparentColor = System.Drawing.Color.Magenta;
            this.executeButton.Name = "executeButton";
            this.executeButton.Size = new System.Drawing.Size(70, 22);
            this.executeButton.Text = " Execute";
            // 
            // cancelButton
            // 
            this.cancelButton.Image = ((System.Drawing.Image)(resources.GetObject("cancelButton.Image")));
            this.cancelButton.ImageTransparentColor = System.Drawing.Color.Magenta;
            this.cancelButton.Name = "cancelButton";
            this.cancelButton.Size = new System.Drawing.Size(23, 22);
            this.cancelButton.ToolTipText = "Cancel Executing Query";
            // 
            // splitContainer
            // 
            this.splitContainer.Dock = System.Windows.Forms.DockStyle.Fill;
            this.splitContainer.Location = new System.Drawing.Point(0, 49);
            this.splitContainer.Name = "splitContainer";
            this.splitContainer.Orientation = System.Windows.Forms.Orientation.Horizontal;
            // 
            // splitContainer.Panel1
            // 
            this.splitContainer.Panel1.Controls.Add(this.queryText);
            // 
            // splitContainer.Panel2
            // 
            this.splitContainer.Panel2.Controls.Add(this.resultsTabContainer);
            this.splitContainer.Size = new System.Drawing.Size(778, 516);
            this.splitContainer.SplitterDistance = 133;
            this.splitContainer.TabIndex = 2;
            // 
            // queryText
            // 
            this.queryText.Dock = System.Windows.Forms.DockStyle.Fill;
            this.queryText.Font = new System.Drawing.Font("Verdana", 9.75F, System.Drawing.FontStyle.Regular, System.Drawing.GraphicsUnit.Point, ((byte)(0)));
            this.queryText.ForeColor = System.Drawing.Color.MediumBlue;
            this.queryText.Location = new System.Drawing.Point(0, 0);
            this.queryText.Multiline = true;
            this.queryText.Name = "queryText";
            this.queryText.Size = new System.Drawing.Size(778, 133);
            this.queryText.TabIndex = 0;
            // 
            // resultsTabContainer
            // 
            this.resultsTabContainer.Controls.Add(this.resultsTab);
            this.resultsTabContainer.Controls.Add(this.statusTab);
            this.resultsTabContainer.Controls.Add(this.rawDataTab);
            this.resultsTabContainer.Dock = System.Windows.Forms.DockStyle.Fill;
            this.resultsTabContainer.Location = new System.Drawing.Point(0, 0);
            this.resultsTabContainer.Name = "resultsTabContainer";
            this.resultsTabContainer.SelectedIndex = 0;
            this.resultsTabContainer.Size = new System.Drawing.Size(778, 379);
            this.resultsTabContainer.TabIndex = 0;
            // 
            // resultsTab
            // 
            this.resultsTab.Controls.Add(this.listView);
            this.resultsTab.Location = new System.Drawing.Point(4, 22);
            this.resultsTab.Name = "resultsTab";
            this.resultsTab.Padding = new System.Windows.Forms.Padding(3);
            this.resultsTab.Size = new System.Drawing.Size(770, 353);
            this.resultsTab.TabIndex = 0;
            this.resultsTab.Text = "Results";
            this.resultsTab.UseVisualStyleBackColor = true;
            // 
            // statusTab
            // 
            this.statusTab.Controls.Add(this.statusText);
            this.statusTab.Location = new System.Drawing.Point(4, 22);
            this.statusTab.Name = "statusTab";
            this.statusTab.Padding = new System.Windows.Forms.Padding(3);
            this.statusTab.Size = new System.Drawing.Size(770, 353);
            this.statusTab.TabIndex = 1;
            this.statusTab.Text = "Status";
            this.statusTab.UseVisualStyleBackColor = true;
            // 
            // statusText
            // 
            this.statusText.BackColor = System.Drawing.SystemColors.ControlLightLight;
            this.statusText.Dock = System.Windows.Forms.DockStyle.Fill;
            this.statusText.Font = new System.Drawing.Font("Verdana", 9.75F, System.Drawing.FontStyle.Regular, System.Drawing.GraphicsUnit.Point, ((byte)(0)));
            this.statusText.Location = new System.Drawing.Point(3, 3);
            this.statusText.Multiline = true;
            this.statusText.Name = "statusText";
            this.statusText.ReadOnly = true;
            this.statusText.Size = new System.Drawing.Size(764, 347);
            this.statusText.TabIndex = 0;
            // 
            // rawDataTab
            // 
            this.rawDataTab.Controls.Add(this.jsonText);
            this.rawDataTab.Location = new System.Drawing.Point(4, 22);
            this.rawDataTab.Name = "rawDataTab";
            this.rawDataTab.Padding = new System.Windows.Forms.Padding(3);
            this.rawDataTab.Size = new System.Drawing.Size(770, 353);
            this.rawDataTab.TabIndex = 2;
            this.rawDataTab.Text = "JSON Response";
            this.rawDataTab.UseVisualStyleBackColor = true;
            // 
            // rawdataText
            // 
            this.jsonText.Dock = System.Windows.Forms.DockStyle.Fill;
            this.jsonText.Font = new System.Drawing.Font("Verdana", 9.75F, System.Drawing.FontStyle.Regular, System.Drawing.GraphicsUnit.Point, ((byte)(0)));
            this.jsonText.Location = new System.Drawing.Point(3, 3);
            this.jsonText.Multiline = true;
            this.jsonText.Name = "rawdataText";
            this.jsonText.Size = new System.Drawing.Size(764, 347);
            this.jsonText.TabIndex = 0;
            // 
            // listView
            // 
            this.listView.Activation = System.Windows.Forms.ItemActivation.OneClick;
            this.listView.Dock = System.Windows.Forms.DockStyle.Fill;
            this.listView.GridLines = true;
            this.listView.HeaderStyle = System.Windows.Forms.ColumnHeaderStyle.Nonclickable;
            this.listView.HoverSelection = true;
            this.listView.Location = new System.Drawing.Point(3, 3);
            this.listView.MultiSelect = false;
            this.listView.Name = "listView";
            this.listView.Size = new System.Drawing.Size(764, 347);
            this.listView.TabIndex = 0;
            this.listView.UseCompatibleStateImageBehavior = false;
            this.listView.View = System.Windows.Forms.View.Details;
            this.listView.VirtualMode = true;
            this.listView.RetrieveVirtualItem += new System.Windows.Forms.RetrieveVirtualItemEventHandler(this.listView_RetrieveVirtualItem);
            // 
            // MainForm
            // 
            this.AutoScaleDimensions = new System.Drawing.SizeF(6F, 13F);
            this.AutoScaleMode = System.Windows.Forms.AutoScaleMode.Font;
            this.ClientSize = new System.Drawing.Size(778, 565);
            this.Controls.Add(this.splitContainer);
            this.Controls.Add(this.toolStrip);
            this.Controls.Add(this.menuStrip);
            this.Icon = ((System.Drawing.Icon)(resources.GetObject("$this.Icon")));
            this.MainMenuStrip = this.menuStrip;
            this.Name = "MainForm";
            this.Text = "PubSubSQL Interactive Query";
            this.FormClosing += new System.Windows.Forms.FormClosingEventHandler(this.MainForm_FormClosing);
            this.Load += new System.EventHandler(this.MainForm_Load);
            this.menuStrip.ResumeLayout(false);
            this.menuStrip.PerformLayout();
            this.toolStrip.ResumeLayout(false);
            this.toolStrip.PerformLayout();
            this.splitContainer.Panel1.ResumeLayout(false);
            this.splitContainer.Panel1.PerformLayout();
            this.splitContainer.Panel2.ResumeLayout(false);
            ((System.ComponentModel.ISupportInitialize)(this.splitContainer)).EndInit();
            this.splitContainer.ResumeLayout(false);
            this.resultsTabContainer.ResumeLayout(false);
            this.resultsTab.ResumeLayout(false);
            this.statusTab.ResumeLayout(false);
            this.statusTab.PerformLayout();
            this.rawDataTab.ResumeLayout(false);
            this.rawDataTab.PerformLayout();
            this.ResumeLayout(false);
            this.PerformLayout();

        }

        #endregion

        private System.Windows.Forms.MenuStrip menuStrip;
        private System.Windows.Forms.ToolStripMenuItem fileMenu;
        private System.Windows.Forms.ToolStripMenuItem newMenu;
        private System.Windows.Forms.ToolStripSeparator toolStripSeparator1;
        private System.Windows.Forms.ToolStripMenuItem exitMenu;
        private System.Windows.Forms.ToolStripMenuItem queryMenu;
        private System.Windows.Forms.ToolStripSeparator customToolStripMenuItem;
        private System.Windows.Forms.ToolStripMenuItem executeMenu;
        private System.Windows.Forms.ToolStripMenuItem cancelMenu;
        private System.Windows.Forms.ToolStripMenuItem helpMenu;
        private System.Windows.Forms.ToolStripMenuItem aboutMenu;
        private System.Windows.Forms.ToolStrip toolStrip;
        private System.Windows.Forms.ToolStripButton connectButton;
        private System.Windows.Forms.ToolStripButton connectLocalButton;
        private System.Windows.Forms.ToolStripButton disconnectButton;
        private System.Windows.Forms.ToolStripButton newButton;
        private System.Windows.Forms.ToolStripSeparator separator1;
        private System.Windows.Forms.ToolStripSeparator separator2;
        private System.Windows.Forms.ToolStripButton executeButton;
        private System.Windows.Forms.ToolStripButton cancelButton;
        private System.Windows.Forms.ToolStripMenuItem connectionMenu;
        private System.Windows.Forms.ToolStripMenuItem connectLocalMenu;
        private System.Windows.Forms.ToolStripMenuItem connectMenu;
        private System.Windows.Forms.ToolStripMenuItem disconnectMenu;
        private System.Windows.Forms.SplitContainer splitContainer;
        private System.Windows.Forms.TextBox queryText;
        private System.Windows.Forms.TabControl resultsTabContainer;
        private System.Windows.Forms.TabPage resultsTab;
        private System.Windows.Forms.TabPage statusTab;
        private System.Windows.Forms.TabPage rawDataTab;
        private System.Windows.Forms.TextBox statusText;
        private System.Windows.Forms.TextBox jsonText;
        private ListViewDoubleBuffered listView;
        private System.Windows.Forms.ToolStripMenuItem simulateMenu;
    }
}

