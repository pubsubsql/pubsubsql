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

import "strconv"
import "bufio"
import "os"
import "strings"
import "net"
import "time"

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

func (c *cli) readInput() {
	// we do not join here because there is no way to return from blocking readLine
	defer c.stoper.Stop(0)
	l := newLineReader("q")
	for l.readLine() {
		if len(l.line) > 0 {
			c.fromStdin <- l.line
		}
	}
}

func (c *cli) connect() bool {
	conn, err := net.Dial("tcp", config.netAddress())
	if err != nil {
		println(err)
		return false
	}
	c.conn = conn
	return true
}

func (c *cli) outputError(err string) {
	println("\nerror: " + err)
}

func (c *cli) writeMessages() {
	c.stoper.Join()
	defer c.stoper.Stop(0)
	writer := newNetMessageReaderWriter(c.conn, nil)
	var message string
	ok := true
	for ok {
		select {
		case message, ok = <-c.toServer:
			if ok {
				bytes := []byte(message)
				err := writer.writeHeaderAndMessage(bytes)
				if err != nil {
					c.outputError(err.Error())
					ok = false
				}
			}
		case <-c.stoper.GetChan():
			ok = false
		}
	}
}

func (c *cli) readMessages() {
	c.stoper.Join()
	defer c.stoper.Stop(0)
	reader := newNetMessageReaderWriter(c.conn, nil)
	ok := true
	for ok {
		bytes, err := reader.readMessage()
		if err != nil {
			c.outputError(err.Error())
			break
		}
		select {
		case c.fromServer <- string(bytes):
		case <-c.stoper.GetChan():
			ok = false
		}
	}
}

func (c *cli) run() {
	c.initPrefix()
	// connect to server
	if !c.connect() {
		return
	}
	// read user input
	go c.readInput()
	go c.readMessages()
	go c.writeMessages()
	//
	cout := bufio.NewWriter(os.Stdout)
	ok := true
	var serverMessage string
	var userInput string
	for ok {
		cout.WriteString(c.prefix)
		cout.Flush()

		select {
		case userInput, ok = <-c.fromStdin:
			if ok {
				c.toServer <- userInput
			}
		case serverMessage, ok = <-c.fromServer:
			if ok {
				cout.WriteString(serverMessage)
				cout.WriteString("\n")
				cout.Flush()
			}
		case <-c.stoper.GetChan():
			ok = false
		}
	}
	c.conn.Close()
	c.stoper.Wait(time.Millisecond * config.WAIT_MILLISECOND_CLI_SHUTDOWN)
}

func (c *cli) initPrefix() {
	def := defaultConfig()
	c.prefix = "pubsubsql"
	if def.IP != config.IP {
		c.prefix += " " + config.netAddress()
	} else if def.PORT != config.PORT {
		c.prefix += ":" + strconv.Itoa(int(config.PORT))
	}
	c.prefix += ">"
}
