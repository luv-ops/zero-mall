# 分布式商城项目

## 介绍
本项目基于微服务架构开发，服务分层职责清晰：
- **API层**：对外HTTP/gRPC入口，负责鉴权、参数校验、请求转发，不处理核心业务逻辑
- **RPC层**：实现核心业务逻辑，完成 MySQL、Redis、消息队列等中间件操作
- 项目已提供业务建表 `.sql` 脚本，可快速初始化数据库；
- `docker‑compose` 用于本地快速部署 RocketMQ 消息中间件和etcd；
- Gateway 网关：聚合多组 API 服务，对外暴露统一端口，前端仅对接网关即可，屏蔽后端多服务多端口细节。

## 部署
- 拉取项目后，进入项目中的docker文件夹，然后执行 docker-compose up -d,快速本地部署rocketmq和etcd
- 进入各个服务的rpc/etc文件夹，把.yaml.template修改一下mysql的password,然后去除.template后缀
- 然后每个服务提供了.sql文件，通过.sql快速创建数据库和各个表
- 由于rocketmq不像kafka一样提供了api可以在程序中创建topic，所以docket启动后，去127.0.0.1:8888控制面板去预先创建topic，不然启动不了程序,topic名称和类型已在下面给出
- 然后进入各个服务的api和rpc文件夹，去执行他们的main函数，一个一个启动(没办法，这是微服务项目，也有想过通过编写.bat去批量启动，但是会创建10几个黑窗口)
- 最后启动gateway文件夹中的main包，部署一下网关，前端可以直接统一和网关的8000端口对接了

## topic
| name                  | messageType |
|-----------------------|-------------|
| balance_insert_topic  | normal      |
| cart_change_topic     | normal      |
| cart_del_topic        | normal      |
| stock_insert_topic    | normal      |
| order_delay_off_topic | delay       |
| stock_return_topic    | normal      |
| stock_tx_frozen_topic | transaction |
| pay_tx_success_topic  | transaction |
| stock_deduct_topic    | normal      |


## 技术栈
`go‑zero` + `gRPC` + `MySQL` + `Redis` + `RocketMQ`

## 项目亮点
### 1. 购物车高性能异步同步方案
> Redis为主存储，MySQL仅做归档落地，应对用户高频次增删改购物车操作
1. 用户修改购物车，执行 **Lua脚本** 操作Redis购物车数据，同时发送RocketMQ消息，携带用户ID与时间戳；
2. 消费端接收消息，将用户ID `SADD` 写入Redis待同步集合；集合天然去重，**同一个用户10s内多次修改只会留存一条用户ID**，避免重复同步；
3. 后台定时Ticker每10秒执行一次同步任务：`SPOP` 批量取出最多50个待同步用户ID；
4. 增量同步归档 MySQL：读取该用户Redis购物车完整数据，与MySQL归档数据做对比，区分新增、修改、删除项，仅做差异更新，减少数据库IO压力。

> 优势：高频修改全部落在Redis，削峰异步落地MySQL，极大降低数据库压力。

### 2.分布式事务最终一致性
> 使用rocketmq的事务消息，先往broker发送半消息，经过tx.commit后，才转为可投递消息到达消费者
1. 创建订单和创建支付记录都用到了事务消息
2. 支付订单时，先查询订单状态和过期时间做提前拦截和防止用户越权
3. 调用payRpc中，开启事务消息，并发送半消息，然后调用balanceRpc扣减余额，再插入支付记录表
4. 如果扣减余额成功，插入支付记录失败，调用balanceRpc的补偿任务退款，然后tx.rollback，半消息不再转为可投递消息

### 3.延迟消息处理
> 不相信延迟消息的延迟时间，因为延迟消息可能在你设置的时间提前投递
1. 先查询订单状态和过期时间做提前拦截，防止延迟消息提前消费
### 4.rocketmq在go-zero中的架构
> 公共包的mq提供消费者和生产者的构造函数，还有各自的config，并提供consumer_group实现了serviceGroup的方法用来统一管理消费者
1. 普通生产者可以直接在svc依赖中进行构造注入，事务消息生产者为了防止循环依赖在main包中初始化，再注入svc
2. 每个服务的消费者有各自的消费逻辑，只需通过构造消费者管理者，就可以在main包中和api服务一样统一start和stop

### 5. 缓存穿透处理方案
> go‑zero Redis封装特性：当Key不存在时，`GetCtx` 返回 `cacheStr=""`、`err=nil`，**无法区分「Key不存在」和「Key存储空字符串」两种场景，原生空字符串占位方案失效**。
1. 不使用空字符串作为空缓存占位，自定义业务特殊标记字符串（魔串）作为空值标记；
2. 查询缓存：读到特殊标记代表数据库确认无该业务数据，直接返回业务 `NotFound`；
3. 返回空字符串代表Key从未创建，放行穿透查询数据库；
4. DB查询不到数据时，写入带较短TTL的特殊标记占位，阻挡后续无效请求打数据库；
5.配合接口层参数合法性校验，双层防护缓存穿透。

### 6.错误传递方式
> 封装xerr和map，对于业务错误，应该清晰的返回给前端，而不是用rpc err 一层一层包装，让前端迷惑。
1. api层调用rpc时的错误应该使用xerr.FromRpcErr，拿到业务错误码，如果不是xerr类型，则直接返回服务错误，不会包含中间件的敏感信息
2. rpc层调用其他rpc层也应该使用xerr.FromRpcErr透传最底层rpc返回的错误。
3. 如果每一个rpc错误无脑返回status.Error,最终返回给前端的响应会类似于这样:rpc error: code = Internal desc = rpc error: code = Internal desc = rpc error: code = Aborted desc =


