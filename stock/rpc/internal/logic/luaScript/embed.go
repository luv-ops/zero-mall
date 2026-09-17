package luaScript

import _ "embed"

//go:embed stockReturn.lua
var StockReturn string

//go:embed preHeatStock.lua
var StockPreHeat string

//go:embed stockDeduct.lua
var StockDeduct string
