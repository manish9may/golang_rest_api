package main

import (
	"golang_rest_api/db"
	"golang_rest_api/routes"
	"sync"

	"github.com/gin-gonic/gin"
)

func main() {
	var wq sync.WaitGroup
	wq.Add(1)
	go func() {
		defer wq.Done()
		db.InitDB()
	}()
	wq.Wait()
	server := gin.Default()
	routes.RegisterRoutes(server)
	server.Run(":5173") // localhost:5173
}
