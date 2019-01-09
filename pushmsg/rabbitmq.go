package pushmsg

import (
	"errors"
	"sync"

	"github.com/streadway/amqp"

	"github.com/lauyang/golib/logs"
)

type RabbitMQ struct {
	uri         string
	locker      sync.Mutex
	conn        *amqp.Connection
	channel     *amqp.Channel
	isConnected bool
}

// 初始化消息推送
func (self *RabbitMQ) Init(uri string) {
	self.uri = uri
}

// 执行消息推送
func (self *RabbitMQ) connect() error {
	self.isConnected = false
	// 关闭旧的
	if nil != self.channel {
		self.channel.Close()
		self.channel = nil
	}

	if nil != self.conn {
		self.conn.Close()
		self.conn = nil
	}

	// 开始连接
	var err error
	self.conn, err = amqp.Dial(self.uri)
	if nil != err {
		logs.Error(err)
		return err
	}

	// 创建通道
	self.channel, err = self.conn.Channel()
	if nil != err {
		logs.Error(err)
		return err
	}

	self.isConnected = true
	return nil
}

// 发送push消息
func (self *RabbitMQ) Push(queue string, msg string) error {
	// 加锁
	self.locker.Lock()
	defer self.locker.Unlock()
	// 判断连接
	if false == self.isConnected {
		err := self.connect()
		if nil != err {
			return err
		}
	}

	if false != self.isConnected {
		err := self.channel.Publish("", queue, true, false, amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte(msg),
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
			// a bunch of application/implementation-specific fields
		})

		if nil != err {
			self.isConnected = false
			logs.Error(err)
		}

		return err
	}

	return errors.New("not connect rabbitmq server")
}
