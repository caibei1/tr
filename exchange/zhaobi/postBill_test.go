package zhaobi

import (
	"fmt"
	"testing"
)

func TestPostBill(t *testing.T) {
	resp := PostBill("2", "YCC", "BTY", "0.0441", "BUY")
	fmt.Println(resp)
}
