package svc

import (
	"context"
	"log"
	"zeromall/cart/rpc/internal/config"
	"zeromall/cart/rpc/internal/logic/luaScript"

	"zeromall/cart/rpc/internal/model"
	"zeromall/common/mq"
	"zeromall/goods/rpc/goodsPb"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config        config.Config
	CartModel     model.CartModel
	GoodsRpc      goodsPb.GoodsClient
	Redis         *redis.Redis
	AddCartSha    string
	BackFillSha   string
	UpdateCartSha string
	Producer      *mq.Producer
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.Mysql.DataSource)
	rdb := redis.MustNewRedis(c.RedisConf)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//加载lua脚本
	sha, err := rdb.ScriptLoadCtx(ctx, luaScript.AddCartScript)
	if err != nil {
		log.Fatalf("SCRIPT LOAD add_cart lua failed: %v", err)
	}
	sha2, err := rdb.ScriptLoadCtx(context.Background(), luaScript.BackFillScript)
	if err != nil {
		log.Fatalf("SCRIPT LOAD back_fill lua failed: %v", err)
	}
	sha3, err := rdb.ScriptLoadCtx(context.Background(), luaScript.UpdateCartScript)
	if err != nil {
		log.Fatalf("SCRIPT LOAD update_cart lua failed: %v", err)
	}
	//注入goodsRpc
	client := zrpc.MustNewClient(c.GoodsRpcConf)

	//挂载全局生产者
	pro, err := mq.NewProducer(c.RocketMqConf.Endpoint, c.RocketMqConf.Producer.TopicSyncFiling)
	if err != nil {
		log.Fatalf("NewProducer failed: %v", err)
	}

	return &ServiceContext{
		Config:        c,
		CartModel:     model.NewCartModel(sqlConn),
		GoodsRpc:      goodsPb.NewGoodsClient(client.Conn()),
		Redis:         rdb,
		AddCartSha:    sha,
		BackFillSha:   sha2,
		UpdateCartSha: sha3,
		Producer:      pro,
	}
}
