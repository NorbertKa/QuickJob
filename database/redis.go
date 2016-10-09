package database

import (
	"errors"
	"fmt"
	"strconv"
	"sync"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/mievstac/QuickJob/config"
	redis "gopkg.in/redis.v4"
)

var (
	ErrConnectRedis = errors.New("Can't connect to redis, Host:%v, Port:%d, DB:%d")
)

type RedisConn struct {
	sync.Mutex
	Host      string
	Port      int
	Password  string
	DB        int
	connected bool
	Client    *redis.Client
}

func NewRedisConn(conf *config.Config) *RedisConn {
	return &RedisConn{
		Host:      conf.Redis.Host,
		Port:      conf.Redis.Port,
		Password:  conf.Redis.Password,
		DB:        conf.Redis.DB,
		connected: false,
	}
}

func (conf *RedisConn) flipConnection() {
	conf.Lock()
	conf.connected = !conf.connected
	conf.Unlock()
}

func (conf *RedisConn) Connect() error {
	conf.Lock()
	client := redis.NewClient(&redis.Options{
		Addr:     conf.Host + ":" + strconv.Itoa(conf.Port),
		Password: conf.Password,
		DB:       conf.DB,
	})
	pong, err := client.Ping().Result()
	if pong != "PONG" || err != nil || client == nil {
		conf.Unlock()
		return fmt.Errorf(ErrConnectRedis.Error(), conf.Host, conf.Port, conf.DB)
	}
	conf.connected = true
	conf.Client = client
	conf.Unlock()
	return nil
}

func (conf *RedisConn) CloseRedisConnection() error {
	if conf.connected == false {
		return errors.New("Closing unconnected Database (connected = false)")
	}
	return conf.Client.Close()
}

func (conf RedisConn) CheckConnection() (bool, error) {
	if conf.connected == false {
		return false, errors.New("No connection (connected == false)")
	}
	pong, err := conf.Client.Ping().Result()
	if pong != "PONG" || err != nil {
		conf.flipConnection()
		return false, err
	}
	return true, nil
}

func (conf RedisConn) forceCheckConnection() error {
	client := redis.NewClient(&redis.Options{
		Addr:     conf.Host + ":" + strconv.Itoa(conf.Port),
		Password: conf.Password,
		DB:       conf.DB,
	})
	pong, err := client.Ping().Result()
	if pong != "PONG" || err != nil || client == nil {
		return fmt.Errorf(ErrConnectRedis.Error(), conf.Host, conf.Port, conf.DB)
	}
	return err
}

func (conf RedisConn) Connected() bool {
	return conf.connected
}

func (conf RedisConn) Validate() (bool, error) {
	conf.Lock()
	connectionBool := true
	var result error
	if conf.Port <= 0 || conf.Port > 65535 {
		connectionBool = false
		result = multierror.Append(result, errors.New("Redis Port Out of Range"))
	}
	connectionErr := conf.forceCheckConnection()
	if connectionErr != nil {
		connectionBool = false
		result = multierror.Append(result, connectionErr)
	}
	if connectionBool == true {
		return true, nil
	}
	return false, result
}
