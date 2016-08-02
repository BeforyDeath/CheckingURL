package core

import (
	"encoding/json"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	"time"
)

var Config config

type config struct {
	Listen   string
	Logger   logger
	DataBase database
	Server   server
	Worker   worker
}

type logger struct {
	Debug   bool
	OutFile bool
	File    *os.File
}

type database struct {
	DriverName     string
	DataSourceName string
}

type worker struct {
	Quantity int
	Timeout  time.Duration
}

type server struct {
	Quantity int
}

func (c *config) Init(fileName string) (*os.File, error) {

	f, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(f, &c); err != nil {
		return nil, err
	}

	if c.Logger.Debug {
		log.SetLevel(log.DebugLevel)
	}

	if c.Logger.OutFile {
		c.Logger.File, err = os.OpenFile("logrus.json", os.O_CREATE|os.O_RDWR, 0666) //os.O_APPEND |
		if err != nil {
			return nil, err
		}
		log.SetOutput(c.Logger.File)
		log.SetFormatter(&log.JSONFormatter{})
		return c.Logger.File, nil
	}

	return nil, nil
}
