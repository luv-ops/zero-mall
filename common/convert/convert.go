package convert

import (
	"errors"
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

func YuanStrToCents(price string) (int64, error) {
	// 1. 从字符串解析为 Decimal 类型
	d, err := decimal.NewFromString(price)
	if err != nil {
		return 0, errors.New("无效的金额格式")
	}

	// 2. 乘以 100，将“元”转换为“分”
	centsDecimal := d.Mul(decimal.NewFromInt(100))

	// 3. 四舍五入到最接近的整数分，并转为 int64
	// Round(0) 表示保留 0 位小数，即四舍五入到整数
	return centsDecimal.Round(0).IntPart(), nil
}
