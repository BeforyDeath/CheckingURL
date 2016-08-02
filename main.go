package main

import (
	"net"
	"net/http"
	"net/rpc"

	"github.com/BeforyDeath/CheckingURL/core"
	"github.com/BeforyDeath/CheckingURL/server"
	"github.com/BeforyDeath/CheckingURL/storage"
	log "github.com/Sirupsen/logrus"
	"time"
)

func main() {

	LoggerFile, err := core.Config.Init("config.json")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer LoggerFile.Close()

	storage.Connect()
	defer storage.Close()

	// запускаем раз в час удаление статусов, сроком больше 24 часов
	go RefreshEvery(1 * time.Hour, storage.StatusDelDay)

	Server := new(server.Server)
	err = rpc.Register(Server)
	if err != nil {
		log.Fatal(err)
	}

	rpc.HandleHTTP()

	rpc, err := net.Listen("tcp", core.Config.Listen)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Serving RPC handler %v", core.Config.Listen)
	err = http.Serve(rpc, nil)
	if err != nil {
		log.Fatal(err)
	}

}

func RefreshEvery(d time.Duration, f func() error) {
	for _ = range time.Tick(d) {
		err := f()
		if err != nil {
			log.Errorf("Refresh error: %v", err)
			break
		}
	}
}