package es

import (
	"context"
	"es/config"
	"es/helper"
	"github.com/olivere/elastic/v7"
	"log"
)

// NewClient 创建一个Es客户端
func NewClient() (*elastic.Client, error) {
	//初始化日志
	helper.InitLogger()

	client, err := elastic.NewClient(
		elastic.SetURL(config.ElasticSearchHost),
		elastic.SetBasicAuth(config.ElasticSearchUserName, config.ElasticSearchPassword),
		elastic.SetSniff(false),            //禁止自动转换地址
		elastic.SetTraceLog(helper.Logger), //记录日志
	)
	if err != nil {
		log.Println("es服务连接失败：" + err.Error())
		return nil, err
	}
	log.Println("es服务连接成功")

	//获取集群版本信息
	info, code, err := client.Ping(config.ElasticSearchHost).Do(context.Background())
	if err != nil {
		log.Println("获取es集群版本信息失败：" + err.Error())
		return nil, err
	}
	log.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	return client, nil
}
