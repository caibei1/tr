package zhaobi

import (
	"fmt"
	"testing"
)

func TestPostBill(t *testing.T) {
	resp := PostBill("46.082098", "YCC", "BTY", "0.040314", "BUY")
	fmt.Println(resp)
}


