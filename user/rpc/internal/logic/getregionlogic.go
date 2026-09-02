package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"zeromall/common/constant"
	"zeromall/user/rpc/internal/svc"
	"zeromall/user/rpc/userpb"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetRegionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRegionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRegionLogic {
	return &GetRegionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetRegionLogic) GetRegion(in *userpb.GetRegionReq) (*userpb.GetRegionResp, error) {
	// todo: add your logic here and delete this line
	if in.Pid == nil || in.Level == nil {
		return nil, status.Error(codes.InvalidArgument, "参数不全")
	}
	if *in.Pid < 0 || *in.Pid > 5000 || *in.Level < 1 || *in.Level > 3 {
		return nil, status.Error(codes.InvalidArgument, "参数错误")
	}
	var list []*userpb.RegionItem
	//查询缓存
	key := fmt.Sprintf(constant.AreaKey, *in.Pid, *in.Level)
	cacheStr, err := l.svcCtx.Redis.GetCtx(l.ctx, key)
	if err == nil && cacheStr != "" {
		//缓存命中
		err = json.Unmarshal([]byte(cacheStr), &list)
		if err != nil {
			l.Logger.Errorf("反序列化失败 err: %s", err.Error())
		} else {
			//反序列化成功
			return &userpb.GetRegionResp{
				List: list,
			}, nil
		}
		//反序列化失败，也代表缓存损坏
	}
	//缓存未命中
	//查询数据库
	area, err := l.svcCtx.AreaModel.SelectFields(l.ctx, *in.Level, *in.Pid)
	if err != nil {
		l.Logger.Errorf(constant.MysqlFailed, "region", "select", err.Error())
		return nil, status.Error(codes.Internal, constant.MiddlewareError)
	}
	//有地区数据，才缓存
	if len(area) > 0 {
		for _, item := range area {
			list = append(list, &userpb.RegionItem{
				Id:    item.Id,
				PId:   item.Pid,
				Name:  item.Name.String,
				Level: item.Level,
			})
		}
		//设置缓存
		data, err := json.Marshal(list)
		if err != nil {
			l.Logger.Errorf("序列化失败 err: %s\", err.Error()")
		} else {
			//序列化成功才写入缓存
			err = l.svcCtx.Redis.SetexCtx(l.ctx, key, string(data), 3600*24*7)
			if err != nil {
				l.Logger.Errorf("缓存写入失败 err: %s", err.Error())
			}
		}
	}

	return &userpb.GetRegionResp{
		List: list,
	}, nil
}
