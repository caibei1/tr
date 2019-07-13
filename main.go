package main

import (
	"./exchange/zhaobi"
	log "github.com/sirupsen/logrus"
	"time"
	"tr/exchange/huobi"
	"tr/util"
)

// 初始金额
var money = 100.0
// 未卖出的单
var order =  make(map[int]float64)

func main()  {

	log.Warn("初始金额:",money)
	log.SetLevel(log.WarnLevel)
	// 找币
	zbc := zhaobi.NewZBClient()
	zbc.Init()

	// 火币
	hbc := huobi.NewHBClient()
	hbc.Init()

	// 等待初始化完成
	time.Sleep(time.Second*5)
	i := 0
	// 监控找币与火币价格差
	go func() {
		for   {
			for k,v := range huobi.Symbol{
				buy := hbc.GetLastBuyPrice(k)
				sell := zbc.GetLastSellPrice(k)
				b := util.IsOnePercent(buy,sell)
				//log.Warn(v,buy,sell)
				if b {
					// 买入
					log.Warn("买入")
					if money >0 {
						if money >= 10.0 {
							// 每次10块
							GuaDan(sell,10,0,1,i,v)

							// 监控卖出
							go WaitSell(zbc,sell,i,k,v)

							i ++

						}
					}
					time.Sleep(time.Minute)
				}
			}
			time.Sleep(time.Second)
		}
	}()


	// 监控ycc 低买高卖
	//go zhaobi.StartYCC(zbc,exchange.YCC)


	select {

	}
}

// t 1买 2 卖  buyPrice挂卖单的时候的买入价格
func GuaDan(price, count, buyPrice float64, t int, mark int,symbol string)  {
	if t == 1 {
		log.Warnf("====买入%s，买入价格：%f, 买入金额：%f, mark: %d", symbol, price, count, mark)
		price = price - 10
		order[mark] = price
	}else {
		log.Warnf("====卖出%s，买出价格：%f, 卖出收益百分比：%f, mark: %d", symbol, (price-buyPrice)/price , count, mark)
		price = price + 10*(1.0+(price-buyPrice)/price - 0.001)
		delete(order,mark)
	}
	log.Warn("余额：%f", price)
	for k,v := range order{
		log.Warn("正在进行中的单子： %d, %f",k, v)
	}
}

func WaitSell(z *zhaobi.ZBClient,buyP float64, mark int,symbol int, symbolStr string)  {
	for {
		buy := z.GetLastBuyPrice(symbol)
		if (buy - buyP)/buyP >= 0.01 {
			log.Warn("====监控到价格大约0.01，价格差为 %f, 卖出", (buy - buyP)/buyP)
			GuaDan(buy,10,buyP,2,mark,symbolStr)
		}
	}

}