package main

import (
	"log"
	"ocenakademik/db"
	"ocenakademik/internal/user"
	"ocenakademik/router"
	"ocenakademik/util"
)

func main() {
	c, err := util.LoadConfig()
	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	hasher := util.NewSha1Hasher(c.JwtSecret)

	db, err := db.NewDatabase(c.DBUrl)
	if err != nil {
		log.Fatalln("Failed start a databse", err)
	}

	userRep := user.NewUserRepository(db.GetDB())
	userSvc := user.NewService(userRep, hasher, []byte(c.JwtSecret), c.TokenTtl)
	userHandler := user.NewHandler(userSvc)

	router.InitRouter(userHandler, c.JwtSecret)
	router.Start("0.0.0.0:8080")
}
