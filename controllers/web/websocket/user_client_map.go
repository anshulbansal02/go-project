package web

type clientMap struct {
	// Stores UserId -> ClientId mapping
	userClient map[string]string
	// Stores ClientId -> UserId mapping
	clientUser map[string]string
}

func NewClientMap() *clientMap {
	return &clientMap{
		userClient: make(map[string]string),
		clientUser: make(map[string]string),
	}
}

func (cm *clientMap) GetUserId(clientId string) (userId string, exists bool) {
	userId, exists = cm.clientUser[clientId]
	return
}

func (cm *clientMap) GetClientId(userId string) (clientId string, exists bool) {
	clientId, exists = cm.userClient[clientId]
	return
}

func (cm *clientMap) Add(clientId string, userId string) {
	cm.clientUser[clientId] = userId
	cm.userClient[userId] = clientId
}

func (cm *clientMap) RemoveUser(userId string) {
	clientId := cm.userClient[userId]

	delete(cm.clientUser, clientId)
	delete(cm.userClient, userId)
}

func (cm *clientMap) RemoveClient(clientId string) {
	userId := cm.clientUser[clientId]

	delete(cm.clientUser, clientId)
	delete(cm.userClient, userId)
}
