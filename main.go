package main

import (
	cutils "common/utils"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/nats-io/go-nats-streaming"
)

const (
	PublisherType  = "publisher"
	SubscriberType = "subscriber"

	Subject = "test_trades"
)

var (
	NatsConn     stan.Conn
	g_actor_type = SubscriberType
	gClientName  string

	SysSymbols = []string{"BTC_USDT", "ETH_USDT"}
	PltSymbols = []string{"btcusdt", "ethusdt"}
	Directions = []string{"buy", "sell"}
)

type PublishTrade struct {
	Platform  string
	SysSymbol string
	PltSymbol string
	TradeId   int64
	Ts        int64
	Price     float64
	Amount    float64
	Direction string
}

func init() {
	const usage = "nats_streming_test [-a actor_type] [-c client_name]"
	flag.StringVar(&g_actor_type, "a", SubscriberType, usage)
	flag.StringVar(&gClientName, "c", "a", usage)

	rand.Seed(time.Now().UnixNano())
}

func publishTrade() {
	logrus.Debug("start publish trades")
	for i := 0; i < 2; i++ {
		go func() {
			for {
				pmsg := PublishTrade{
					Platform:  "binance",
					SysSymbol: SysSymbols[rand.Uint32()%2],
					PltSymbol: PltSymbols[rand.Uint32()%2],
					TradeId:   time.Now().UnixNano() / 1e6,
					Ts:        time.Now().Unix(),
					Price:     4000.1,
					Amount:    1000.1,
					Direction: Directions[rand.Uint32()%2],
				}
				pbody, err := json.Marshal(&pmsg)
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"err":          err,
						"PublishTrade": pmsg,
					}).Errorf("encode PublishTrade err")
					continue
				}
				err = NatsConn.Publish(Subject, pbody)
				if err != nil {
					logrus.WithError(err).Errorf("publish trade to natsd server err")
					continue
				}

				time.Sleep(10 * time.Millisecond)
			}
		}()
		time.Sleep(100 * time.Millisecond)
	}
}

func subscribeTrade() {
	NatsConn.Subscribe(Subject, subTradeCallback, stan.DurableName("durable_trade"))
	logrus.Debug("subscribeTrade success, wait for trades")
}

func subTradeCallback(msg *stan.Msg) {
	var trade PublishTrade
	err := json.Unmarshal(msg.MsgProto.Data, &trade)
	if err != nil {
		logrus.WithError(err).Errorf("unmarshal nats streaming trade msg err")
		return
	}

	logrus.Debugf("recieve trade: %v", trade)
}

func main() {
	defer cutils.MyRecovery()

	flag.Parse()
	logrus.SetLevel(logrus.DebugLevel)

	var err error
	switch g_actor_type {
	case PublisherType:
		NatsConn, err = stan.Connect("test-cluster", fmt.Sprintf("trades-pub-%s", gClientName))
		if err != nil {
			panic(err)
		}
		defer NatsConn.Close()

		publishTrade()
	case SubscriberType:
		NatsConn, err = stan.Connect("test-cluster", fmt.Sprintf("trades-sub-%s", gClientName))
		if err != nil {
			panic(err)
		}
		defer NatsConn.Close()

		subscribeTrade()
	default:
		panic("wrong actor_type, use nats_streaming_test [-a publisher[subscriber]]")
	}

	select {}
}
