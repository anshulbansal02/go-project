package chat

import (
	"fmt"
	"math/rand"
	"sync"
)

var source = struct {
	m        map[string]int
	mu       sync.Mutex
	minValue int
}{
	m:        make(map[string]int),
	minValue: 10e5,
}

var generateChatId = func(seed string) int {
	source.mu.Lock()
	defer source.mu.Unlock()

	_, e := source.m[seed]
	if !e {
		source.m[seed] = rand.Int() + source.minValue
	}

	source.m[seed] += 1
	return source.m[seed]

}

func getNamespaceKey(conversationId string) string {
	return fmt.Sprintf("ent:chat:%v", conversationId)
}
