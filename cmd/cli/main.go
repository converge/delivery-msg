package main

import (
	"delivery-msg"
	"delivery-msg/config"
	"github.com/charmbracelet/log"
	"github.com/nats-io/nats.go"
	"time"
)

func main() {

	cfg := config.ReadConfig()
	nc, err := nats.Connect(delivery_msg.NATSHostPost, nats.UserInfo(cfg.NATSUser, cfg.NATSPassword))
	if err != nil {
		log.Error(err)
		return
	}
	defer nc.Close()

	sub, err := nc.SubscribeSync("nats_development")
	if err != nil {
		log.Error(err)
		return
	}

	for {
		msg, err := sub.NextMsg(90 * time.Minute)
		if err != nil {
			log.Error(err)
			return
		}
		log.Infof("Received a message: (subject: %s) %s", msg.Subject, string(msg.Data))
	}

}
