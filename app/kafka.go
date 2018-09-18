package app

import (
	"github.com/Shopify/sarama"
	"github.com/domac/trans-broker/logger"
	"time"
)

const KafkaOutputFrequency = 500

type KafkaOutput struct {
	brokers  []string
	producer sarama.AsyncProducer
}

func NewKafkaProducer(kafkaBrokers []string) (*KafkaOutput, error) {
	c := sarama.NewConfig()
	c.Producer.RequiredAcks = sarama.WaitForLocal

	//c.Producer.Partitioner = sarama.NewHashPartitioner
	//c.Producer.Return.Successes = true
	//c.Producer.Return.Errors = true

	c.Producer.Compression = sarama.CompressionSnappy
	c.Producer.Flush.Frequency = KafkaOutputFrequency * time.Millisecond

	logger.Log().Infof("new kafka producer, brokers : %v", kafkaBrokers)
	producer, err := sarama.NewAsyncProducer(kafkaBrokers, c)
	if err != nil {
		return nil, err
	}

	out := &KafkaOutput{
		brokers:  kafkaBrokers,
		producer: producer,
	}

	//handleError should receive errors
	go out.handleError()

	return out, nil
}

func (out *KafkaOutput) handleError() {
	var err *sarama.ProducerError
	for {
		err = <-out.producer.Errors()
		if err != nil {
			logger.Log().Errorf("producer message error, partition:%d offset:%d key:%v valus:%s error(%v)", err.Msg.Partition, err.Msg.Offset, err.Msg.Key, err.Msg.Value, err.Err)
		}
	}
}

//数据写入
func (out *KafkaOutput) Write(topic string, packet *Packet) (err error) {
	data := packet.body
	out.producer.Input() <- &sarama.ProducerMessage{Topic: topic, Value: sarama.ByteEncoder(data)}
	return
}

//批量数据写入
func (out *KafkaOutput) BulkWrite(topic string, packets []*Packet) (err error) {
	for _, packet := range packets {
		data := packet.body
		out.producer.Input() <- &sarama.ProducerMessage{Topic: topic, Value: sarama.ByteEncoder(data)}
	}
	return
}
