package main

import (
	"apiserver/internal/app/engine"
	"apiserver/internal/app/model"
	"apiserver/internal/app/service"
	"fmt"
	compound_registry "github.com/douyu/jupiter/pkg/registry/compound"
	etcdv3_registry "github.com/douyu/jupiter/pkg/registry/etcdv3"
	"log"

	"github.com/douyu/jupiter"
)

func main() {
	eng := engine.NewEngine()
	eng.SetRegistry(
		compound_registry.New(
			etcdv3_registry.StdConfig("bj01").Build(),
		),
	)
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
