package main

import (
	"bufio"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/BeforyDeath/CheckingURL/core"
	"github.com/BeforyDeath/CheckingURL/storage"
	log "github.com/Sirupsen/logrus"
)

// (https?:\/\/)?([\w-\.]+\.[a-z]{2,6}\.?)

func main() {
	//*/
	LoggerFile, err := core.Config.Init("config.json")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer LoggerFile.Close()

	storage.Connect()
	defer storage.Close()

	//*/

	file, _ := os.Open("migrate/list.txt")
	defer file.Close()

	f := bufio.NewReader(file)
	d, _ := regexp.Compile(`([\w-\.]+\.[a-z]{2,6}\.?)`)

	domains := make(map[string]int)
	urls := make(map[string]int)
	var domain_id int
	var url_id int

	for {
		link, err := f.ReadString('\n')
		if err != nil {
			break
		}

		domain := d.FindString(link)
		if _, ok := domains[domain]; !ok {
			domain_id++
			domains[domain] = domain_id
		}

		link = strings.Replace(link, "\n", "", -1)
		if _, ok := urls[link]; !ok {
			url_id++
			urls[link] = domains[domain]
		}

		log.Infof("%v |%v| %v |%v|", domain, domains[domain], link, urls[link])
	}
	dump_domain(domains)
	dump_url(urls)

}

func dump_domain(domains map[string]int) {
	sql := "INSERT INTO domain(id, name) VALUES "
	value := []interface{}{}

	for name, id := range domains {
		sql += "(?, ?),"
		value = append(value, id, name)
	}
	sql = sql[0 : len(sql)-1]

	stmt, err := storage.Db.Prepare(sql)
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(value...)
	if err != nil {
		log.Fatal(err)
	}
}

func dump_url(urls map[string]int) {
	sql := "INSERT INTO url(link, domain_id, last_datetime) VALUES "
	value := []interface{}{}

	for link, domain_id := range urls {
		sql += "(?, ?, ?),"
		value = append(value, link, domain_id, time.Now())
	}
	sql = sql[0 : len(sql)-1]

	stmt, err := storage.Db.Prepare(sql)
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(value...)
	if err != nil {
		log.Fatal(err)
	}
}
