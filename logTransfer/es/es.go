package es

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"strings"
	"time"
)

// 初始化ES，准备接收kafka那边发来的数据

type LogData struct{
	Topic string `json:"topic"`
	Data string `json:"data"`
}

var (
	client *elastic.Client
	ch chan  *LogData
)

func Init(address string,chanSize int,nums int)(err error){
	if !strings.HasPrefix(address,"http://"){
		address = "http://" + address
	}
	client, err = elastic.NewClient(elastic.SetURL(address))
	if err != nil {
		// Handle error
		return
	}
	ch = make(chan *LogData,chanSize)
	for i:=0;i<nums;i++{
		go SendToES()
	}

	return
}

// 发送数据到ES
func SendToESChan(msg *LogData) {
	// ?
	ch <- msg
}

func SendToES(){
	// 链式操作
	for{
		select {
		case msg := <- ch:
			put1, err := client.Index().Index(msg.Topic).Type("xxx").BodyJson(msg).Do(context.Background())
			if err != nil {
				// Handle error
				fmt.Println(err)
				continue
			}
			fmt.Printf("Indexed %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
		default:
			time.Sleep(time.Second)
		}
	}

}