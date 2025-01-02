package upp

import (
	"time"

	"github.com/alexflint/go-arg"
	"github.com/loveuer/upp/pkg/tool"
)

type _env struct {
	Debug      bool   `arg:"env:DEBUG"`
	ListenHttp string `arg:"env:LISTEN_HTTP"`
}

var env = &_env{}

func init() {
	time.Local = time.FixedZone("CST", 8*3600)

	arg.MustParse(env)

	if env.Debug {
		tool.TablePrinter(env)
	}
}
