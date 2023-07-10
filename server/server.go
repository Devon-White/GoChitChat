package server

import (
	"context"
	"fmt"
	"log"
	"os"
)

func Init() {

	router := SetupRoutes()

	go func() {
		fmt.Println("Starting server...")
		if err := router.Run(":8080"); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	if err := Run(context.Background(), router); err != nil {
		log.Fatal(err)
	}
}
