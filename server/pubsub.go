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

import "fmt"

// pubsub  
type pubSub struct {
	head *subscription
}

func (p *pubSub) hasSubscriptions() bool {
	return p.head != nil
}

func (p *pubSub) add(s *subscription) {
	s.next = p.head
	p.head = s
}

type pubsubVisitor func(s *subscription) bool

func (p *pubSub) visit(v pubsubVisitor) {
	prev := p.head
	for sub := p.head; sub != nil; sub = sub.next {
		if !sub.active() || !v(sub) {
			if sub == p.head {
				p.head = sub.next
				prev = p.head
			} else {
				prev.next = sub.next
			}
		} else {
			prev = sub
		}
	}
}

func (p *pubSub) count() int {
	i := 0
	f := func(s *subscription) bool {
		i++
		return true
	}
	p.visit(f)
	return i
}

func (p *pubSub) publish(r response) {
	f := func(s *subscription) bool {
		fmt.Println(r)
		return true
	}
	p.visit(f)
}

// subscription represents individual client subscription
type subscription struct {
	next   *subscription // next node
	sender *responseSender
}

// factory
func newSubscription(sender *responseSender) *subscription {
	return &subscription{
		next:   nil,
		sender: sender,
	}
}

//
func (s *subscription) active() bool {
	return s.sender != nil
}

//
func (s *subscription) deactivate() {
	s.sender = nil
}
