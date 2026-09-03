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
	AreaNotFound        = "地区不存在"
	PhoneIllegal        = "手机号非法"
)

// 后端打印信息
const (
	RedisFailed  = "%s Redis failed"
	MysqlFailed  = "%s Mysql %s failed"
	WhereFailed  = "%s %s error"
	UnmarshalErr = "unmarshal err in %s: err:%v"
	MarshalErr   = "marshal err in %s: err:%v"
	RpcError     = "rpc err in %s:%v"
)

// redis
const (
	AreaKey            = "area:pid:%d:level:%d"
	UserInfoKey        = "user:info:"
	GoodsInfoKey       = "goods:info:"
	ShortTTL           = 5 * 60
	LongTTL            = 60 * 60
	CartKey            = "cart:"
	PendingSyncCartKey = "pending:sync:users"
	//缓存穿透
	RedisEmptyValue    = "_EMPTY_VALUE" //解决go中使用redis.get时，访问不存在的key，err是nil问题
	DefaultReceiveArea = "area:default:"
)
