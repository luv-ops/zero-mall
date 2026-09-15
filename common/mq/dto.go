package mq

type CartChangeMsg struct {
	UserId    string `json:"usrId"`
	TimeStamp int64  `json:"timeStamp"`
}
type CartDelMsg struct {
	UserId    string   `json:"usrId"`
	GoodsIds  []string `json:"goodsIds"`
	TimeStamp int64    `json:"timeStamp"`
}
type OrderOffMessage struct {
	OrderNo   int64  `json:"orderNo"`
	UserId    string `json:"userId"`
	TimeStamp int64  `json:"timeStamp"`
}
type OrderTransactionMsg struct {
	OrderNo int64 `json:"orderNo"`
}
type InsertStockMsg struct {
	Stock     int64  `json:"stock"`
	GoodsId   string `json:"goodsId"`
	TimeStamp int64  `json:"timeStamp"`
}
type ReturnStockItem struct {
	GoodsId string `json:"goodsId"`
	Num     int64  `json:"num"`
}
type ReturnStockMsg struct {
	List []*ReturnStockItem `json:"list"`
}
type FrozenStockMsg struct {
	List []*FrozenItem `json:"list"`
}
type FrozenItem struct {
	OrderNo int64  `json:"orderNo"`
	GoodsId string `json:"goodsId"`
	Num     int64  `json:"num"`
}
