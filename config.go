package main

import (
	"strconv"
	"strings"
	"time"

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

func (c Config) ReadCheck(check Check) (int64, error) {
	key := c.CheckPath + check.Code
	resp, err := s3.GetObject(c.Bucket, key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(string(resp), 10, 64)
}

func (c Config) IsCheckStale(check Check) (bool, error) {
	ts, err := c.ReadCheck(check)
	if err != nil {
		return false, err
	}
	if ts+check.Stale < time.Now().Unix() {
		return true, nil
	}
	return false, nil
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
