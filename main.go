package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	httpServer := gin.Default()
	Run(httpServer)
}

//hello
