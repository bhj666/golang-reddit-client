package persistance

import (
	"aws-example/config"
	"fmt"
	"github.com/jinzhu/gorm"
	"net/url"
)

func openConnection() *gorm.DB {

	dsn := url.URL{
		User:     url.UserPassword(config.DB_USER, config.DB_PASSWORD),
		Scheme:   config.DB_SCHEME,
		Host:     fmt.Sprintf("%s:%d", config.DB_HOST, config.DB_PORT),
		Path:     "postgres",
		RawQuery: (&url.Values{"sslmode": []string{"disable"}}).Encode(),
	}
	db, err := gorm.Open("postgres", dsn.String())
	if err != nil {
		panic(err)
	}
	return db
}
