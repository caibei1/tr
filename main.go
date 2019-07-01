package main

import (
	"./exchange/zhaobi"
	log "github.com/sirupsen/logrus"
	"time"
	"fmt"
)

func main()  {
	log.SetLevel(log.InfoLevel)
	test()
}


func test()  {
	zhaobi.Init()
	zhaobi.GetZBIndexInfo()
	zhaobi.GetZBMarketInfo(10,1)
	zbc := zhaobi.NewZBClient()
	for true {
		time.Sleep(time.Second)
		fmt.Println(zbc.GetLastBuyPrice(1))
		fmt.Println(zbc.GetLastSellPrice(1))
		fmt.Println(zbc.GetLastSuccessRMBPrice(1))
		fmt.Println(zbc.GetOpen(1))
	}

}