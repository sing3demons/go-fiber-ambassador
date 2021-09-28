package main

import (
	"github.com/bxcodec/faker/v3"
	"github.com/sing3demons/main/database"
	"github.com/sing3demons/main/models"
)

func main() {
	database.Connect()
	db := database.GetDB()
	numberOfUser := 30
	// users := make([]models.User, numberOfUser)
	for i := 1; i < numberOfUser; i++ {
		ambassador := models.User{
			FirstName:    faker.FirstName(),
			LastName:     faker.LastName(),
			Email:        faker.Email(),
			IsAmbassador: true,
		}
		ambassador.HashPassword("1234")

	
		db.Create(&ambassador)
	}
}
