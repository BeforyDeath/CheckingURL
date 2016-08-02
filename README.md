## Checking the access URL

ПС .. Пока в один поток, без учёта параллельных запросов

Импортируем зависимости
```
go get
```

Создаём таблицы 
```
CREATE TABLE `domain` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(45) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `history` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `url_id` int(11) NOT NULL,
  `url` varchar(512) NOT NULL,
  `code` int(11) NOT NULL,
  `datetime` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_history_url1_idx` (`url_id`),
  CONSTRAINT `fk_history_url1` FOREIGN KEY (`url_id`) REFERENCES `url` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `url` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `domain_id` int(11) NOT NULL,
  `link` varchar(512) NOT NULL,
  `last_datetime` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_url_domain_idx` (`domain_id`),
  CONSTRAINT `fk_url_domain` FOREIGN KEY (`domain_id`) REFERENCES `domain` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
```

Далее заливаем дамп из `migrate/dump.sql` или запускаем 
```
go run migrate/main.go
```

Cобираем и запускаем приложение
```
go build
./CheckingURL
```

Отрываем пару консолей и запускаем в них воркеры
```
go run worker/main.go
```
