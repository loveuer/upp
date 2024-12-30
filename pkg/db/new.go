package db

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/glebarez/sqlite"
	"github.com/loveuer/upp/pkg/log"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(uri string) (*gorm.DB, error) {
	ins, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	var (
		username = ""
		password = ""
		tx       *gorm.DB
	)

	if ins.User != nil {
		username = ins.User.Username()
		password, _ = ins.User.Password()
	}

	switch ins.Scheme {
	case "sqlite":
		path := strings.TrimPrefix(uri, ins.Scheme+"://")
		log.Debug("db.New: type = %s, path = %s", ins.Scheme, path)
		tx, err = gorm.Open(sqlite.Open(path))
	case "mysql", "mariadb":
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", username, password, ins.Host, ins.Path, ins.RawQuery)
		log.Debug("db.New: type = %s, dsn = %s", ins.Scheme, dsn)
		tx, err = gorm.Open(mysql.Open(dsn))
	case "pg", "postgres", "postgresql":
		opts := make([]string, 0)
		for key, val := range ins.Query() {
			opts = append(opts, fmt.Sprintf("%s=%s", key, val))
		}
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s %s", ins.Hostname(), username, password, ins.Path, ins.Port(), strings.Join(opts, " "))
		log.Debug("db.New: type = %s, dsn = %s", ins.Scheme, dsn)
		tx, err = gorm.Open(postgres.Open(dsn))
	default:
		return nil, fmt.Errorf("invalid database type(uri_scheme): %s", ins.Scheme)
	}

	if err != nil {
		return nil, err
	}

	return tx, nil
}
