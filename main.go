package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/akerl/go-lambda/apigw/events"
	"github.com/akerl/go-lambda/mux"
	"github.com/akerl/go-lambda/s3"
)

var (
	triggerRegex = regexp.MustCompile(`^/trigger$`)
	reportRegex  = regexp.MustCompile(`^/report$`)
	defaultRegex = regexp.MustCompile(`^/.*$`)
)

var config Config

func loadConfig() {
	cf, err := s3.GetConfigFromEnv(&config)
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
		mux.NewRoute(triggerRegex, triggerFunc),
		&mux.Route{
			Path: reportRegex,
			SimpleReceiver: mux.SimpleReceiver{
				HandleFunc: reportFunc,
				AuthFunc:   reportAuthFunc,
			},
		},
		mux.NewRoute(defaultRegex, defaultFunc),
		&mux.SimpleReceiver{
			HandleFunc: cronFunc,
			AuthFunc:   cronAuthFunc,
		},
	)
	mux.Start(d)
}

func reportFunc(req events.Request) (events.Response, error) {
	//TODO add report
	return events.Succeed("Placeholder")
}

func reportAuthFunc(req events.Request) (events.Response, error) {
	//TODO add report auth
	return events.Response{}, nil
}

func triggerFunc(req events.Request) (events.Response, error) {
	code := req.QueryStringParameters["code"]
	if code == "" {
		return events.Fail("No code provided")
	}

	check, found := config.CheckFromCode(code)
	if !found {
		return events.Fail("Invalid code")
	}

	err := config.WriteCheck(check)
	if err != nil {
		return events.Fail("Failed to write check")
	}
	return events.Succeed("Check updated")
}

func defaultFunc(req events.Request) (events.Response, error) {
	return events.Redirect("https://"+req.Headers["Host"]+"/report", 303)
}

func cronFunc(req events.Request) (events.Response, error) {
	for _, c := range config.Checks {
		ok, err := config.IsCheckStale(c)
		if err != nil {
			return events.Fail("Failed to parse check")
		}
		if !ok {
			err := config.Alert(c)
			if err != nil {
				return events.Fail("Failed to alert for check")
			}
		}
	}
	return events.Succeed("Successful cron run")
}

func cronAuthFunc(req events.Request) (events.Response, error) {
	if req.RequestContext.AccountID == "" {
		return events.Response{}, nil
	}
	return events.Response{}, fmt.Errorf("request not allowed via API Gateway")
}
