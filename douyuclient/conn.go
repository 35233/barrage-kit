package douyuclient

import (
	"encoding/binary"
	"github.com/35233/barrage-kit/stt"
	"net"
	"time"
)

type ConnectStatus byte

const (
	ConnectConnected ConnectStatus = iota
	ConnectDisconnected
)

type clientConn struct {
	conn    net.Conn
	client  *client
	status  ConnectStatus
	roomId  string
	revBuff []byte
}

func newClientConn(conn net.Conn, client *client, roomId string) *clientConn {
	return &clientConn{conn, client, ConnectConnected, roomId, nil}
}

func (clientConn *clientConn) init() error {
	clientConn.client.waitGroup.Add(3)
	reverseChan := make(chan *Message)
	go clientConn.reverseConnMessage(reverseChan)
	loginStr := stt.Encode(map[string]interface{}{
		"type": "loginreq",
	})

	if _, err := clientConn.Send(CreateRawSendMsg(loginStr)); err != nil {
		logger.Println("send loginreq error", err)
		return err
	}
	loginRes := <-reverseChan
	logger.Println("loginResContent", loginRes.Text())
	joinGroupStr := stt.Encode(map[string]interface{}{
		"type": "joingroup",
		"rid":  clientConn.roomId,
		"gid":  "-9999",
	})
	if _, err := clientConn.Send(CreateRawSendMsg(joinGroupStr)); err != nil {
		logger.Println("send joingroup error", err)
		return err
	}
	go clientConn.startHeartBeat()
	go clientConn.forwardChan(reverseChan)
	return nil
}

func (clientConn *clientConn) tryReconnect() {
	clientConn.client.tryReconnect(clientConn)
}

func (clientConn *clientConn) forwardChan(reverseChan <-chan *Message) {
	defer clientConn.client.waitGroup.Done()
	defer logger.Println("end", startSpan("forwardChan"))
	for message := range reverseChan {
		clientConn.client.message <- message
	}
}

func (clientConn *clientConn) isRunning() bool {
	return clientConn.client.status == ClientRunning && clientConn.status == ConnectConnected
}

func (clientConn *clientConn) reverseConnMessage(ch chan *Message) {
	defer clientConn.client.waitGroup.Done()
	defer logger.Println("end", startSpan("reverseConnMessage"))
	clientConn.revBuff = make([]byte, 0)
	for clientConn.isRunning() {
		tmpBuff := make([]byte, 1024)
		if n, err := clientConn.conn.Read(tmpBuff); err == nil {
			var currBuff []byte
			if len(clientConn.revBuff) > 0 {
				revLen := len(clientConn.revBuff)
				currBuff = make([]byte, revLen+n)
				copy(currBuff, clientConn.revBuff)
				copy(currBuff[revLen:], tmpBuff[:n])
			} else {
				currBuff = tmpBuff[:n]
			}
			for len(currBuff) >= 4 {
				curPackLen := binary.LittleEndian.Uint32(currBuff)
				if curPackLen <= uint32(len(currBuff))-4 {
					message := currBuff[:curPackLen+4]
					ch <- &Message{time.Now().UnixNano(), message}
					currBuff = currBuff[curPackLen+4:]
				} else {
					break
				}
			}
			clientConn.revBuff = currBuff
		} else {
			logger.Println("reverseConnMessage error", err)
			break
		}
	}
	close(ch)
	clientConn.tryReconnect()
}

func (clientConn *clientConn) RoomId() string {
	return clientConn.roomId
}

func (clientConn *clientConn) Send(buff []byte) (int, error) {
	n, err := clientConn.conn.Write(buff)
	if err != nil {
		logger.Println("Send error", err)
	}
	return n, err
}

func (clientConn *clientConn) Close() error {
	clientConn.status = ConnectDisconnected
	return clientConn.conn.Close()
}

func (clientConn *clientConn) startHeartBeat() {
	defer clientConn.client.waitGroup.Done()
	defer logger.Println("end", startSpan("startHeartBeat"))
	mrklStr := stt.Encode(map[string]interface{}{
		"type": "mrkl",
	})
	for clientConn.isRunning() {
		if _, err := clientConn.Send(CreateRawSendMsg(mrklStr)); err != nil {
			logger.Println("send mrkl error", err)
			if !clientConn.isRunning() {
				break
			}
		}
		time.Sleep(30 * time.Second)
	}
}
