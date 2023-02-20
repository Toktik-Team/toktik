// Progress:
//   WIP

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"testing"
	"time"
	"toktik/test/e2e/controller"
)

func TestMessageServer(t *testing.T) {
	e := newExpect(t)
	userIdA, _ := getTestUserToken(testUserA, e)
	userIdB, _ := getTestUserToken(testUserB, e)

	connA, err := net.Dial("tcp", "127.0.0.1:9090")
	if err != nil {
		fmt.Printf("Connect server failed: %v\n", err)
		return
	}
	connB, err := net.Dial("tcp", "127.0.0.1:9090")
	if err != nil {
		fmt.Printf("Connect server failed: %v\n", err)
		return
	}

	createChat(userIdA, connA, userIdB, connB)

	go readMessage(connB)
	sendMessage(userIdA, userIdB, connA)
}

func readMessage(conn net.Conn) {
	defer conn.Close()

	var buf [256]byte
	for {
		n, err := conn.Read(buf[:])
		if n == 0 {
			if err == io.EOF {
				break
			}
			fmt.Printf("Read message failed: %v\n", err)
			continue
		}

		var event = controller.MessagePushEvent{}
		_ = json.Unmarshal(buf[:n], &event)
		fmt.Printf("Read message：%+v\n", event)
	}
}

func sendMessage(fromUserId int, toUserId int, fromConn net.Conn) {
	defer fromConn.Close()

	for i := 0; i < 3; i++ {
		time.Sleep(time.Second)
		sendEvent := controller.MessageSendEvent{
			UserId:     int64(fromUserId),
			ToUserId:   int64(toUserId),
			MsgContent: "Test Content",
		}
		data, _ := json.Marshal(sendEvent)
		_, err := fromConn.Write(data)
		if err != nil {
			fmt.Printf("Send message failed: %v\n", err)
			return
		}
	}
	time.Sleep(time.Second)
}

func createChat(userIdA int, connA net.Conn, userIdB int, connB net.Conn) {
	chatEventA := controller.MessageSendEvent{
		UserId:   int64(userIdA),
		ToUserId: int64(userIdB),
	}
	chatEventB := controller.MessageSendEvent{
		UserId:   int64(userIdB),
		ToUserId: int64(userIdA),
	}
	eventA, _ := json.Marshal(chatEventA)
	eventB, _ := json.Marshal(chatEventB)
	_, err := connA.Write(eventA)
	if err != nil {
		fmt.Printf("Create chatA failed: %v\n", err)
		return
	}
	_, err = connB.Write(eventB)
	if err != nil {
		fmt.Printf("Create chatB failed: %v\n", err)
		return
	}
}
