package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)



func main() {
	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Cat API",
		AppName:       "CAT API JEJE",
	})

	app.Use(logger.New())

	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	app.Get("/stack", func(c *fiber.Ctx) error {
		return c.JSON(c.App().Stack())
	})

	app.Get("/cat", getCat)
	app.Get("/cat/:id", getCatPreview)
	app.Listen(":3007")
}

func getCat(c *fiber.Ctx) error {
	// Make the GET request to the cat API
	resp, err := http.Get("https://api.thecatapi.com/v1/images/search")
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// Create a variable to hold the JSON data
	var jsonMap []map[string]interface{} // Or use a specific struct for better type safety

	// Parse the JSON response into jsonMap
	if err := json.Unmarshal(body, &jsonMap); err != nil {
		log.Fatalln(err)
	}

	fmt.Println(jsonMap)

	// Set the response content type and return the parsed JSON
	c.Accepts("application/json")
	return c.JSON(jsonMap)
}

func getCatPreview(c *fiber.Ctx) error {
	// Make the GET request to the cat API
	id := c.Params("id")

	url := "https://cdn2.thecatapi.com/images/" + id + ".jpg"
	resp, err := http.Get(url)

	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	reader := io.NopCloser(strings.NewReader(string(body)))

	// Set the response content type and return the image stream
	c.Set("Content-Type", "image/jpeg")
	return c.SendStream(reader)
}

