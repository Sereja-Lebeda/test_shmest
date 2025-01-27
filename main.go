package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Эндпоинт для вебхука
	app.Post("/webhook", func(c *fiber.Ctx) error {
		// Получаем тело запроса
		body := c.Body()

		// Добавляем произвольные сообщения в консоль
		fmt.Println("\nNew hook triggered!")

		// Попробуем обработать тело как JSON
		var formattedJSON map[string]interface{}
		if err := json.Unmarshal(body, &formattedJSON); err != nil {
			// Если тело не JSON, пробуем декодировать как URL-encoded данные
			parsedData, parseErr := url.ParseQuery(string(body))
			if parseErr != nil {
				// Если данные не распознаны, выводим их как есть
				log.Printf("Webhook received (raw): %s", body)
			} else {
				// Преобразуем URL-encoded данные в map
				decodedData := make(map[string]interface{})
				for key, values := range parsedData {
					if len(values) > 1 {
						decodedData[key] = values
					} else {
						decodedData[key] = values[0]
					}
				}

				// Форматируем и выводим как JSON
				prettyBody, _ := json.MarshalIndent(decodedData, "", "  ")
				log.Println("Webhook received (decoded URL-encoded):")
				log.Println(string(prettyBody))
			}
		} else {
			// Если тело валидное JSON, форматируем его с отступами
			prettyBody, _ := json.MarshalIndent(formattedJSON, "", "  ")
			log.Printf("Webhook received (formatted JSON):\n%s", prettyBody)
		}

		return c.SendString("Webhook received successfully!")
	})

	// Запускаем сервер
	log.Println("Server is running on http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}
