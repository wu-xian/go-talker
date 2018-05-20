package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/wu-xian/go-talker/common"
	"github.com/wu-xian/go-talker/proto"
	"github.com/wu-xian/go-talker/terminal"

	"github.com/wu-xian/go-talker/log"

	"gopkg.in/ini.v1"
)

var (
	IP                 string
	Port               int
	UName              string
	Key                string
	connection         *net.TCPConn
	messageChan        chan *proto.MessageWarpper = make(chan *proto.MessageWarpper, 0)
	shouldStopChan     chan struct{}              = make(chan struct{}, 1)
	messagePublishChan chan string                = make(chan string, 0)
	stopIt             chan os.Signal             = make(chan os.Signal, 1)
	ClientId           int32
)

type EmptyStruct struct{}

const MESSAGE_MAX_LENGTH = 2048

func main() {
	err := Init()
	if err != nil {
		fmt.Println(err)
		return
	}
	dialer := net.Dialer{
		Timeout: time.Duration(5) * time.Second,
	}
	conn, err := dialer.Dial("tcp", IP+":"+strconv.Itoa(Port))
	connection = conn.(*net.TCPConn)
	defer connection.Close()
	log.Logger.Info("has been connected to the server...")
	if err != nil {
		log.Logger.Info("unable to connect to the server : %s:%d", IP, Port)
		return
	}

	err = Login(connection)
	if err != nil {
		fmt.Println(err)
		return
	}

	go MessageReceiver(connection)
	go MessagePublisher(connection)

	go func() {
		signal.Notify(stopIt, os.Interrupt, os.Kill)
	}()

	terminal.LoopClientUI(messageChan, messagePublishChan)

	shouldStopChan <- EmptyStruct{}
	_ = Logout(connection)
	//connection.CloseRead()
	fmt.Println("application stopped")
}

func MessageReceiver(conn *net.TCPConn) {
	bytes := make([]byte, MESSAGE_MAX_LENGTH)
	for {
		count, err := conn.Read(bytes)
		//conn.CloseRead()
		if err != nil {
			log.Logger.Info("unable to read message", err)
			return
		}
		message := MessageInterpreter(bytes[:count])
		if err != nil {
			log.Logger.Info("invalid message:", err)
			return
		}
		log.Logger.Info("messagereceiver get:", message)
		if message == nil {
			continue
		}
		switch message.Type {
		case proto.COMMUNICATION_TYPE_ClientReceived:
			{
				messageChan <- message
			}
		case proto.COMMUNICATION_TYPE_ClientLogin:
			{
				messageChan <- message
			}
		case proto.COMMUNICATION_TYPE_ClientLogout:
			{
				if message.MessageClientLogout.Id != ClientId {
					messageChan <- message
				}
			}
		case proto.COMMUNICATION_TYPE_OnlineClients:
			{
				messageChan <- message
			}
		default:
			{
				log.Logger.Info("default message.", message.Type)
			}
		}
	}
}

func MessagePublisher(conn *net.TCPConn) {
	for {
		select {
		case content := <-messagePublishChan:
			{

				log.Logger.Info("send content", content)
				//content := ""
				//fmt.Scanln(&content)
				message := proto.MessageWarpper{
					Type: proto.COMMUNICATION_TYPE_ClientSend,
					MessageClientSend: &proto.MessageClientSend{
						Content: content,
					},
				}
				bytess, err := message.MessageMarshal()
				log.Logger.Info("send bytes ", bytess)
				if err != nil {
					return
				}
				_, err = conn.Write(bytess)
				//err = conn.CloseWrite()
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		case _ = <-shouldStopChan:
			{
				return
			}
		}
	}
}

//MessageInterpreter 获取包装壳
func MessageInterpreter(bytes []byte) (msg *proto.MessageWarpper) {
	warpper := &proto.MessageWarpper{}
	err := warpper.MessageUnmarshal(bytes)
	if err != nil {
		log.Logger.Info("err warpper:", err)
		return nil
	}
	return warpper
}

//Login 客户端登陆
func Login(conn *net.TCPConn) error {
	loginMessage := proto.MessageWarpper{
		Type: proto.COMMUNICATION_TYPE_LoginRequest,
		MessageLoginRequest: &proto.MessageLoginRequest{
			Name:  UName,
			Token: "",
		},
	}
	bytes, err := loginMessage.MessageMarshal()
	if err != nil {
		log.Logger.Info("err login message", err)
		return err
	}
	for i := 0; i < 3; i++ {
		count, err := conn.Write(bytes)
		common.CheckError(err)
		if count >= MESSAGE_MAX_LENGTH {
			return errors.New("message too large")
		}

		readBytes := make([]byte, 20)
		count, err = conn.Read(readBytes)
		msg := MessageInterpreter(readBytes[:count])
		if msg.Type == proto.COMMUNICATION_TYPE_LoginResponse &&
			msg.MessageLoginResponse.Succeed {
			ClientId = msg.MessageLoginResponse.Id
			break
		}
		if i == 2 {
			return errors.New("wrong login message")
		}
		time.Sleep(time.Second * 3)
	}
	log.Logger.Info("login succeeded")
	return nil
}

//Login 客户端登出
func Logout(conn *net.TCPConn) error {
	logoutMessage := &proto.MessageWarpper{
		Type:                 proto.COMMUNICATION_TYPE_LogoutRequest,
		MessageLogoutRequest: &proto.MessageLogoutRequest{},
	}
	bytes, err := logoutMessage.MessageMarshal()
	if err != nil {
		log.Logger.Info("logout err:", err)
		return err
	}
	conn.Write(bytes)

	return nil
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func Init() error {
	cfg, err := ini.Load("talker.conf")
	if err != nil {
		return errors.New("failure to load config file:talker.conf")
	}
	serverSection := cfg.Section("server")
	IP = serverSection.Key("ip").String()
	Port, err = serverSection.Key("port").Int()
	if err != nil {
		return errors.New("failure to load config : server.port")
	}
	clientSection := cfg.Section("client")
	UName = clientSection.Key("name").String()
	if UName == "" {
		panic("name is empty")
	}
	Key = clientSection.Key("key").String()
	return nil
}
