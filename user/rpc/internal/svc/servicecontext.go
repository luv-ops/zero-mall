package svc

import (
	"zeromall/common/mq"
	"zeromall/goods/rpc/goodsPb"
	"zeromall/user/rpc/internal/config"
	"zeromall/user/rpc/internal/model"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config          config.Config
	UserModel       model.UserModel
	RecAddressModel model.UserReceiveAddressModel
	Redis           *redis.Redis
	AreaModel       model.AreaModel
	GoodsRpc        goodsPb.GoodsClient
	Producer        *mq.Producer
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.Mysql.DataSource)
	rds := redis.MustNewRedis(c.RedisConf)
	pg := mq.ProducerConfig{
		Endpoint: c.RocketMqConf.Endpoint,
		Topics:   []string{c.RocketMqConf.Topics.TopicInsertBalance},
	}
	pro, err := mq.NewProducer(&pg)
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:          c,
		UserModel:       model.NewUserModel(sqlConn),
		Redis:           rds,
		AreaModel:       model.NewAreaModel(sqlConn),
		RecAddressModel: model.NewUserReceiveAddressModel(sqlConn),
		Producer:        pro,
	}
}
