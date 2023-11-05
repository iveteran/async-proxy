package mq

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/adjust/rmq/v4"

	"matrix.works/async-proxy/backend"
	"matrix.works/async-proxy/conf"
	"matrix.works/async-proxy/logger"
)

const (
	DEFAULT_QUEUE_NAME    = "fap_default"
	DEFAULT_UNACKED_LIMIT = 10
	DEFAULT_NUM_CONSUMERS = 1
	DEFAULT_MESSAGE_TTL   = 60

	STATUS_QUEUING = 1
)

type Consumer struct {
	name  string
	count int
}

type MqConsumer struct {
	Consumer
	queueInfo *conf.QueueInfo
	creator   *MqConsumerCreator
}

func NewMqConsumer(
	mqName string, queueInfo *conf.QueueInfo,
	c *MqConsumerCreator,
) *MqConsumer {
	return &MqConsumer{
		Consumer: Consumer{
			name:  mqName,
			count: 0,
		},
		queueInfo: queueInfo,
		creator:   c,
	}
}

func (c *MqConsumer) Consume(delivery rmq.Delivery) {
	c.count++
	logger.Logger.Printf("[Consume] connection %s queue %s consumed %d messages\n",
		*c.creator.CC.GetConnection(), c.name, c.count)
	logger.Logger.Printf("[Consume] queue info: %+v\n", c.queueInfo)

	//fmt.Printf(">>> delivery payload: %+v\n", delivery.Payload())
	results := make(map[string]interface{})
	err := json.Unmarshal([]byte(delivery.Payload()), &results)
	if err != nil {
		logger.Logger.Printf("[Consume] Decode mq message failed: %s\n", err.Error())
		return
	}
	//fmt.Printf(">>> results: %v\n", results)

	msgTs := int64(results["ts"].(float64))
	msgId := results["msg_id"].(string)
	beHost := results["backend"].(string)
	requestBytesStr := results["request"].(string)

	if msgId == "" || beHost == "" || len(requestBytesStr) == 0 {
		logger.Logger.Printf("[Consume] Have invalid parameter")
		return
	}

	nowTs := time.Now().Unix()
	if nowTs-msgTs > c.queueInfo.MessageTTL {
		logger.Logger.Printf("[Consume] the message(%s) is timeout, drop it", msgId)
		return
	}

	requestBytes, err := hex.DecodeString(requestBytesStr)
	if err != nil {
		fmt.Println("[Consume] Unable to convert hex to byte. ", err)
	}
	//fmt.Printf(">>> requestBytes: %v\n", requestBytes)

	var status int
	var errmsg string
	err = backend.SendMessage(beHost, requestBytes)
	if err == nil {
		status = 1
		logger.Logger.Printf("[Info] Successed to send request to backend server\n")
		if err := delivery.Ack(); err != nil {
			logger.Logger.Printf("[Error] Ack delivery failed: %s\n", err.Error())
		}
	} else {
		status = -1
		errmsg = err.Error()
		logger.Logger.Printf("[Error] Failed to send request to backend server: %s\n", err.Error())
		if err := delivery.Reject(); err != nil {
			logger.Logger.Printf("[Error] Reject delivery failed: %s\n", err.Error())
		}
	}
	fmt.Printf("[Consume] Send message, status: %d, errmsg: %s\n", status, errmsg)

	// TODO: report reqeust processing status
}

type MqConsumerCreator struct {
	CC *ConsumerCreator
}

func (c *MqConsumerCreator) Init(mqConnName string) error {
	const mqCfgName = "mq"
	if _, exist := conf.Cfg.Redises[mqCfgName]; !exist {
		return errors.New("The configure of message queue server does not exist")
	}
	host := conf.Cfg.Redises[mqCfgName].Host
	port := conf.Cfg.Redises[mqCfgName].Port
	db := conf.Cfg.Redises[mqCfgName].Database
	if host == "" || port == 0 {
		return errors.New("The host or port of message queue server are invalid")
	}
	mqCfg := conf.Cfg.Mq

	queues := mqCfg.TopicQueues
	if mqCfg.DefaultQueue.Name == "" {
		mqCfg.DefaultQueue = conf.QueueInfo{
			Name:         DEFAULT_QUEUE_NAME,
			UnackedLimit: DEFAULT_UNACKED_LIMIT,
			NumConsumers: DEFAULT_NUM_CONSUMERS,
			MessageTTL:   DEFAULT_MESSAGE_TTL,
		}
	}
	queues[mqCfg.DefaultQueue.Name] = mqCfg.DefaultQueue

	mqConsumers := make(map[string]rmq.Consumer)
	for topic, queueInfo := range queues {
		fmt.Printf("[Init] create consumer for: topic: %s, queue info: %+v\n", topic, queueInfo)
		mqConsumers[topic] = NewMqConsumer(topic, &queueInfo, c)
	}

	var err error
	c.CC, err = NewConsumerCreator(host, port, db,
		mqCfg.DefaultQueue.UnackedLimit, mqCfg.DefaultQueue.NumConsumers)
	if err != nil {
		return err
	}

	err = c.CC.Create(mqConnName, mqConsumers)
	if err != nil {
		return err
	}

	c.CC.Cleanup()

	return nil
}

func (c *MqConsumerCreator) Cleanup() {
	c.CC.StopConsuming()
	c.CC.Cleanup()
	logger.Logger.Printf("[MqConsumerCreator.Cleanup] Clear connections and queues")
}
