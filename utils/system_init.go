package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB    *gorm.DB
	Redis *redis.Client
)

func InitRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConns"),
	})
	redisCtx := Redis.Context()
	pong, err := Redis.Ping(redisCtx).Result()
	if err != nil {
		fmt.Println("Init Redis 。。。。", err)
	} else {
		fmt.Println("config Redis inited", pong)
	}
}

func InitConfig() {

	viper.SetConfigName("app")

	viper.AddConfigPath("config")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("config app inited")

}

func InitMySQL() {

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	DB, _ = gorm.Open(mysql.Open(viper.GetString("mysql.dns")), &gorm.Config{
		Logger: newLogger,
	})
	fmt.Println("config MySQL inited")
}

const (
	PublishKey = "websocket"
)

// Publish 發布消息到Redis
func Publish(ctx context.Context, channel string, msg string) error {
	var err error
	fmt.Println("Published 。。。。", msg)
	err = Redis.Publish(ctx, channel, msg).Err()
	return err
}

// Subscribe 訂閱Redis
func Subscribe(ctx context.Context, channel string) (string, error) {
	sub := Redis.Subscribe(ctx, channel)
	fmt.Println("Subscribe 。。。。", ctx)

	msg, err := sub.ReceiveMessage(ctx)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println("Subscribe 。。。。", msg.Payload)
	return msg.Payload, err
}
