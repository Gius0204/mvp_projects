package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/segmentio/kafka-go"
)

type Notification struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

var notifications []Notification
var mu sync.Mutex

func consumeKafkaMessages() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "notifications",
		GroupID: "notification-service",
	})

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println("Error al leer mensaje de Kafka:", err)
			continue
		}

		var notif Notification
		if err := json.Unmarshal(m.Value, &notif); err != nil {
			log.Println("Error al parsear mensaje:", err)
			continue
		}

		// Guardar en memoria
		mu.Lock()
		notifications = append(notifications, notif)
		mu.Unlock()

		fmt.Println("Nueva notificación:", notif)
	}
}

func main() {
	app := fiber.New()

	// Ruta para obtener notificaciones
	app.Get("/notifications", func(c *fiber.Ctx) error {
		mu.Lock()
		defer mu.Unlock()
		return c.JSON(notifications)
	})

	// Iniciar consumidor Kafka en segundo plano
	go consumeKafkaMessages()

	// Iniciar servidor
	log.Fatal(app.Listen(":3004"))
}
