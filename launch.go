package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

//连接客户端
var clients = make(map[*websocket.Conn]bool)

//广播通道
var broadcasts = make(chan Message)

type Message struct {
	Id      uuid.UUID `json:"uuid"`
	Name    string    `json:"name"`
	Message string    `json:"message"`
}

//配置连接升级设置
var upgrader = websocket.Upgrader{
	HandshakeTimeout: 0,
	ReadBufferSize:   0,
	WriteBufferSize:  0,
	WriteBufferPool:  nil,
	Subprotocols:     []string{},
	Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
	},
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	EnableCompression: false,
}

func main() {
	//建立简单的文件系统
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	//websocket 路由配置
	http.HandleFunc("/ws", handleConnections)

	http.HandleFunc("/uuid", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(uuid.NewV1().String()))
	})

	//启动信息监听
	go handleMessage()

	//启动服务
	log.Printf("Http Service Start On localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	//将初始化get请求升级为weboscket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Handleing Connect Failure: %v\n", err)
	}
	//函数结束时关闭连接
	defer ws.Close()

	//注册新的连接
	clients[ws] = true

	//循环工作
	for {
		var msg Message
		//读取json信息，并转换为message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Close Websocket: Error: %v\n", err)
			delete(clients, ws)
			break
		}
		//将信息发送到广播通道
		broadcasts <- msg
	}
}
func handleMessage() {
	for {
		//获取最新的信息，并且广播到说有连接上的客户端
		msg := <-broadcasts
		//发送信息
		for client := range clients {
			//将信息写入每个客户端
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Failure Write Message Into Client : Error: %v\n", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
