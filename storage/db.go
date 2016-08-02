package storage

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/BeforyDeath/CheckingURL/core"
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

func Connect() {
	var err error
	Db, err = sql.Open(core.Config.DataBase.DriverName, core.Config.DataBase.DataSourceName)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = Db.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func Close() {
	Db.Close()
}

func Update(UID string, list map[int]SrvStatus) error {
	var str string = ""
	for id, _ := range list {
		str += strconv.Itoa(id) + ","
	}
	str = str[0 : len(str)-1]

	log.Infof("DB [worker:%v] Обновляю:%v", UID, len(list))

	dt := time.Now().Format("2006-01-02 15:04:05")
	_, err := Db.Exec("UPDATE url SET last_datetime='" + dt + "' WHERE id IN (" + str + ")")
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func GetList() ([]*Url, error) {
	log.Infof("DB читаю:\t%v", core.Config.Server.Quantity)

	rows, err := Db.Query("SELECT id, domain_id, link, Last_datetime FROM url WHERE Last_datetime < (NOW() - interval 5 MINUTE) LIMIT 1000")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*Url, 0)

	for rows.Next() {
		e := new(Url)
		err := rows.Scan(&e.Id, &e.Domain_id, &e.Link, &e.Last_datetime)
		if err != nil {
			return nil, err
		}
		result = append(result, e)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	log.Infof("DB получил:\t%v", len(result))
	return result, nil
}


