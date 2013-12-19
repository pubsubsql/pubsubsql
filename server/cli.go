/* Copyright (C) 2013 CompleteDB LLC.
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

package pubsubsql

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type lineReader struct {
	reader *bufio.Reader
	quit   string
	line   string
}

func newLineReader(quit string) *lineReader {
	return &lineReader{
		reader: bufio.NewReader(os.Stdin),
		quit:   quit,
	}
}

func (l *lineReader) readLine() bool {
	line, err := l.reader.ReadString('\n')
	l.line = strings.TrimSpace(line)
	if err != nil {
		return false
	}
	return l.line != l.quit
}

type cli struct {
	prefix     string
	stoper     *Stoper
	fromStdin  chan string
	fromServer chan string
	toServer   chan string
	conn       net.Conn
}

func newCli() *cli {
	return &cli{
		stoper:     NewStoper(),
		fromStdin:  make(chan string),
		fromServer: make(chan string),
		toServer:   make(chan string),
	}
}

func (this *cli) readInput() {
	// we do not join the stoper because there is no way to return from blocking readLine
	defer this.stoper.Stop(0)
	cin := newLineReader("q")
	for cin.readLine() {
		if len(cin.line) > 0 {
			this.fromStdin <- cin.line
		}
	}
}

func (this *cli) connect() bool {
	conn, err := net.Dial("tcp", config.netAddress())
	if err != nil {
		this.outputError(err)
		return false
	}
	this.conn = conn
	return true
}

func (this *cli) outputError(err error) {
	fmt.Println("error: ", err)
}

func (this *cli) writeMessages() {
	this.stoper.Join()
	defer this.stoper.Stop(0)
	writer := newNetMessageReaderWriter(this.conn, nil)
	var message string
	ok := true
	for ok {
		select {
		case message, ok = <-this.toServer:
			if ok {
				bytes := []byte(message)
				err := writer.writeHeaderAndMessage(bytes)
				if err != nil {
					this.outputError(err)
					ok = false
				}
			}
		case <-this.stoper.GetChan():
			ok = false
		}
	}
}

func (this *cli) readMessages() {
	this.stoper.Join()
	defer this.stoper.Stop(0)
	reader := newNetMessageReaderWriter(this.conn, nil)
	ok := true
	for ok {
		bytes, err := reader.readMessage()
		if err != nil {
			this.outputError(err)
			break
		}
		select {
		case this.fromServer <- string(bytes):
		case <-this.stoper.GetChan():
			ok = false
		}
	}
}

func (this *cli) run() {
	this.initPrefix()
	// connect to server
	if !this.connect() {
		return
	}
	// read user input
	go this.readInput()
	go this.readMessages()
	go this.writeMessages()
	//
	cout := bufio.NewWriter(os.Stdout)
	ok := true
	var serverMessage string
	var userInput string
	for ok {
		cout.WriteString(this.prefix)
		cout.Flush()

		select {
		case userInput, ok = <-this.fromStdin:
			if ok {
				this.toServer <- userInput
			}
		case serverMessage, ok = <-this.fromServer:
			if ok {
				cout.WriteString(serverMessage)
				cout.WriteString("\n")
				cout.Flush()
			}
		case <-this.stoper.GetChan():
			ok = false
		}
	}
	this.conn.Close()
	this.stoper.Wait(time.Millisecond * config.WAIT_MILLISECOND_CLI_SHUTDOWN)
}

func (this *cli) initPrefix() {
	def := defaultConfig()
	this.prefix = "pubsubsql"
	if def.IP != config.IP {
		this.prefix += " " + config.netAddress()
	} else if def.PORT != config.PORT {
		this.prefix += ":" + strconv.Itoa(int(config.PORT))
	}
	this.prefix += ">"
}
