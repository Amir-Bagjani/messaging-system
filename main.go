package main

import (
	"log"
	"messaging-system/delivery"
	"messaging-system/pkg/environmentvariable"
	"messaging-system/pkg/rabbitmq"
	"messaging-system/repository/dbrepository"
	"messaging-system/service/messageservice"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	environmentvariable.NewEnv()

	databaseURL := environmentvariable.GetEnv("DATABASE_URL")
	rabbitMQURL := environmentvariable.GetEnv("RABBITMQ_URL")

	db := dbrepository.New(databaseURL)
	messageDbRepo := dbrepository.NewMessageRepository(db)

	rabbitMQ, err := rabbitmq.NewRabbitMQ(rabbitMQURL, "user_messages")
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()

	messageService := messageservice.NewMessageService(messageDbRepo)
	wsHandler := delivery.NewWebSocketHandler(messageService)
	httpHandler := delivery.NewHTTPHandler(messageService, rabbitMQ)

	messageService.StartRabbitMQConsumer(rabbitMQ, "user_messages")

	http.HandleFunc("/ws", wsHandler.HandleWebSocket)
	http.HandleFunc("/send", httpHandler.HandleSendMessage)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
