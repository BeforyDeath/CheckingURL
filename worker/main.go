package main

import (
	"fmt"
	"net/http"
	"net/rpc"
	"strings"
	"time"

	"github.com/BeforyDeath/CheckingURL/core"
	"github.com/BeforyDeath/CheckingURL/server"
	"github.com/BeforyDeath/CheckingURL/storage"
	log "github.com/Sirupsen/logrus"
	"github.com/satori/go.uuid"
)

var res server.Response
var req server.Request

var push = make(map[int]storage.SrvStatus)

var cCurl http.Client

func main() {

	UID := uuid.NewV4().String()
	log.Infof("Запуск worker:%v", UID)

	LoggerFile, err := core.Config.Init("config.json")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer LoggerFile.Close()

	cCurl = http.Client{
		Timeout: time.Duration(time.Millisecond * core.Config.Worker.Timeout),
	}

	client, err := rpc.DialHTTP("tcp", core.Config.Listen)
	if err != nil {
		log.Fatal(err)
	}

	req.UID = UID
	req.Quantity = core.Config.Worker.Quantity

	for {

		res.Sleep = false
		res.List = make(map[string]int)

		log.Infof("Запрашиваю:\t%v", core.Config.Worker.Quantity)
		err = client.Call("Server.Get", req, &res)
		if err != nil {
			log.Fatal(err)
		}

		if res.Sleep {
			time.Sleep(time.Second * 1)
			log.Warnf("Опс! Мне ответили `%v` сплю 1 сек", res.Msg)

		} else {
			log.Infof("Получил:\t%v", len(res.List))

			for v, k := range res.List {
				log.Infof("Обрабатываю id:%v \t url: %v", k, v)
				code := curl(v)
				log.Infof("Ответ:%v", code)
				push[k] = storage.SrvStatus{
					Url:  v,
					Code: code,
				}
			}
		}

		// todo push
		if len(push) >= 1 {
			log.Infof("Возвращаю: %v записей", len(push))

			req.SrvStatus = push
			err = client.Call("Server.Set", req, &res)
			if err != nil {
				log.Fatal(err)
			}
			for k := range push {
				delete(push, k)
			}
		}

	}
/*
	log.Info("Freeze ...")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, syscall.SIGTERM)
	for {
		select {
		case <-sigChan:
			log.Info("Exit (CTRL-C)")
			return
		}
	}
*/
}

func curl(url string) int {
	res, err := cCurl.Head(url)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			return 0
		}
		if strings.Contains(err.Error(), "Client.Timeout") {
			fmt.Println(err)
			return 503
		}
		return 1
	}
	return res.StatusCode
}
