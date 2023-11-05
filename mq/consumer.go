package mq

import (
	"errors"
	"fmt"
	"time"

	"github.com/adjust/rmq/v4"

	mylog "matrix.works/async-proxy/logger"
)

type ConsumingInfo struct {
	unackedLimit int64
	numConsumers int64
}

type ConsumerCreator struct {
	server     *RmqServerInfo
	consuming  *ConsumingInfo
	connection rmq.Connection
}

func NewConsumerCreator(host string, port, db int, unackedLimit, numConsumers int64,
) (*ConsumerCreator, error) {
	if host == "" || port == 0 {
		return nil, errors.New("The host or port of message queue server are invalid")
	}

	cc := &ConsumerCreator{
		server: &RmqServerInfo{
			host: host,
			port: port,
			db:   db,
		},
		consuming: &ConsumingInfo{
			unackedLimit: unackedLimit,
			numConsumers: numConsumers,
		},
	}
	return cc, nil
}

func (cc *ConsumerCreator) Create(mqConnName string, mqConsumers map[string]rmq.Consumer) (err error) {

	logger := mylog.GetStdoutLogger()

	mqAddress := fmt.Sprintf("%s:%d", cc.server.host, cc.server.port)
	logger.Printf("[ConsumerCreator.Create] mq address: %s/%d\n", mqAddress, cc.server.db)

	cc.connection, err = rmq.OpenConnection(mqConnName, "tcp", mqAddress, cc.server.db, nil)
	if err != nil {
		return err
	}
	if cc.connection == nil {
		return errors.New("Failed to open MQ connection")
	}
	logger.Printf("[ConsumerCreator.Create] Open connection: %s", mqConnName)

	for mqName, mqConsumer := range mqConsumers {
		q, err := cc.connection.OpenQueue(mqName)
		if err != nil {
			logger.Printf("Open queue failed: %s", err.Error())
			continue
		}
		logger.Printf("[ConsumerCreator.Create] Open queue: %s", mqName)
		q.StartConsuming(cc.consuming.unackedLimit, time.Second)
		for i := int64(0); i < cc.consuming.numConsumers; i++ {
			cName := fmt.Sprintf("%s_consumer#%d", mqName, i)
			q.AddConsumer(cName, mqConsumer)
			logger.Printf("Add consumer: %s", cName)
		}
	}

	logger.Printf("[ConsumerCreator.Create] Waiting on %d queues...", len(mqConsumers))
	return err
}

func (cc *ConsumerCreator) GetConnection() *rmq.Connection {
	return &cc.connection
}

func (cc *ConsumerCreator) StopConsuming() {
	finChan := cc.connection.StopAllConsuming()
	<-finChan
}

func (cc *ConsumerCreator) Cleanup() {
	cleaner := rmq.NewCleaner(cc.connection)
	cleaner.Clean()

	logger := mylog.GetStdoutLogger()
	logger.Printf("[ConsumerCreator destroy] Clear MQ stale connections and queues\n")
}
