package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/hanjm/natsbeat/config"
)

type NatsBeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	conf := config.DefaultConfig
	if err := cfg.Unpack(&conf); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}
	logp.Info("natsbeat using config:%+v", conf)
	bt := &NatsBeat{
		done:   make(chan struct{}),
		config: conf,
	}
	return bt, nil
}

func (bt *NatsBeat) Run(b *beat.Beat) error {
	logp.Info("natsbeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	err = Init(bt.config.NatsHost, bt.config.NatsMonitorPort, int(bt.config.Period.Seconds()))
	if err != nil {
		return err
	}

	for {
		select {
		case <-bt.done:
			return nil
		case stats := <-engine.StatsCh:
			event := beat.Event{
				Timestamp: time.Now(),
				Fields: common.MapStr{
					"type":           b.Info.Name,
					"hostname":       b.Info.Hostname,
					"cpu":            stats.Varz.CPU,
					"mem":            stats.Varz.Mem,
					"up_time":        stats.Varz.Uptime,
					"num_conns":      stats.Connz.NumConns,
					"in_msgs":        stats.Varz.InMsgs,
					"out_msgs":       stats.Varz.OutMsgs,
					"in_bytes":       stats.Varz.InBytes,
					"out_bytes":      stats.Varz.OutBytes,
					"slow_consumers": stats.Varz.SlowConsumers,
					"in_msgs_rate":   stats.Rates.InMsgsRate,
					"out_msgs_rate":  stats.Rates.OutMsgsRate,
					"in_bytes_rate":  stats.Rates.InBytesRate,
					"out_bytes_rate": stats.Rates.OutBytesRate,
				},
			}
			bt.client.Publish(event)
			//logp.Info("Event sent")
		}
	}
}

func (bt *NatsBeat) Stop() {
	Close()
	bt.client.Close()
	close(bt.done)
}
