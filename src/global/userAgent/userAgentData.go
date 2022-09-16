package userAgent

type UserAgentData struct {
	AgentAppId string `json:"agentAppId"`
	SocketId   string `json:"socketId"`
	UserId     int64  `json:"userId"`
	LoginAt    int64  `json:"loginAt"`
}
