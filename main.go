package main

import (
	"log"
	"regexp"

	"github.com/akerl/go-lambda/apigw/events"
	"github.com/akerl/go-lambda/mux"
	"github.com/akerl/go-lambda/s3"


var (
	triggerRegex = regexp.MustCompile(`^/trigger$`)
	reportRegex  = regexp.MustCompile(`^/report$`)
	defaultRegex = regexp.MustCompile(`^/.*$`)
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

	reportRoute := &mux.Route{
		Path: reportRegex,
		SimpleReceiver: mux.SimpleReceiver{
			HandleFunc: reportFunc,
			AuthFunc:   reportAuthFunc,
		},
	}

	d := mux.NewDispatcher(
		mux.NewRoute(triggerRegex, triggerFunc),
		reportRoute,
		mux.NewRoute(defaultRegex, defaultFunc),
		&mux.SimpleReceiver{
			HandleFunc: cronFunc,
			AuthFunc:   cronAuthFunc,
		},
	)
	mux.Start(d)
}

func reportFunc(req events.Request) (events.Response, error) {
}

func reportAuthFunc(req events.Request) (events.Response, error) {
}

func triggerFunc(req events.Request) (events.Response, error) {
}

func defaultFunc(req events.Request) (events.Response, error) {
	return events.Redirect("https://"+req.Headers["Host"] + "/report", 303)
}

func cronFunc(req events.Request) (events.Response, error) {
}

func cronAuthFunc(req events.Request) (events.Response, error) {
}
