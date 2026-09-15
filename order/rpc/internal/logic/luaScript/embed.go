package luaScript

import _ "embed"

//go:embed stockFrozen.lua
var StockFrozen string

//go:embed stockReturn.lua
var StockReturn string
