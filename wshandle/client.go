package wshandle

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type WsConn struct {
	Connecting   bool
	Connected    bool            // 连接状态
	MsgSendCh    chan string     // 消息发送通道
	MsgReceiveCh chan string     // 消息接收通道
	Done         chan struct{}   // 发送结束通道
	Exit         chan bool       // 程序退出信号
	Interrupt    chan os.Signal  // 终端中止信号
	Conn         *websocket.Conn // 主连接
	ConnMux      sync.Mutex      // 连接互斥锁
}

var wsconn *WsConn

func (c *WsConn) keepAlive() {
	for {
		if !c.Connected && !c.Connecting {
			c.Connecting = true
			New()
		}
	}
}

func (c *WsConn) messageReceiveHandler() {
	go c.keepAlive()
	defer close(c.Done)
	for {
		if c.Connected {
			_, msg, err := c.Conn.ReadMessage()
			if err != nil {
				// 读取信息出错，连接已经意外断开
				log.Println(err)
				c.Connected = false
			}
			c.MsgReceiveCh <- string(msg)
		}
	}
}

func (c *WsConn) messageSendHandler() {
	for {
		// 循环监听发送
		select {
		case <-c.Done:
			return
		case t := <-c.MsgSendCh:
			err := c.Conn.WriteMessage(websocket.TextMessage, []byte(t))
			if err != nil {
				log.Println("write:", err)
				return
			}
		// 来自终端的中断运行请求
		case <-c.Interrupt:
			// 向 websocket 发起关闭连接任务
			err := c.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				os.Exit(1)
			}
			select {
			// 等到了结果，直接退出
			case <-c.Done:
			// 如果等待 1s 还是拿不到结果，不再等待，超时退出
			case <-time.After(time.Second):
			}
			os.Exit(1)
			// return
		}
	}
}

func createWsConn() *WsConn {
	// 设置终端中断通道
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "wss", Host: "api.leo.moe", Path: "/v2/ipGeoWs"}
	// log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	wsconn = &WsConn{
		Conn:         c,
		Connected:    true,
		Connecting:   false,
		MsgSendCh:    make(chan string),
		MsgReceiveCh: make(chan string),
	}
	log.Println("WebSocket 连接意外断开，正在尝试重连...")

	if err != nil {
		log.Println("dial:", err)
		return wsconn
	}
	// defer c.Close()
	// 将连接写入WsConn，方便随时可取
	wsconn.Done = make(chan struct{})
	go wsconn.messageReceiveHandler()
	go wsconn.messageSendHandler()
	return wsconn
}

func New() *WsConn {
	return createWsConn()
}

func GetWsConn() *WsConn {
	return wsconn
}
