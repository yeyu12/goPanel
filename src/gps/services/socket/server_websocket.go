package socket

type ServerWebsocketManager struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	UnRegister chan *Client
}

var ServerWsManager = ServerWebsocketManager{
	Broadcast:  make(chan []byte),
	Register:   make(chan *Client),
	UnRegister: make(chan *Client),
	Clients:    make(map[*Client]bool),
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
				close(conn.WsRead)
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
