namespace PubSubSQLGUI
{
    partial class ConnectForm
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
            this.groupBox = new System.Windows.Forms.GroupBox();
            this.hostText = new System.Windows.Forms.TextBox();
            this.line = new System.Windows.Forms.GroupBox();
            this.cancelButton = new System.Windows.Forms.Button();
            this.okButton = new System.Windows.Forms.Button();
            this.portUpDown = new System.Windows.Forms.NumericUpDown();
            this.portLabel = new System.Windows.Forms.Label();
            this.hostLabel = new System.Windows.Forms.Label();
            this.groupBox.SuspendLayout();
            ((System.ComponentModel.ISupportInitialize)(this.portUpDown)).BeginInit();
            this.SuspendLayout();
            // 
            // groupBox
            // 
            this.groupBox.Controls.Add(this.hostText);
            this.groupBox.Controls.Add(this.line);
            this.groupBox.Controls.Add(this.cancelButton);
            this.groupBox.Controls.Add(this.okButton);
            this.groupBox.Controls.Add(this.portUpDown);
            this.groupBox.Controls.Add(this.portLabel);
            this.groupBox.Controls.Add(this.hostLabel);
            this.groupBox.Location = new System.Drawing.Point(7, -1);
            this.groupBox.Margin = new System.Windows.Forms.Padding(4, 4, 4, 4);
            this.groupBox.Name = "groupBox";
            this.groupBox.Padding = new System.Windows.Forms.Padding(4, 4, 4, 4);
            this.groupBox.Size = new System.Drawing.Size(303, 185);
            this.groupBox.TabIndex = 27;
            this.groupBox.TabStop = false;
            // 
            // hostText
            // 
            this.hostText.Location = new System.Drawing.Point(109, 23);
            this.hostText.Margin = new System.Windows.Forms.Padding(4, 4, 4, 4);
            this.hostText.Name = "hostText";
            this.hostText.Size = new System.Drawing.Size(177, 22);
            this.hostText.TabIndex = 0;
            this.hostText.Text = "localhost";
            // 
            // line
            // 
            this.line.Location = new System.Drawing.Point(16, 133);
            this.line.Margin = new System.Windows.Forms.Padding(4, 4, 4, 4);
            this.line.Name = "line";
            this.line.Padding = new System.Windows.Forms.Padding(4, 4, 4, 4);
            this.line.Size = new System.Drawing.Size(272, 4);
            this.line.TabIndex = 29;
            this.line.TabStop = false;
            // 
            // cancelButton
            // 
            this.cancelButton.Anchor = System.Windows.Forms.AnchorStyles.None;
            this.cancelButton.DialogResult = System.Windows.Forms.DialogResult.Cancel;
            this.cancelButton.Location = new System.Drawing.Point(76, 146);
            this.cancelButton.Margin = new System.Windows.Forms.Padding(4, 4, 4, 4);
            this.cancelButton.Name = "cancelButton";
            this.cancelButton.Size = new System.Drawing.Size(100, 28);
            this.cancelButton.TabIndex = 3;
            this.cancelButton.Text = "Cancel";
            // 
            // okButton
            // 
            this.okButton.Anchor = System.Windows.Forms.AnchorStyles.None;
            this.okButton.DialogResult = System.Windows.Forms.DialogResult.OK;
            this.okButton.Location = new System.Drawing.Point(188, 146);
            this.okButton.Margin = new System.Windows.Forms.Padding(4, 4, 4, 4);
            this.okButton.Name = "okButton";
            this.okButton.Size = new System.Drawing.Size(100, 28);
            this.okButton.TabIndex = 4;
            this.okButton.Text = "OK";
            this.okButton.Click += new System.EventHandler(this.okButton_Click_1);
            // 
            // portUpDown
            // 
            this.portUpDown.Location = new System.Drawing.Point(109, 55);
            this.portUpDown.Margin = new System.Windows.Forms.Padding(4, 4, 4, 4);
            this.portUpDown.Maximum = new decimal(new int[] {
            65535,
            0,
            0,
            0});
            this.portUpDown.Minimum = new decimal(new int[] {
            1,
            0,
            0,
            0});
            this.portUpDown.Name = "portUpDown";
            this.portUpDown.Size = new System.Drawing.Size(125, 22);
            this.portUpDown.TabIndex = 1;
            this.portUpDown.Value = new decimal(new int[] {
            7777,
            0,
            0,
            0});
            // 
            // portLabel
            // 
            this.portLabel.AutoSize = true;
            this.portLabel.Location = new System.Drawing.Point(16, 58);
            this.portLabel.Margin = new System.Windows.Forms.Padding(4, 0, 4, 0);
            this.portLabel.Name = "portLabel";
            this.portLabel.Size = new System.Drawing.Size(51, 17);
            this.portLabel.TabIndex = 5;
            this.portLabel.Text = "PORT:";
            // 
            // hostLabel
            // 
            this.hostLabel.AutoSize = true;
            this.hostLabel.Location = new System.Drawing.Point(16, 23);
            this.hostLabel.Margin = new System.Windows.Forms.Padding(4, 0, 4, 0);
            this.hostLabel.Name = "hostLabel";
            this.hostLabel.Size = new System.Drawing.Size(51, 17);
            this.hostLabel.TabIndex = 4;
            this.hostLabel.Text = "HOST:";
            // 
            // ConnectForm
            // 
            this.AutoScaleDimensions = new System.Drawing.SizeF(8F, 16F);
            this.AutoScaleMode = System.Windows.Forms.AutoScaleMode.Font;
            this.ClientSize = new System.Drawing.Size(316, 190);
            this.ControlBox = false;
            this.Controls.Add(this.groupBox);
            this.FormBorderStyle = System.Windows.Forms.FormBorderStyle.FixedDialog;
            this.Margin = new System.Windows.Forms.Padding(4, 4, 4, 4);
            this.Name = "ConnectForm";
            this.Text = "Connect";
            this.Load += new System.EventHandler(this.SimulatorForm_Load);
            this.groupBox.ResumeLayout(false);
            this.groupBox.PerformLayout();
            ((System.ComponentModel.ISupportInitialize)(this.portUpDown)).EndInit();
            this.ResumeLayout(false);

        }

        #endregion

        private System.Windows.Forms.GroupBox groupBox;
        private System.Windows.Forms.GroupBox line;
        private System.Windows.Forms.Button cancelButton;
        private System.Windows.Forms.Button okButton;
        private System.Windows.Forms.NumericUpDown portUpDown;
        private System.Windows.Forms.Label portLabel;
        private System.Windows.Forms.Label hostLabel;
        private System.Windows.Forms.TextBox hostText;

    }
}