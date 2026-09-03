### goctl根据.api文件生成目录
```
goctl api go -api api/user.api -dir .
goctl api go -api api/goods.api -dir .
```

### goctl根据.proto文件生成目录
```
goctl rpc protoc proto/user.proto --go_out=. --go-grpc_out=. --zrpc_out=.
```
### goctl在线连接mysql自动创建表
```
goctl model mysql datasource --url "root:3110940369w@tcp(127.0.0.1:3306)/zero_user" -t "user" --dir "internal/model
```

### 启动rpc和api 相同命令
```
 go run user.go -f etc/user.yaml
 go run goods.go -f etc/goods.yaml
 go run order.go -f etc/order.yaml
 go run cart.go -f etc/cart.yaml
```

### 某些字段需要设置可选，比如请求体某些字段需要进行可选
``` api
type ChangeInfoReq {
	UserName *string `json:"username,optional,omitempty"`
	Avatar   *string `json:"avatar,optional,omitempty"`
	Age      *uint32 `json:"age,optional,omitempty"`
	Sex      *uint32 `json:"sex,optional,omitempty"`
	Region   *string `json:"region,optional,omitempty"`
}
```
``` protobuf
message ChangeInfoReq{
  string userId=1;
  optional string username=2;
  optional string avatar=3;
  optional string region=4;
  optional uint32 age=5;
  optional uint32 sex=6;
}
```

### 数据库update不能仅通过err!=nil判断是否更新成功，还应该拿到影响行数是否>0

### api层的get请求参数应该设置 form而不是json
```api
type GetRegionReq {
	Level int64  `form:"level"`
	PId   *int64 `form:"pId,optional,omitempty"`
}
```
### api层或proto需要设置切片类型,都需要进行一层包装
```api
type RegionItem {
	Id    int64  `json:"id"`
	PId   int64  `json:"pId"`
	Name  string `json:"name"`
	Level int64  `json:"level"`
}

type GetRegionResp {
	List []*RegionItem `json:"list"`
}
```

```protobuf3
message RegionItem{
  int64 id=1;
  int64 pId=2;
  string name=3;
  int64 level=4;
}
message GetRegionResp{
  repeated RegionItem list=1;
}
```
### sqlx问题
```go
err := m.conn.QueryRowPartialCtx(ctx, user, sqlStr, userId)
err := m.conn.QueryRowCtx(ctx, user, sqlStr, userId)

区别:如果使用1，允许批量查询
	如果使用2，那你的message体中的字段必须和查询字段一样，全量查询

```
### 修改文件根目录名字，会出现 新名字[旧名字]
    在资源目录找到文件夹，删除里面的.idea就行

### 把一个schema多个表快速迁移到不同scheama
    rename table 原schema.表 to 新scheama.表

### rocketMq in go-zero
``` go
    //创建生产者时，必须设置withTopics防止死锁，并且这个bug不会输出任何信息，导致你的主程序启动不了
    producer, err := rmq_client.NewProducer(&rmq_client.Config{
		Endpoint: Endpoint,
		Credentials: &credentials.SessionCredentials{
			AccessKey:    AccessKey,
			AccessSecret: SecretKey,
		},
	},
		rmq_client.WithTopics(Topic),
	)
```
### lua脚本，注意数字类型一定要使用tonumber函数，不然可能导致反序列失败
    比如:json unmarshal error  json: cannot unmarshal string into Go struct field CartItem.selected of type int64

### golang使用redis时，如果getCtx ,key不存在会返回空串和nil，key存在但是值为空也是返回空串,nil