package constant

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
	CartKey            = "cart:"
	PendingSyncCartKey = "pending:sync:users"
	OrderPreviewKey    = "order:preview"
	//缓存穿透
	RedisEmptyValue    = "_EMPTY_VALUE" //解决go中使用redis.get时，访问不存在的key，err是nil问题
	DefaultReceiveArea = "area:default:"
	GoodsBaseKey       = "goods:base:"
	//大促商品库存预热key
	StockGoodsKey = "stock:goods:"
	MinShortTTL   = 30
	ShortTTL      = 5 * 60
	LongTTL       = 60 * 60
)
