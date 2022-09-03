package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

var cache = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

var ctx = context.Background()

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("It is working 👊")
	})

	app.Get("/:id", verifyCache, func(c *fiber.Ctx) error {
		id := c.Params("id")
		res, err := http.Get("https://jsonplaceholder.typicode.com/users/" + id)
		if err != nil {
			return err
		}

		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		cacheErr := cache.Set(ctx, id, body, 10*time.Second).Err()
		if cacheErr != nil {
			return cacheErr
		}

		data := toJson(body)
		return c.JSON(fiber.Map{"Data": data})
	})

	app.Listen(":3000")
}

func verifyCache(c *fiber.Ctx) error {
	id := c.Params("id")
	val, err := cache.Get(ctx, id).Bytes()
	if err != nil {
		return c.Next()
	}

	data := toJson(val)
	return c.JSON(fiber.Map{"Cached": data})
}

func toJson(val []byte) User {
	user := User{}
	err := json.Unmarshal(val, &user)
	if err != nil {
		panic(err)
	}

	return user
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Address  `json:"address"`
	Phone    string `json:"phone"`
	Website  string `json:"website"`
	Company  `json:"company"`
}

type Address struct {
	Street  string `json:"street"`
	Suite   string `json:"suite"`
	City    string `json:"city"`
	Zipcode string `json:"zipcode"`
	Geo     `json:"geo"`
}

type Geo struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

type Company struct {
	Name        string `json:"name"`
	CatchPhrase string `json:"catchPhrase"`
	Bs          string `json:"bs"`
}
