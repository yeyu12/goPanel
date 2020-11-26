package socket

import log "github.com/sirupsen/logrus"

type ServerWebsocketManager struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	UnRegister chan *Client
}

func (manager *ServerWebsocketManager) Start() {
	defer func() {
		recover()
	}()

	for {
		select {
		case conn := <-manager.Register:
			manager.Clients[conn] = true
		case conn := <-manager.UnRegister:
			if _, ok := manager.Clients[conn]; ok {
				if err := conn.Socket.Close(); err != nil {
					log.Error(err)
				}
				if conn.RelayConn != nil {
					if err := conn.RelayConn.Close(); err != nil {
						log.Error(err)
					}
				}
				if err := conn.RelayListener.Close(); err != nil {
					log.Error(err)
				}

				ControlManager.PushRecoveryPort(conn.RelayPort)

				close(conn.wsRead)
				close(conn.Send)
				delete(manager.Clients, conn)
			}
		case message := <-manager.Broadcast:
			for conn := range manager.Clients {
				select {
				case conn.Send <- message:
				default:
					close(conn.Send)
					delete(manager.Clients, conn)
				}
			}
		}
	}
}

func (manager *ServerWebsocketManager) SendAll(message []byte, ignore *Client) {
	for conn := range manager.Clients {
		if conn != ignore {
			conn.Send <- message
		}
	}
}
