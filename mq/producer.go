package mq

import (
	"errors"
	"fmt"

	"github.com/adjust/rmq/v4"

	mylog "matrix.works/async-proxy/logger"
)

type ProducerCreator struct {
	server     *RmqServerInfo
	connection rmq.Connection
	mqMap      map[string]rmq.Queue
}

func NewProducerCreator(host string, port, db int) (*ProducerCreator, error) {
	if host == "" || port == 0 {
		return nil, errors.New("Missed required parameters of MQ configure")
	}
	pc := &ProducerCreator{
		server: &RmqServerInfo{
			host: host,
			port: port,
			db:   db,
		},
	}
	return pc, nil
}

func (pc *ProducerCreator) Create(mqConnName string, mqIdList []string) (err error) {

	logger := mylog.GetStdoutLogger()

	mqAddress := fmt.Sprintf("%s:%d", pc.server.host, pc.server.port)
	logger.Printf("[ProducerCreator] mq address: %s/%d\n", mqAddress, pc.server.db)

	pc.connection, err = rmq.OpenConnection(mqConnName, "tcp", mqAddress, pc.server.db, nil)
	if err != nil {
		return err
	}
	if pc.connection == nil {
		return errors.New("Failed to open MQ connection")
	}
	logger.Printf("Open connection: %s", mqConnName)

	pc.mqMap = make(map[string]rmq.Queue)
	for _, mqId := range mqIdList {
		mq, err := pc.connection.OpenQueue(mqId)
		if err != nil {
			logger.Printf("Open queue failed: %s", err.Error())
			continue
		}
		pc.mqMap[mqId] = mq
		logger.Printf("Open queue: %s", mqId)
	}
	return err
}

func (pc *ProducerCreator) GetQueue(mqId string) rmq.Queue {
	if mq, exist := pc.mqMap[mqId]; exist {
		return mq
	} else {
		return nil
	}
}

func (pc *ProducerCreator) Cleanup() {
	cleaner := rmq.NewCleaner(pc.connection)
	cleaner.Clean()
	logger := mylog.GetStdoutLogger()
	logger.Printf("[ProducerCreator.Destroy] Clear MQ connection\n")
}
