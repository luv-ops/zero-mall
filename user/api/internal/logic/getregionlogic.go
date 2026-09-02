// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"zeromall/user/api/internal/svc"
	"zeromall/user/api/internal/types"
	"zeromall/user/rpc/userpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRegionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRegionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRegionLogic {
	return &GetRegionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRegionLogic) GetRegion(req *types.GetRegionReq) (resp *types.GetRegionResp, err error) {
	// todo: add your logic here and delete this line
	res, err := l.svcCtx.UserRpc.GetRegion(l.ctx, &userpb.GetRegionReq{
		Level: req.Level,
		Pid:   req.PId,
	})
	if err != nil {
		return nil, err
	}
	var list []*types.RegionItem
	for _, item := range res.List {
		list = append(list, &types.RegionItem{
			Id:    item.Id,
			PId:   item.PId,
			Level: item.Level,
			Name:  item.Name,
		})
	}
	return &types.GetRegionResp{
		List: list,
	}, nil
}
