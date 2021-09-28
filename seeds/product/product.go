package main

import (
	"math/rand"

	"github.com/bxcodec/faker/v3"
	"github.com/sing3demons/main/database"
	"github.com/sing3demons/main/models"
)

func main() {
	database.Connect()
	db := database.GetDB()
	numberOfProduct := 100000
	products := make([]models.Product, numberOfProduct)
	for i := 1; i < numberOfProduct; i++ {
		product := models.Product{
			Title:       faker.Name(),
			Description: faker.Username(),
			Image:       faker.URL(),
			Price:       float64(rand.Intn(90) + 10),
		}

		products[i] = product
	}

	db.CreateInBatches(products, 1000)
}
