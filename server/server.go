package server

import (
	"github.com/BeforyDeath/CheckingURL/storage"
	log "github.com/Sirupsen/logrus"
)

type Request struct {
	Sleep     bool
	Quantity  int
	SrvStatus map[int]storage.SrvStatus
	UID       string
}

type Response struct {
	List  map[string]int
	Sleep bool
	Msg   string
}

type Server int

func (r *Server) Get(req *Request, res *Response) error {
	list, err := storage.Cache.Get(req.UID, req.Quantity)
	if err != nil {
		res.Sleep = true
		res.Msg = err.Error()
		log.Info(err)
		return nil
	}
	res.Sleep = false
	res.List = list
	return nil
}

func (r *Server) Set(req *Request, res *Response) error {
	storage.Update(req.UID, req.SrvStatus)
	storage.StatusUpdate(req.UID, req.SrvStatus)
	return nil
}
