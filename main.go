package main

import (
	"fmt"
	"log"

	"github.com/akerl/go-lambda/s3"
	"github.com/aws/aws-lambda-go/lambda"
)

type config struct {
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
	lambda.Start(handler)
}
