package exchange

import (
	"sync"
)

type ClientMap struct {
	mu sync.RWMutex
	// Stores UserId -> ClientId mapping
	userClient map[string]string
	// Stores ClientId -> UserId mapping
	clientUser map[string]string
}

func NewClientMap() *ClientMap {
	return &ClientMap{
		userClient: make(map[string]string),
		clientUser: make(map[string]string),
	}
}

func (cm *ClientMap) GetClientIds(userIds []string) []string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	clientIds := make([]string, len(userIds))

	for i, userId := range userIds {
		clientIds[i] = cm.userClient[userId]
	}

	return clientIds
}

func (cm *ClientMap) GetUserId(clientId string) (userId string, exists bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	userId, exists = cm.clientUser[clientId]
	return
}

func (cm *ClientMap) GetClientId(userId string) (clientId string, exists bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	clientId, exists = cm.userClient[userId]
	return
}

func (cm *ClientMap) Add(clientId string, userId string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.clientUser[clientId] = userId
	cm.userClient[userId] = clientId

}

func (cm *ClientMap) RemoveUser(userId string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	clientId := cm.userClient[userId]

	delete(cm.clientUser, clientId)
	delete(cm.userClient, userId)
}

func (cm *ClientMap) RemoveClient(clientId string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	userId := cm.clientUser[clientId]

	delete(cm.clientUser, clientId)
	delete(cm.userClient, userId)
}
