package socket

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
				_ = conn.Socket.Close()
				_ = conn.RelayListener.Close()
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
