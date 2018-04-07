package beater

import (
	"fmt"

	natstop "github.com/nats-io/nats-top/util"
)

var engine *natstop.Engine

func Init(host string, port int, delay int) (err error) {
	engine = natstop.NewEngine(host, port, 1024, delay)
	engine.SetupHTTP()
	_, err = engine.Request("/varz")
	if err != nil {
		return fmt.Errorf("nats monitor api endpoint connect error:%s, host:%s, port:%d", err, host, port)
	}
	go engine.MonitorStats()
	return nil
}

func Close() {
	if engine != nil {
		close(engine.ShutdownCh)
	}
}
