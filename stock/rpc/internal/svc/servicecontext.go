package svc

import (
	"context"
	"log"
	"zeromall/stock/rpc/internal/config"
	"zeromall/stock/rpc/internal/logic/luaScript"
	"zeromall/stock/rpc/internal/model"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config          config.Config
	StockModel      model.StockModel
	Redis           *redis.Redis
	StockReturnSha  string
	PreHeatStockSha string
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.MustNewConn(c.Mysql)
	rdb := redis.MustNewRedis(c.RedisConf)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sha, err := rdb.ScriptLoadCtx(ctx, luaScript.StockReturn)
	if err != nil {
		log.Fatalf("SCRIPT LOAD stock_return luaScript failed: %v", err)
	}
	sha2, err := rdb.ScriptLoadCtx(ctx, luaScript.StockPreHeat)
	if err != nil {
		log.Fatalf("SCRIPT LOAD stock_preheat luaScript failed: %v", err)
	}
	return &ServiceContext{
		Config:          c,
		StockModel:      model.NewStockModel(conn),
		Redis:           rdb,
		StockReturnSha:  sha,
		PreHeatStockSha: sha2,
	}
}
