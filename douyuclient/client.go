package douyuclient

import (
	"errors"
	"net"
	"sync"
	"time"
)

type Client interface {
	AddRoom(roomId string)
	Start() (<-chan *Message, error)
	Status() byte
	Stop()
}

const (
	ClientInit = iota
	ClientRunning
	ClientStopping
	ClientStopped
)

const (
	TaskTypeCreate = iota
	TaskTypeDeleteAll
	TaskTypeReconnect
)

type connectTask struct {
	taskType   byte
	roomId     string
	clientConn *clientConn
}

type client struct {
	address               string
	status                byte
	roomIds               []string
	connMap               map[string]*clientConn
	bufferSize            uint32
	message               chan *Message
	connectTaskChan       chan *connectTask
	waitGroup             sync.WaitGroup
	connectTaskHandleQuit chan byte
}

func New(address string, bufferSize uint32) Client {
	return &client{
		address:    address,
		status:     ClientInit,
		roomIds:    make([]string, 0),
		connMap:    make(map[string]*clientConn),
		bufferSize: bufferSize,
	}
}

func (client *client) AddRoom(roomId string) {
	client.roomIds = append(client.roomIds, roomId)
	client.emitCreateConn(roomId, 0)
}

func (client *client) Start() (<-chan *Message, error) {
	if client.status == ClientRunning {
		return nil, errors.New("Client is running")
	}
	if client.status == ClientStopping {
		return nil, errors.New("Client is Stopping")
	}
	client.status = ClientRunning
	client.message = make(chan *Message, client.bufferSize)
	client.connectTaskChan = make(chan *connectTask, 10)
	client.connectTaskHandleQuit = make(chan byte)

	client.waitGroup.Add(1)
	go client.connectTaskHandle()

	for _, roomId := range client.roomIds {
		client.emitCreateConn(roomId, 0)
	}
	return client.message, nil
}

func (client *client) Status() byte {
	return client.status
}

func (client *client) Stop() {
	client.status = ClientStopping
	client.connectTaskHandleQuit <- 1
	client.waitGroup.Wait()
	close(client.message)
	client.status = ClientStopped
}

func (client *client) emitCreateConn(roomId string, afterTime time.Duration) {
	if afterTime > 0 {
		logger.Printf("CreateConn after %ds\n", afterTime)
		time.Sleep(afterTime * time.Second)
	}
	if client.status == ClientRunning {
		client.connectTaskChan <- &connectTask{taskType: TaskTypeCreate, roomId: roomId}
	}
}

func (client *client) tryReconnect(conn *clientConn) {
	logger.Println("tryReconnect after 5s")
	time.Sleep(5 * time.Second)
	if client.status == ClientRunning {
		client.connectTaskChan <- &connectTask{taskType: TaskTypeReconnect, clientConn: conn}
	}
}

func (client *client) connectTaskHandle() {
	defer logger.Println("end", startSpan("connectTaskHandle"))
	defer client.waitGroup.Done()
LOOP:
	for {
		select {
		case task := <-client.connectTaskChan:
			if client.status != ClientRunning {
				continue
			}
			switch task.taskType {
			case TaskTypeCreate:
				roomId := task.roomId
				logger.Println("connect roomId", roomId)
				if conn, err := net.Dial("tcp", client.address); err == nil {
					clientConn := newClientConn(conn, client, roomId)
					client.connMap[roomId] = clientConn
					if err = clientConn.init(); err != nil {
						delete(client.connMap, roomId)
						go client.emitCreateConn(roomId, 10)
					}
				} else {
					logger.Println("connect error", roomId, err)
					go client.emitCreateConn(roomId, 10)
				}
			case TaskTypeDeleteAll:
				client.closeAllConnect()
			case TaskTypeReconnect:
				logger.Println("do reconnect", task.clientConn.roomId)
				conn := task.clientConn
				if err := conn.Close(); err != nil {
					logger.Println("close conn error", err)
				}
				for k, v := range client.connMap {
					if v == conn {
						delete(client.connMap, k)
					}
				}
				client.emitCreateConn(conn.roomId, 0)
			}
		case <-client.connectTaskHandleQuit:
			logger.Println("connectTaskHandleQuit")
			client.closeAllConnect()
			break LOOP
		}
	}
}

func (client *client) closeAllConnect() {
	for k, v := range client.connMap {
		if err := v.Close(); err != nil {
			logger.Println("close conn error", err)
		}
		delete(client.connMap, k)
	}
}
