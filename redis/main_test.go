package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	//TODO...可以在这里进行初始化部分参数
	retCode := m.Run() //执行测试
	os.Exit(retCode)
}

func TestSet(t *testing.T) {
	ctx := context.TODO()
	err := RedisCli.Set(ctx, "a", true, time.Minute*5)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("success...")
}

func TestGet(t *testing.T) {
	ctx := context.TODO()
	result := RedisCli.Get(ctx, "a")
	if result.Err() != nil {
		fmt.Printf("redis get error: %s\n", result.Err().Error())
		return
	}
	fmt.Printf("%s\n", result.Val())
	fmt.Println(result.Int64())
	fmt.Println(result.Bool())
}

func TestGetLock(t *testing.T) {
	ctx := context.TODO()
	err := RedisCli.GetLock(ctx, "b", time.Hour)
	fmt.Println("err.....")
	fmt.Println(err)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}
	fmt.Println("ok....")
}

func TestReleaseLock(t *testing.T) {
	ctx := context.TODO()
	err := RedisCli.ReleaseLock(ctx, "b")
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}
	fmt.Println("ok....")
}

//===================== 发布订阅模式 =====================

const channelName1 = "Li:channel_001"
const channelName2 = "Li:channel_002"

// 发布订阅模式_消费者消费启动
func pubSubStartConsume(chanName string) {
	fmt.Println("消费者订阅消息启动...")

	//客户端对通道进行订阅
	ctx := context.TODO()
	pubSub := RedisCli.Client.Subscribe(ctx, chanName)

	//等待确认订阅成功
	_, err := pubSub.Receive(ctx)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("消费者订阅消息订阅成功，正在接收消息...")

	for msg := range pubSub.Channel() {
		fmt.Println(msg.String())
	}
}

// 发布订阅模式_消费者1消费启动
func TestRunPubSubStartConsume1(t *testing.T) {
	pubSubStartConsume(channelName1)
}

// 发布订阅模式_消费者2消费启动
func TestRunPubSubStartConsume2(t *testing.T) {
	pubSubStartConsume(channelName2)
}

// 发布订阅模式_消息发送
func TestPubSubPublishMessage(t *testing.T) {
	ctx := context.TODO()
	//发布消息
	result := RedisCli.Client.Publish(ctx, channelName1, channelName1+" ---- hello world ---- "+time.Now().Format(time.DateTime))
	fmt.Printf("%+v \n", result)
	fmt.Println(result.Result())
	fmt.Println(result.Val())
	fmt.Println(result.Name())
	fmt.Println(result.Args())
	fmt.Println(result.FullName())
	fmt.Println(result.Err())
}

// ===================== Redis Stream 模式 =====================
const StreamName1 = "stream_001"

// 添加队列消息
func TestStreamXAddMessage(t *testing.T) {
	ctx := context.TODO()
	messages := make([]*redis.XAddArgs, 0, 10)
	messages = append(messages,
		&redis.XAddArgs{
			Stream: StreamName1,
			ID:     "*", //消息 id，我们使用 * 表示由 redis 生成，可以自定义，但是要自己保证递增性。
			Values: map[string]interface{}{"name": "LiXianPei", "age": 18},
		},
		&redis.XAddArgs{
			Stream: StreamName1,
			ID:     "*", //消息 id，我们使用 * 表示由 redis 生成，可以自定义，但是要自己保证递增性。
			Values: map[string]interface{}{"name": "WangXiEr", "age": 20},
		},
	)
	for _, item := range messages {
		result, err := RedisCli.Client.XAdd(ctx, item).Result()
		fmt.Println(result, err)
	}
}

// 消费队列消息
func TestStreamConsumeMessage(t *testing.T) {
	for {
		ctx := context.TODO()
		messages, err := RedisCli.Client.XRead(ctx, &redis.XReadArgs{
			Streams: []string{StreamName1, "0"}, // 流名称和起始ID
			Block:   1,                          // 阻塞模式，0 表示无限等待，即使流数据中为空也会继续等待，否则（设置为1）当流中为空时直接报错：Error reading from stream: read tcp 192.168.220.1:49629->192.168.220.128:6379: i/o timeout
			Count:   1,                          // 每次读取的消息数量
		}).Result()

		if err != nil {
			fmt.Println("Error reading from stream:", err)
			return
		}

		re := RedisCli.Client.XLen(ctx, StreamName1)
		fmt.Println("StreamName1--------长度：", re.Val())

		for _, stream := range messages {
			fmt.Printf("Received message: %v\n", stream.Stream)
			if len(stream.Messages) > 0 {
				for _, message := range stream.Messages {
					fmt.Printf("%+v\n", message)
					res := RedisCli.Client.XAck(ctx, StreamName1, "", message.ID)
					fmt.Printf("XAck-Result:%+v\n", res)
					//消息确认后可以从流中删除
					if res.Err() == nil {
						//消息如果不删除会继续存在缓存中，只有删除后才不会重复处理，否则for循环将会一直读取到
						RedisCli.Client.XDel(ctx, StreamName1, message.ID)
					}
					//尝试从
				}
			}
		}
		time.Sleep(time.Second)
	}
}
