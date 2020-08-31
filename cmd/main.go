package main

import (
	"apiserver/internal/app/engine"
	"apiserver/internal/app/model"
	"apiserver/internal/app/service"
	"fmt"
	"log"

	"github.com/douyu/jupiter"
)

func main() {
	eng := engine.NewEngine()
	eng.RegisterHooks(jupiter.StageAfterStop, func() error {
		fmt.Println("exit jupiter app ...")
		return nil
	})

	model.Init()
	service.Init()
	if err := eng.Run(); err != nil {
		log.Fatal(err)
	}
}
