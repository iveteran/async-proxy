package mq

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"matrix.works/fmx-async-proxy/conf"
	"matrix.works/fmx-common/mq"
)

type MqProducer struct {
	mq *mq.ProducerCreator
}

func NewMqProducer() *MqProducer {
	return &MqProducer{}
}

func createMQs(mqConnName string, routeList ...string) *mq.ProducerCreator {
	if _, exist := conf.Cfg.Redises["mq"]; !exist {
		log.Fatal("[createMQs] Error: The configure of MQ does not exist")
		return nil
	}
	redisMQ := conf.Cfg.Redises["mq"]
	host := redisMQ.Host
	port := redisMQ.Port
	db := redisMQ.Database

	mq, err := mq.NewProducerCreator(host, port, db)
	if err != nil {
		log.Fatal("[createMQs] Error: The configure of MQ does not exist")
		return nil
	}
	log.Printf("[createMQs] Create MQs [%s] %+v\n", mqConnName, routeList)
	mq.Create(mqConnName, routeList)

	CleanupStale(mq)

	return mq
}

func CleanupStale(mq *mq.ProducerCreator) {
	mq.Cleanup()
	log.Printf("[Cleanup] Clear MQ stale connection and queues\n")
}

func (this *MqProducer) Enqueue(
	topic, backend string, requestBytes []byte,
) (string, error) {
	msgId := generateMessageId()

	defaultQueue := DEFAULT_QUEUE_NAME
	if conf.Cfg.Mq.DefaultQueue.Name != "" {
		defaultQueue = conf.Cfg.Mq.DefaultQueue.Name
	}
	if this.mq == nil {
		topicList := getTopicList()
		topicList = append(topicList, defaultQueue)
		this.mq = createMQs(conf.Cfg.Mq.ConnectionName, topicList...)
	}

	delivery := fmt.Sprintf("{ \"ts\": %d, \"msg_id\": \"%s\", \"backend\": \"%s\", \"request\": \"%x\" }",
		time.Now().Unix(), msgId, backend, requestBytes)

	//fmt.Printf("[Enqueue] delivery: %+v\n", delivery)

	mq := this.mq.GetQueue(topic)
	fmt.Printf("[Enqueue] get mq object(topic: %s): %+v\n", topic, mq)
	if mq == nil {
		mq = this.mq.GetQueue(defaultQueue)
	}

	err := mq.Publish(delivery)
	if err != nil {
		return "", err
	}

	// TODO: report reqeust processing status

	return msgId, nil
}

func getTopicList() (result []string) {
	for topic := range conf.Cfg.Mq.TopicQueues {
		result = append(result, topic)
	}
	return result
}

func generateMessageId() string {
	id := uuid.New()
	return id.String()
}
