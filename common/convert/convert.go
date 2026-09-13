package convert

import (
	"fmt"

	"github.com/shopspring/decimal"
)

func CentsToYuanStr(priceCent int64) string {
	sign := ""
	if priceCent < 0 {
		sign = "-"
		priceCent = -priceCent // 转为绝对值
	}
	yuan := priceCent / 100
	cent := priceCent % 100
	return fmt.Sprintf("%s%d.%02d", sign, yuan, cent)
}

func YuanStrToCents(price string) int64 {
	priceStr, _ := decimal.NewFromString(price)
	return priceStr.Round(2).Mul(decimal.NewFromInt(100)).IntPart()
}
