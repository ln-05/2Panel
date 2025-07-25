package main

import (
	router2 "api/router"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	router2.Router(r)
	r.Run(":8889")
}
