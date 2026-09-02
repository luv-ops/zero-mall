package mq

type CartChangeMsg struct {
	UserId    string `json:"usrId"`
	TimeStamp int64  `json:"timeStamp"`
}
