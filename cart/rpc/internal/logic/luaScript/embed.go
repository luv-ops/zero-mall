package luaScript

import _ "embed"

//go:embed addCart.lua
var AddCartScript string

//go:embed backFill.lua
var BackFillScript string

//go:embed updateCart.lua
var UpdateCartScript string
