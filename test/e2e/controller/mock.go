package controller

type MessageSendEvent struct {
	UserId     int64
	ToUserId   int64
	MsgContent string
}

type MessagePushEvent struct{}
