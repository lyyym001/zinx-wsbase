package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/sun-fight/zinx-websocket/znet"

	"github.com/gorilla/websocket"
)

var _addr = flag.String("addr", "192.168.0.22:8999", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *_addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			//------重点看这  接收服务器消息
			msgType, ioReader, err := c.NextReader()
			if err != nil {
				fmt.Println("get read reader error ", err)
				return
			}
			//读取客户端的Msg head
			dataPack := znet.NewDataPack()
			headData := make([]byte, dataPack.GetHeadLen())
			if _, err := io.ReadFull(ioReader, headData); err != nil {
				fmt.Println("read msg head error ", err)
				return
			}
			//拆包，得到msgID 和 datalen 放在msg中
			msg, err := dataPack.Unpack(headData)
			if err != nil {
				fmt.Println("unpack error ", err)
				return
			}
			msg.SetMsgType(msgType)

			//根据 dataLen 读取 data，放在msg.Data中
			var data []byte
			if msg.GetDataLen() > 0 {
				data = make([]byte, msg.GetDataLen())
				if _, err := io.ReadFull(ioReader, data); err != nil {
					fmt.Println("read msg data error ", err)
					return
				}
			}
			msg.SetData(data)
			msg.ToString()
		}
	}()

	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()
	var count int
	for {
		select {
		case <-done:
			return
		case _ = <-ticker.C:
			//------重点看这  每秒向服务器发送消息
			count++
			msgID := 1001
			text := "执行登录"
			if count%2 == 0 {
				msgID = 1002
				text = "执行退出登录"
			}
			msgPackage := znet.NewBinaryMsgPackage(uint16(msgID), []byte(text))
			pack, err := znet.NewDataPack().Pack(msgPackage)
			if err != nil {
				log.Println("write:", err)
				return
			}
			err = c.WriteMessage(websocket.BinaryMessage, pack)
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
