package main

import (
	"./exchange/zhaobi"
	log "github.com/sirupsen/logrus"
	"time"
	"tr/exchange/huobi"
	"tr/util"
	"tr/exchange"
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
	time.Sleep(time.Second*5)

	// 监控找币与火币价格差
	go func() {
		for   {
			for k,_ := range huobi.Symbol{
				buy := hbc.GetLastBuyPrice(k)
				sell := zbc.GetLastSellPrice(k)
				b := util.IsOnePercent(buy,sell)
				//log.Warn(v,buy,sell)
				if b {
					time.Sleep(time.Minute)
				}
			}
			time.Sleep(time.Second)
		}
	}()


	// 监控ycc 低买高卖
	go zhaobi.StartYCC(zbc,exchange.YCC)


	select {

	}
}


