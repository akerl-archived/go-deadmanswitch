package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/akerl/go-lambda/s3"
	s3api "github.com/aws/aws-sdk-go-v2/service/s3"
)

// Check defines an individual item to monitor
type Check struct {
	Name  string
	Code  string
	Stale int64
}

// CheckSet defines a set of checks
type CheckSet []Check

// Config defines a set of checks and their metadata
type Config struct {
	Bucket   string
	Path     string
	Checks   CheckSet
	Alerts   []string
	s3client *s3api.S3
}

// CheckFromCode finds a check matching the given code
func (c Config) CheckFromCode(code string) (Check, bool) {
	for _, check := range c.Checks {
		if check.Code == code {
			return check, true
		}
	}
	return Check{}, false
}

// Alert sends an alert for the failing check
func (c Config) Alert(check Check) error {
	//TODO handle alerts
	return nil
}

func (c Config) loadS3Client() error {
	if c.s3client != nil {
		return nil
	}
	var err error
	c.s3client, err = s3.Client()
	return err
}

// ReadCheck returns the last timestamp for the check
func (c Config) ReadCheck(check Check) (int64, error) {
	key := c.Path + check.Code
	resp, err := s3.GetObject(c.Bucket, key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(string(resp), 10, 64)
}

// IsCheckStale reads a check timestamp and compares it to the stale window
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

// WriteCheck updates the stored timestamp
func (c Config) WriteCheck(check Check) error {
	if err := c.loadS3Client(); err != nil {
		return err
	}

	key := c.Path + check.Code
	ts := strconv.FormatInt(time.Now().Unix(), 10)

	input := &s3api.PutObjectInput{
		Body:   strings.NewReader(ts),
		Bucket: &c.Bucket,
		Key:    &key,
	}
	s3req := c.s3client.PutObjectRequest(input)
	_, err := s3req.Send()
	return err
}
