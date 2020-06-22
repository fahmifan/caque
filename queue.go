package caque

import (
	"sort"
	"sync"
)

// Queue :nodoc:
type Queue struct {
	keys          []string
	mapKeyToIndex map[string]int
	mu            sync.RWMutex
}

// Append :nodoc:
func (q *Queue) Append(key string) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if _, ok := q.mapKeyToIndex[key]; ok {
		return
	}

	q.keys = append(q.keys, key)
	sort.Strings(q.keys)
	q.mapKeyToIndex[key] = sort.SearchStrings(q.keys, key)
}

// Pop :nodoc:
func (q *Queue) Pop() string {
	q.mu.Lock()
	key := q.keys[0]
	q.keys = q.keys[1:]
	delete(q.mapKeyToIndex, key)
	q.mu.Unlock()
	return key
}

// DeleteKey make sure the keys are sorted before
func (q *Queue) DeleteKey(key string) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	if _, ok := q.mapKeyToIndex[key]; !ok {
		return false
	}

	idx := sort.SearchStrings(q.keys, key)
	if idx == len(q.keys) {
		return false
	}

	if idx == 0 {
		q.keys = q.keys[1:]
		return true
	}

	copy(q.keys[idx:], q.keys[idx+1:])
	q.keys = q.keys[:len(q.keys)-1]
	return true
}

// Size :nodoc:
func (q *Queue) Size() int {
	return len(q.keys)
}
