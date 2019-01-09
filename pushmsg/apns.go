package pushmsg

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/lauyang/golib/logs"
)

type APNS struct {
	server      string
	port        int
	certFile    string
	keyFile     string
	locker      sync.Mutex
	isConnected bool
	conn        net.Conn
	channel     *tls.Conn
	index       int
}

// 初始化
func (self *APNS) Init(server string, port int, certFile string, keyFile string) {
	self.server = server
	self.port = port
	self.certFile = certFile
	self.keyFile = keyFile
}

// 执行连接
func (self *APNS) connect() error {
	// 连接成功
	if self.isConnected {
		return nil
	}

	// 执行连接
	cert, err := tls.LoadX509KeyPair(self.certFile, self.keyFile)
	if nil != err {
		logs.Error(err)
		return err
	}

	conf := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: false,
		ServerName:         self.server,
	}

	self.conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", self.server, self.port))
	if nil != err {
		logs.Error(err)
		return err
	}

	//  创建通道
	self.channel = tls.Client(self.conn, conf)

	//  发送握手
	err = self.channel.Handshake()
	if nil != err {
		logs.Error(err)
		return err
	}

	// 调试信息
	state := self.channel.ConnectionState()
	logs.Debug("conn state", state)

	self.isConnected = true
	logs.Info("connect to apns server")

	go self.recv()
	return nil
}

// 接收协程
func (self *APNS) recv() {
	// 接收错误
	defer func() {
		// 加锁
		self.locker.Lock()
		defer self.locker.Unlock()

		self.isConnected = false
		self.channel.Close()
		self.conn.Close()

		self.conn = nil
		self.channel = nil

		if err := recover(); nil != err {
			logs.Error(err)
		}
	}()

	// 进入消息接循环
	for {
		readb := [1024]byte{}
		n, err := self.channel.Read(readb[:])
		if n < 1 && nil != err {
			logs.Error(err)
			break
		}

		logs.Debug(readb[:n])
	}
}

// 发布消息
func (self *APNS) Push(token string, msg string) error {
	// 加锁
	self.locker.Lock()
	defer self.locker.Unlock()

	err := self.connect()
	if nil != err {
		return err
	}

	// 发送
	// 构造 pdu
	buffer := bytes.NewBuffer([]byte{})
	frame := bytes.NewBuffer([]byte{})

	btoken, _ := hex.DecodeString(token)
	if len(btoken) != 32 {
		tmp := make([]byte, 32)
		copy(tmp[:], btoken)
		btoken = tmp
	}

	bpayload := []byte(msg)
	// create frame
	// device token
	binary.Write(frame, binary.BigEndian, uint8(1))
	binary.Write(frame, binary.BigEndian, uint16(len(btoken)))
	binary.Write(frame, binary.BigEndian, btoken)
	//frame.Write(btoken)
	// playload
	binary.Write(frame, binary.BigEndian, uint8(2))
	binary.Write(frame, binary.BigEndian, uint16(len(bpayload)))
	binary.Write(frame, binary.BigEndian, bpayload)
	//frame.Write(bpayload)
	// identifier
	binary.Write(frame, binary.BigEndian, uint8(3))
	binary.Write(frame, binary.BigEndian, uint16(4))
	self.index++
	binary.Write(frame, binary.BigEndian, uint32(self.index))
	// expiration
	binary.Write(frame, binary.BigEndian, uint8(4))
	binary.Write(frame, binary.BigEndian, uint16(4))
	binary.Write(frame, binary.BigEndian, uint32(time.Now().UTC().Unix()+60*60))

	bbody := frame.Bytes()
	// 命令
	binary.Write(buffer, binary.BigEndian, uint8(2))
	binary.Write(buffer, binary.BigEndian, uint32(len(bbody)))
	binary.Write(buffer, binary.BigEndian, bbody)

	pdu := buffer.Bytes()

	_, err = self.channel.Write(pdu)
	return err
}
