
### 介绍
    这是一个分布式商城项目，每一个服务中api层负责鉴权和参数检验，rpc层主要用来写业务和操作mysql，redis，消息队列
    对于已开发好的服务，已经提供了.sql文件，docker-compose目前用来本地部署rocketmq
    gateway将在未来对于不同端口的api服务，统一端口让前端只需要对接网关端口就行
### 技术栈
    go-zero+grpc+mysql+redis+rocketmq
### 详细亮点
    1.购物车功能: 1.以redis存储为主，mysql只归档，当用户修改购物车时，会调用lua脚本后，会生产一条消息，携带用户id和时间戳，消费者会拿到用户id并且将其redis sadd
    到待同步集合，ticker会每10s从redis，spop 50条待同步用户id，然后同步归档到mysql中。此设计天然能应对用户频繁修改购物车，因为10s内用户的频繁修改都会sadd到redis，但是redis 集合数据结构
    会自动去重。2.同步mysql时采用增量修改，就是先取出redis数据，然后查询mysql，进行对比，哪些需要修改，哪些需要插入，哪些需要删除。
    2.缓存穿透问题: golang使用redis时，访问不存在的key,得到的值是空字符串，并且error是nil，所以需要自己建立一个业务缓存空值
    
