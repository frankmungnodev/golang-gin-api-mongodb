package main

import (
	"franky/go-api-gin/envs"
	"franky/go-api-gin/mongodb"
	"franky/go-api-gin/routes"
	"log"

	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
)

func main() {

	envs.LoadEnvs()
	mongodb.InitMongoDB()

	if envs.ENV == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	server := gin.Default()

	// Routers
	authrouter := &routes.AuthRouter{}

	// Mount routers
	apiV1Router := server.Group("/api/v1")
	authrouter.DefineV1Routes(apiV1Router)

	server.Run(":8080")
	log.Fatal(autotls.Run(server, "localhost"))
}
