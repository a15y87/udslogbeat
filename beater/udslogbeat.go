package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/a15y87/udslogbeat/config"
	"net"
)

type Udslogbeat struct {
	done     chan struct{}
	config   config.Config
	client   publisher.Client
	udsocket net.Listener
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}
	udsocket, err := net.Listen("unix", config.SocketPath)
	if err != nil {
		return nil, fmt.Errorf("Error create unix domain socket: %v", err)
	}

	bt := &Udslogbeat{
		done:     make(chan struct{}),
		config:   config,
		udsocket: udsocket,
	}
	return bt, nil
}

func (bt *Udslogbeat) Run(b *beat.Beat) error {
	logp.Info("udslogbeat is running! Hit CTRL-C to stop it.")

	bt.client = b.Publisher.Connect()
	buf := make([]byte, bt.config.MaxMessageSize)

	ticker := time.NewTicker(bt.config.Period)
	counter := 1
	for {
		select {
		case <-bt.done:
			return nil
		default:
		}
		fd, err := bt.udsocket.Accept()
		if err != nil {
			logp.Err(err.Error())
			return err
		}
		nr, err := fd.Read(buf)
		message := string(buf[0:nr])
		fd.Close()
		logp.Info(message)

		event := common.MapStr{
			"@timestamp": common.Time(time.Now()),
			"type":       b.Name,
			"message":    string(message[:]),
			"counter":    counter,
		}
		bt.client.PublishEvent(event)
		logp.Info("Event sent")
		counter++
	}
}

func (bt *Udslogbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
