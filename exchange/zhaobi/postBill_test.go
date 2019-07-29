package zhaobi

import (
	"fmt"
	"testing"
)

func TestPostBill(t *testing.T) {
	resp := PostBill("0.0001", "YCC", "usdt", "0.099999", "SELL")
	fmt.Println(resp)
}
