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

package server

import "testing"

// sends consumed tokens to the channel
type chanTokenConsumer struct {
	channel chan *token
}

func (consumer *chanTokenConsumer) Consume(t *token) {
	consumer.channel <- t
	if t.typ == tokenTypeEOF {
		close(consumer.channel)
	}
}

func validateTokens(t *testing.T, expected []token, tokens chan *token) {
	breakLoop := false
	for _, e := range expected {
		g := <-tokens
		if e.typ != g.typ {
			t.Errorf("expected type " + e.typ.String() + " but got " + g.typ.String() + " value: " + g.val)
			breakLoop = true
		}
		if e.val != g.val {
			t.Errorf("expected value " + e.val + " but got " + g.val)
			breakLoop = true
		}
		if breakLoop {
			break
		}
	}
}
