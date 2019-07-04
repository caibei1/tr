package main

import (
	"./exchange/zhaobi"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
	"tr/exchange/huobi"
	"tr/util"
)

func main()  {
	log.SetLevel(log.WarnLevel)
	// 找币
	zbc := zhaobi.NewZBClient()
	zbc.Init()

	// 火币
	hbc := huobi.NewHBClient()
	hbc.Init()

	// 等待初始化完成
	time.Sleep(time.Second*3)

	// 监控找币与火币价格差
	go func() {
		for   {
			for k,v := range huobi.Symbol{
				buy := hbc.GetLastBuyPrice(k)
				sell := zbc.GetLastSellPrice(k)
				b := util.IsOnePercent(buy,sell)
				log.Warn(v,buy,sell)
				if b {
					time.Sleep(time.Minute)
					fmt.Println("=====")
				}
			}
			time.Sleep(time.Second)
		}
	}()

	select {

	}
}


