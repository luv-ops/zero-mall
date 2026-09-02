package constant

// grpc返回错误信息
const (
	MiddlewareError     = "服务繁忙，请稍后重试"
	MsgSmsTooFrequently = "短信发送过于频繁，请稍后再试"
	UserNotFound        = "用户不存在"
	UserExist           = "用户已存在"
	PasswordError       = "密码错误"
	CaptchaError        = "验证码错误"
	PermissionSellError = "未获得售卖权"
	GoodsArgError       = "商品参数错误"
	GoodsNotFound       = "商品不存在"
)

// 后端打印信息
const (
	RedisFailed = "%s Redis failed"
	MysqlFailed = "%s Mysql %s failed"
	WhereFailed = "%s %s error"
)

// redis key
const (
	AreaKey            = "area:pid:%d:level:%d"
	CartKey            = "cart:"
	PendingSyncCartKey = "pending:sync:users"
)
