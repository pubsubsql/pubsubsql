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
            this.licenseText = new System.Windows.Forms.TextBox();
            this.okButton = new System.Windows.Forms.Button();
            this.copyrightLable = new System.Windows.Forms.Label();
            this.panel1 = new System.Windows.Forms.Panel();
            this.panel1.SuspendLayout();
            this.SuspendLayout();
            // 
            // licenseText
            // 
            this.licenseText.Anchor = ((System.Windows.Forms.AnchorStyles)((((System.Windows.Forms.AnchorStyles.Top | System.Windows.Forms.AnchorStyles.Bottom) 
            | System.Windows.Forms.AnchorStyles.Left) 
            | System.Windows.Forms.AnchorStyles.Right)));
            this.licenseText.BackColor = System.Drawing.SystemColors.Control;
            this.licenseText.BorderStyle = System.Windows.Forms.BorderStyle.FixedSingle;
            this.licenseText.Location = new System.Drawing.Point(6, 3);
            this.licenseText.Multiline = true;
            this.licenseText.Name = "licenseText";
            this.licenseText.Size = new System.Drawing.Size(279, 107);
            this.licenseText.TabIndex = 28;
            this.licenseText.Text = resources.GetString("licenseText.Text");
            // 
            // okButton
            // 
            this.okButton.Anchor = ((System.Windows.Forms.AnchorStyles)((System.Windows.Forms.AnchorStyles.Bottom | System.Windows.Forms.AnchorStyles.Right)));
            this.okButton.DialogResult = System.Windows.Forms.DialogResult.OK;
            this.okButton.Location = new System.Drawing.Point(186, 133);
            this.okButton.Name = "okButton";
            this.okButton.Size = new System.Drawing.Size(99, 23);
            this.okButton.TabIndex = 27;
            this.okButton.Text = "OK";
            // 
            // copyrightLable
            // 
            this.copyrightLable.AutoSize = true;
            this.copyrightLable.Location = new System.Drawing.Point(3, 0);
            this.copyrightLable.Name = "copyrightLable";
            this.copyrightLable.Size = new System.Drawing.Size(0, 13);
            this.copyrightLable.TabIndex = 4;
            // 
            // panel1
            // 
            this.panel1.BorderStyle = System.Windows.Forms.BorderStyle.FixedSingle;
            this.panel1.Controls.Add(this.okButton);
            this.panel1.Controls.Add(this.licenseText);
            this.panel1.Controls.Add(this.copyrightLable);
            this.panel1.Dock = System.Windows.Forms.DockStyle.Fill;
            this.panel1.Location = new System.Drawing.Point(0, 0);
            this.panel1.Name = "panel1";
            this.panel1.Size = new System.Drawing.Size(290, 161);
            this.panel1.TabIndex = 28;
            // 
            // AboutForm
            // 
            this.AutoScaleDimensions = new System.Drawing.SizeF(6F, 13F);
            this.AutoScaleMode = System.Windows.Forms.AutoScaleMode.Font;
            this.ClientSize = new System.Drawing.Size(290, 161);
            this.ControlBox = false;
            this.Controls.Add(this.panel1);
            this.FormBorderStyle = System.Windows.Forms.FormBorderStyle.FixedDialog;
            this.Name = "AboutForm";
            this.Text = "About PubSubSQL Interactive Query";
            this.Load += new System.EventHandler(this.AboutForm_Load);
            this.panel1.ResumeLayout(false);
            this.panel1.PerformLayout();
            this.ResumeLayout(false);

        }

        #endregion

        private System.Windows.Forms.Button okButton;
        private System.Windows.Forms.Label copyrightLable;
        private System.Windows.Forms.TextBox licenseText;
        private System.Windows.Forms.Panel panel1;

    }
}