package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	Client *mongo.Client
	DB     *mongo.Database
}

var mg MongoInstance

const dbName = "fiber-hrm"
const mongoURI = "mongodb://localhost:27017" + dbName

type Employee struct {
	ID     int     `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string  `json:"name"`
	Salary float64 `json:"salary"`
	Age    int     `json:"age"`
}

func Connect() error {

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = client.connect(ctx)
	db := client.Database(dbName)
	if err != nil {
		return err
	}
	mg = MongoInstance{
		Client: client,
		DB:     db,
	}
	return nil
}

func main() {

	if err := Connect(); err != nil {
		log.Fatal(err)

	}

	app := fiber.New()

	app.Get("/employee", func(c *fiber.Ctx) error {

		query := bson.D{}

		cur, err := mg.DB.Collection("employees").Find(c.Context(), query)
		if err != nil {
			return err.Status(500).SendString(err.Error())
		}
		var employees []Employee = make([]Employee, 0)

		if err := cur.All(c.Context(), &employees); err != nil {
			return err.Status(500).SendString(err.Error())
		}
		return c.JSON(employees)

	})

	app.Post("/employee")
	app.Put("/employee/:id")
	app.Delete("/employee:id")

}
