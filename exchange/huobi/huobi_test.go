package huobi

import (
	"fmt"
	"testing"
	"time"
)

func TestGetHBInfo(t *testing.T) {

	//info := GetHBInfo(3)
	//fmt.Println(info)
	c := NewHBClient()
	for true {
		time.Sleep(time.Second*2)
		fmt.Println(c.GetLastBuyPrice(3))
		fmt.Println(c.GetLastSellPrice(3))
	}
}
