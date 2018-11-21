package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/akerl/go-lambda/apigw/events"
	"github.com/akerl/go-lambda/mux"
	"github.com/akerl/go-lambda/s3"
	s3api "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

const (
	checkPath = "checks/"
)

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
var bucket string

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
	bucket = os.Getenv("S3_BUCKET")

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
	code := req.QueryStringParameters["code"]
	if code == "" {
		return events.Fail("No code provided")
	}

	if _, err := uuid.Parse(code); err != nil {
		return events.Fail("Invalid code")
	}

	s3client, err := s3.Client()
	if err != nil {
		return events.Fail("Failed to load S3 client")
	}

	key := checkPath + code

	data := string(time.Now().Unix())

	input := &s3api.PutObjectInput{
		Body:   bytes.NewReader([]byte(data)),
		Bucket: &bucket,
		Key:    &key,
	}
	s3req := s3client.PutObjectRequest(input)
	_, err = s3req.Send()
	if err != nil {
		return events.Fail("Failed to post")
	}
	return events.Succeed("Updated!")
}

func defaultFunc(req events.Request) (events.Response, error) {
	return events.Redirect("https://"+req.Headers["Host"]+"/report", 303)
}

func cronFunc(req events.Request) (events.Response, error) {
}

func cronAuthFunc(req events.Request) (events.Response, error) {
	if req.RequestContext.AccountID == "" {
		return events.Response{}, nil
	}
	return events.Response{}, fmt.Errorf("request not allowed via API Gateway")
}
