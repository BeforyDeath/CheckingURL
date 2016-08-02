package storage

import (
	"time"

	log "github.com/Sirupsen/logrus"
)

type SrvStatus struct {
	Url  string
	Code int
}

func StatusUpdate(UID string, list map[int]SrvStatus) error {
	log.Infof("STATUS [worker:%v] Сохраняю:%v", UID, len(list))

	sql := "INSERT INTO history(url_id, url, code, datetime) VALUES "
	value := []interface{}{}

	for id, status := range list {
		sql += "(?, ?, ?, ?),"
		value = append(value, id, status.Url, status.Code, time.Now())
	}
	sql = sql[0 : len(sql)-1]

	stmt, err := Db.Prepare(sql)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(value...)
	if err != nil {
		return err
	}
	return nil
}

func StatusDelDay() error {
	log.Info("STATUS Удаляю записи старше 24 часов")
	_, err := Db.Exec("DELETE FROM history WHERE datetime < DATE_SUB(NOW(), INTERVAL 24 HOUR)")
	if err != nil {
		return err
	}
	return nil
}
