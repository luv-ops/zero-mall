package intercepter

import (
	"context"
	"errors"
	"zeromall/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrorInterceptor(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	resp, err := handler(ctx, req)
	if err == nil {
		return resp, nil
	}

	// 业务错误
	var e *xerr.CodeError
	if errors.As(err, &e) {
		return nil, status.Error(codes.Code(e.Code), e.Msg)
	}

	// 系统/中间件错误
	logx.WithContext(ctx).Errorf("rpc internal error, method: %s, err: %+v", info.FullMethod, err)
	return nil, status.Error(codes.Internal, xerr.NewCodeError(xerr.ServerErr).Error())
}
