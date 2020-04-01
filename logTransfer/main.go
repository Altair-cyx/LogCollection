package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"log_transfer/config"
	"log_transfer/es"
	"log_transfer/kafka"
)

// log transfer
// 将日志数据从kafka取出来发往ES



func main(){
	// 0.加载配置文件
	var cfg = new(config.AppConf)
	err := ini.MapTo(cfg,"./config/config.ini")
	if err != nil{
		fmt.Println("init config failed,err:",err)
		return
	}
	fmt.Printf("cfg:%v\n",cfg)
	// 1.初始化ES
	// 1.1初始化一个ES连接的client
	// 1.2对外提供一个往ES写数据的函数
	err = es.Init(cfg.ESConf.Address,cfg.ESConf.ChanSize,cfg.ESConf.Nums)
	if err != nil{
		fmt.Printf("init ES client failed,err:%v",err)
		return
	}
	fmt.Println("init es success.")
	// 2.初始化kafka
	// 2.1 连接kafka，创建分区的消费者
	// 2.2 每个分区的消费者分别取出数据，通过SendToES()将数据发往ES
	err = kafka.Init([]string{cfg.KafkaConf.Address},cfg.KafkaConf.Topic)
	if err != nil{
		fmt.Printf("init kafka consumer failed,err:%v",err)
		return
	}
	select {

	}
}