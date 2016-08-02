package storage

import (
	"errors"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
)

var Cache cache

type cache struct {
	sync.Mutex
	List  map[string]int
	Sleep bool
}

type Url struct {
	Id            int
	Domain_id     int
	Link          string
	Last_datetime time.Time
}

func (c *cache) Get(UID string, quantity int) (list map[string]int, err error) {

	if len(c.List) == 0 {
		if c.Sleep == false {
			c.get()
		}
		return nil, errors.New("CACHE [worker:" + UID + "] Подождите пока обновлюсь")
	}

	log.Infof("CACHE отрезаю:%v", quantity)
	i := 0
	c.Lock()
	list = make(map[string]int)
	for v, k := range c.List {
		if i >= quantity {
			break
		}
		i++
		list[v] = k
		delete(c.List, v)
	}
	c.Unlock()

	log.Infof("CACHE [worker:%v] Отдаю:%v осталось:%v", UID, len(list), len(c.List))
	return list, nil
}

func (c *cache) get() error {
	u, err := GetList()
	if err != nil {
		return err
	}
	if len(u) > 0 {
		log.Infof("CACHE кэширую:%v", len(u))
		c.Lock()
		c.List = make(map[string]int)
		for _, url := range u {
			c.List[url.Link] = url.Id
		}
		c.Sleep = false
		c.Unlock()
	} else {
		log.Info("CACHE в DB нечё нет, спю 1 минуту")
		c.Lock()
		c.Sleep = true
		c.Unlock()

		timer := time.NewTimer(time.Minute)
		go func() {
			<-timer.C
			log.Info("CACHE проснулись!")
			c.Lock()
			c.Sleep = false
			c.Unlock()
		}()
	}
	return nil
}
