package main

import (
	"strconv"
	"strings"

	"github.com/akerl/go-lambda/s3"
	s3api "github.com/aws/aws-sdk-go-v2/service/s3"
)

type Check struct {
	Name  string
	Code  string
	Stale int
}

type CheckSet []Check

type Config struct {
	Bucket   string
	Path     string
	Checks   CheckSet
	Alerts   []string
	s3client *s3api.S3
}

func (cs CheckSet) CheckFromCode(code string) (check, bool) {
	for _, c := range cs {
		if c.Code == code {
			return c, true
		}
	}
	return Check{}, false
}

func (c Config) loadS3Client() error {
	if c.s3client != nil {
		return nil
	}
	var err error
	c.s3client, err = s3.Client()
	return err
}

func (c Config) WriteCheck(check Check) error {
	if err := c.loadS3Client(); err != nil {
		return err
	}

	key := c.CheckPath + check.Code
	ts := strconv.FormatInt(time.Now().Unix())

	input := &s3api.PutObjectInput{
		Body:   strings.NewReader(ts),
		Bucket: &c.Bucket,
		Key:    &key,
	}
	s3req := c.s3client.PutObjectRequest(input)
	_, err = s3req.Send()
	return err
}
