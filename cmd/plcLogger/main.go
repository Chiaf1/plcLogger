package main

import (
	"fmt"

	"github.com/chiaf1/plclogger/internal/config"
)

func main() {
	var conf config.Config
	err := conf.Load("./config.yaml")
	if err != nil {
		panic(err)
	}
	fmt.Println(conf)
	per := conf.GetPeriodicLogTags()
	fmt.Println(per)
}
