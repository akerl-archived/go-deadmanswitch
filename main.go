package main

import (
	"fmt"
	"log"

	"github.com/akerl/go-lambda/mux"
	"github.com/akerl/go-lambda/s3"
	"github.com/aws/aws-lambda-go/lambda"
)

type check struct {
	Name  string
	Code  string
	Stale int
}

type config struct {
	Checks []check
	Alert  string
}

var c config

func handler() error {
	fmt.Println("We did it!")
	return nil
}

func loadConfig() {
	cf, err := s3.GetConfigFromEnv(&c)
	if err != nil {
		log.Print(err)
		panic(err)
	}
	cf.OnError = func(_ *s3.ConfigFile, err error) {
		log.Print(err)
	}
	cf.Autoreload(60)
}

func main() {
	loadConfig()

	d := mux.NewDispatcher(
		&mux.SimpleReceiver{
			HandleFunc: handleFunc,
			AuthFunc:   authFunc,
		},
	)
	mux.Start(d)
}
