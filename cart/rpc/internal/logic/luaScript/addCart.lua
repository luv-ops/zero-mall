local cartKey=KEYS[1]
local GoodsField=ARGV[1]
local AddNum=tonumber(ARGV[2])
local Selected=tonumber(1)
local oldStr=redis.call("HGET",cartKey,GoodsField)

if oldStr then
    --商品存在，则累计num，返回0，不需要调用goodsRpc
    local item= cjson.decode(oldStr)
    item.num=item.num+AddNum
    redis.call("HSET",cartKey,GoodsField,cjson.encode(item))
    return 0
else
    --商品不存在，则返回1，需要调用goodsRpc，回填
    local total=redis.call("HLEN",cartKey)
    if total>=30 then
        return redis.error_reply("购物车已满，最多30件")
    else
        local placeholder=cjson.encode({
            num=AddNum,
            selected=Selected,
            name="",
            cover="",
            price=""
        })
        redis.call("HSET",cartKey,GoodsField,placeholder)
        return 1
    end

end