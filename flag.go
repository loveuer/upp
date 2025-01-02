package upp

import (
	"flag"
	"time"
)

type __flag struct {
	debug  bool
	listen struct {
		http string
	}
}

var _flag = &__flag{}

func init() {
	time.Local = time.FixedZone("CST", 8*3600)

	flag.BoolVar(&_flag.debug, "debug", false, "debug mode")
	flag.StringVar(&_flag.listen.http, "listen.http", "localhost:8080", "")

	flag.Parse()
}
