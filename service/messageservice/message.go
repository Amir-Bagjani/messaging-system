package messageservice

import (
	"encoding/json"
	"log"
	"messaging-system/entity"
	"messaging-system/pkg/rabbitmq"
	"messaging-system/repository/dbrepository"
	"sync"
)

type MessageService struct {
	clients map[string][]entity.User 
	mu          sync.Mutex
	dbRepo *dbrepository.MessageRepository
}

func NewMessageService(dbRepo *dbrepository.MessageRepository) *MessageService {
	return &MessageService{
		clients: make(map[string][]entity.User),
		dbRepo:  dbRepo,
	}
}

func (s *MessageService) Register(user entity.User) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Append the new connection to the user's connection list
	s.clients[user.ID] = append(s.clients[user.ID], user)
	log.Printf("User %s connected on a new device", user.ID)

	// Send message history to the new connection only
	messages, err := s.GetMessages(user.ID)
	if err != nil {
		log.Printf("Failed to retrieve messages for user %s: %v", user.ID, err)
		return
	}
	if len(messages) > 0 {
		messageHistory, err := json.Marshal(messages)
		if err != nil {
			log.Printf("Failed to marshal message history for user %s: %v", user.ID, err)
			return
		}
		if err := user.Conn.WriteMessage(1, messageHistory); err != nil {
			log.Printf("Error sending message history to user %s: %v", user.ID, err)
		}
	}
}

func (s *MessageService) GetMessages(userID string) ([]entity.Message, error) {
	return s.dbRepo.GetMessagesByUserID(userID)
}

func (s *MessageService) Unregister(userID string, conn entity.WebSocketConn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	connections, ok := s.clients[userID]
	if !ok {
		return
	}

	remainingConnections := []entity.User{}
	for _, userConn := range connections {
		if userConn.Conn != conn {
			remainingConnections = append(remainingConnections, userConn)
		} else {
			userConn.Conn.Close()
			log.Printf("User %s disconnected on one device", userID)
		}
	}

	if len(remainingConnections) == 0 {
		delete(s.clients, userID)
	} else {
		s.clients[userID] = remainingConnections
	}
}


func (s *MessageService) SendMessage(msg entity.Message) error {
	if err := s.dbRepo.StoreMessage(msg.UserID, msg.Message); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	updatedMessages, err := s.GetMessages(msg.UserID)
	if err != nil {
		log.Printf("Failed to retrieve updated messages for user %s: %v", msg.UserID, err)
		return err
	}

	messageHistory, err := json.Marshal(updatedMessages)
	if err != nil {
		log.Printf("Failed to marshal updated message history for user %s: %v", msg.UserID, err)
		return err
	}

	connections, ok := s.clients[msg.UserID]
	if !ok {
		log.Printf("User %s not connected", msg.UserID)
		return nil
	}

	for _, userConn := range connections {
		if err := userConn.Conn.WriteMessage(1, messageHistory); err != nil {
			log.Printf("Error sending updated message list to user %s on one device: %v", msg.UserID, err)
			continue
		}
	}

	return nil
}


func (s *MessageService) StartRabbitMQConsumer(rabbitMQ *rabbitmq.RabbitMQ, queueName string) {
	msgs, err := rabbitMQ.Consume(queueName)
	if err != nil {
		log.Fatalf("Failed to start RabbitMQ consumer: %v", err)
	}

	go func() {
		for d := range msgs {
			var msg entity.Message
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				log.Printf("Error decoding message: %v", err)
				continue
			}

			s.SendMessage(msg)
		}
	}()
}
