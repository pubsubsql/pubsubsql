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
 * You should have idxeived a copy of the GNU Affero General Public License
 * along with PubSubSQL.  If not, see <http://www.gnu.org/licenses/>.
 */

package pubsubsql

// tagItem is a holder for tags and pubsub 
type tagItem struct {
	head   *tag
	pubsub pubSub
}

type tagMap struct {
	tags map[string]*tagItem
}

func (t *tagMap) init() {
	t.tags = make(map[string]*tagItem)
}

func (t *tagMap) getTag(key string) *tag {
	tagitem := t.tags[key]
	if tagitem != nil {
		return tagitem.head
	}
	return nil
}

// getAddTagItem returns tagItem by key. 
// Create new tagItem and adds it to map if does not exist.
func (t *tagMap) getAddTagItem(key string) *tagItem {
	item := t.tags[key]
	if item == nil {
		item = new(tagItem)
		t.tags[key] = item
	}
	return item
}

// addTag adds tag and returns added tag and pubsub
func (t *tagMap) addTag(key string, idx int) (*tag, *pubSub) {
	item := t.getAddTagItem(key)
	if item.head == nil {
		item.head = addTag(nil, idx)
		return item.head, &item.pubsub
	}
	return addTag(item.head, idx), &item.pubsub
}

// addSubscription adds subscription and returns it 
func (t *tagMap) addSubscription(key string, sender *responseSender) *subscription {
	item := t.getAddTagItem(key)
	sub := newSubscription(sender)
	item.pubsub.add(sub)
	return sub
}

// containsTag returns true only if there is a valid head for a given tagItem
func (t *tagMap) containsTag(key string) bool {
	item := t.tags[key]
	if item != nil && item.head != nil {
		return true
	}
	return false
}

// removeTag removes tagItem only if there are no active subscriptions
func (t *tagMap) removeTag(key string) {
	item := t.tags[key]
	if item != nil && item.pubsub.count() == 0 {
		delete(t.tags, key)
	}
}
