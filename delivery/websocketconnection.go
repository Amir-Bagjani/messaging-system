package delivery

import (
	"log"
	"messaging-system/entity"
	"messaging-system/service/messageservice"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type WebSocketHandler struct {
	messageService *messageservice.MessageService
}

func NewWebSocketHandler(messageService *messageservice.MessageService) *WebSocketHandler {
	return &WebSocketHandler{messageService: messageService}
}

func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		http.Error(w, "Missing userId parameter", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	h.messageService.Register(entity.User{ID: userID, Conn: conn})

	go func() {
		defer h.messageService.Unregister(userID, conn)
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				log.Printf("User %s disconnected", userID)
				break
			}
		}
	}()
}

