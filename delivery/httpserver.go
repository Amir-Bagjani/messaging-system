package delivery

import (
	"fmt"
	"messaging-system/pkg/rabbitmq"
	"messaging-system/service/messageservice"
	"net/http"
)

type HTTPHandler struct {
	messageService *messageservice.MessageService
	rabbitMQ       *rabbitmq.RabbitMQ
}

func NewHTTPHandler(messageService *messageservice.MessageService, rabbitMQ *rabbitmq.RabbitMQ) *HTTPHandler {
	return &HTTPHandler{messageService: messageService, rabbitMQ: rabbitMQ}
}

func (h *HTTPHandler) HandleSendMessage(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	message := r.URL.Query().Get("message")

	if userID == "" || message == "" {
		http.Error(w, "Missing userId or message parameter", http.StatusBadRequest)
		return
	}

	if err := h.rabbitMQ.Publish("user_messages", message); err != nil {
		http.Error(w, "Failed to queue message", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Message queued for user %s", userID)
}