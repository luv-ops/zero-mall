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
	OrderNo   string `json:"orderNo"`
	UserId    string `json:"userId"`
	TimeStamp int64  `json:"timeStamp"`
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
