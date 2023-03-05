package main

import (
	"fmt"
	"task5-btpns-stella/router"

	swagfile "github.com/swaggo/files"
	ginswag "github.com/swaggo/gin-swagger"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println(os.Getenv("GRPC_PORT"))

	ginApi, port := router.Route()
	if ginApi == nil {
		log.Fatal("Error GinRoute Initialized")
	}

	url := ginswag.URL(fmt.Sprintf("http://localhost:%s/swagger/doc.json", port))
	ginApi.GET("/swagger/*any", ginswag.WrapHandler(swagfile.Handler, url))
	_ = ginApi.Run(":" + port)

}
