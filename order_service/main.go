package main

import (
	"fmt"
	"log"

	"github.com/palashbhasme/order_service/utils"
)

func main() {
	fmt.Println("Jello")

	logger, err := utils.InitLogger()
	if err != nil {
		log.Fatal("Failed to initialze logger")
	}
	defer logger.Sync()

}
