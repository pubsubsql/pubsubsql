namespace PubSubSQLGUI
{
    partial class AboutForm
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
            System.ComponentModel.ComponentResourceManager resources = new System.ComponentModel.ComponentResourceManager(typeof(AboutForm));
            this.groupBox = new System.Windows.Forms.GroupBox();
            this.license = new System.Windows.Forms.TextBox();
            this.okButton = new System.Windows.Forms.Button();
            this.copyrightLable = new System.Windows.Forms.Label();
            this.groupBox.SuspendLayout();
            this.SuspendLayout();
            // 
            // groupBox
            // 
            this.groupBox.Controls.Add(this.license);
            this.groupBox.Controls.Add(this.okButton);
            this.groupBox.Controls.Add(this.copyrightLable);
            this.groupBox.Location = new System.Drawing.Point(5, -1);
            this.groupBox.Name = "groupBox";
            this.groupBox.Size = new System.Drawing.Size(305, 182);
            this.groupBox.TabIndex = 27;
            this.groupBox.TabStop = false;
            // 
            // license
            // 
            this.license.BackColor = System.Drawing.SystemColors.Control;
            this.license.BorderStyle = System.Windows.Forms.BorderStyle.FixedSingle;
            this.license.Location = new System.Drawing.Point(14, 42);
            this.license.Multiline = true;
            this.license.Name = "license";
            this.license.Size = new System.Drawing.Size(280, 90);
            this.license.TabIndex = 28;
            this.license.Text = resources.GetString("license.Text");
            // 
            // okButton
            // 
            this.okButton.Anchor = System.Windows.Forms.AnchorStyles.None;
            this.okButton.DialogResult = System.Windows.Forms.DialogResult.OK;
            this.okButton.Location = new System.Drawing.Point(195, 150);
            this.okButton.Name = "okButton";
            this.okButton.Size = new System.Drawing.Size(99, 23);
            this.okButton.TabIndex = 27;
            this.okButton.Text = "OK";
            // 
            // copyrightLable
            // 
            this.copyrightLable.AutoSize = true;
            this.copyrightLable.Location = new System.Drawing.Point(12, 19);
            this.copyrightLable.Name = "copyrightLable";
            this.copyrightLable.Size = new System.Drawing.Size(181, 13);
            this.copyrightLable.TabIndex = 4;
            this.copyrightLable.Text = "Copyright (C) 2013 CompleteDB LLC.";
            // 
            // AboutForm
            // 
            this.AutoScaleDimensions = new System.Drawing.SizeF(6F, 13F);
            this.AutoScaleMode = System.Windows.Forms.AutoScaleMode.Font;
            this.ClientSize = new System.Drawing.Size(317, 190);
            this.ControlBox = false;
            this.Controls.Add(this.groupBox);
            this.FormBorderStyle = System.Windows.Forms.FormBorderStyle.FixedDialog;
            this.Name = "AboutForm";
            this.Text = "About PubSubSQL Interactive Query";
            this.Load += new System.EventHandler(this.AboutForm_Load);
            this.groupBox.ResumeLayout(false);
            this.groupBox.PerformLayout();
            this.ResumeLayout(false);

        }

        #endregion

        private System.Windows.Forms.GroupBox groupBox;
        private System.Windows.Forms.Button okButton;
        private System.Windows.Forms.Label copyrightLable;
        private System.Windows.Forms.TextBox license;

    }
}