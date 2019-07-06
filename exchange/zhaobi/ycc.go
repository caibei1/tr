package zhaobi

import (
	"tr/exchange"
	log "github.com/sirupsen/logrus"
	"time"
)

//针对ycc的策略  也适用于差价大的币
//主要逻辑，当卖单和买单差价过大时，出现高价买单，则卖出，同时挂单低价买入

// 是否监听的标志位
var flagYCC = 0

func StartYCC(c exchange.ExClient,sy int)  {

	flagYCC = 1
	// 当买卖差价大于百分之8
	for {

		sell := c.GetLastSellPrice(sy)
		buy := c.GetLastBuyPrice(sy)

		// 控制上限，防止程序出bug时候跑
		if (sell - buy)/buy > 0.08 && (sell - buy)/buy < 0.2{
			log.Warnf(">>>检测到%s币价差距过大，开始监听。",Symbol[sy])
			log.Warnf("买：%f, 卖: %f, 当前价格差：%f",buy,sell,(sell - buy)/buy*100.0)
			Listen(c,sy,sell,buy)
		}
		time.Sleep(time.Second)
	}
}

func Listen(c exchange.ExClient,sy int,s, b float64)  {
	//
	mark := time.Now().String()
	for {
		//sell := c.GetLastSellPrice(sy)
		buy := c.GetLastBuyPrice(sy)
		if (b - buy)/b > 0.03 {
			// 价格下跌百分之3，放弃监听
			log.Warn("价格下跌百分之3，放弃监听")
			flagYCC = 0
			break
		}
		sells := c.GetSellP(sy)
		// 第一个挂单卖价和第二个的价格差
		t1 := (sells[1].Price - sells[0].Price)/sells[0].Price
		t2 := (sells[2].Price - sells[0].Price)/sells[0].Price
		// 监测别人低价卖
		if  (t1 > 0.05 && t1 < 0.2) || (t2 > 0.05 && t2 < 0.2){
			// todo 判断是否有余额
			log.Warnf(">>>监测到有人低价卖")
			count := 0.0
			// 最高只买100美元
			if sells[0].Price * sells[0].Count > 100.0 {
				count = 100.0 / sells[0].Price
			}else {
				count = sells[0].Count
			}
			// 买入
			log.Warnf(">>>mark:%s, 买入%s, 买入价格%f， 买入数量: %f",mark, Symbol[sy], sells[0].Price, count)
			// 监控卖单，等待卖出
			go ListenSell(c,sy,sells[0].Price,count,mark)
		}
		//  检测别人高价买
		buys := c.GetBuyP(sy)
		// 第一个挂单卖价和第二个的价格差
		t1 = (buys[0].Price - buys[1].Price)/buys[0].Price
		t2 = (buys[0].Price - buys[2].Price)/buys[0].Price
		t3 := (s -b)/b
		if  (t1 > 0.05 && t1 < 0.2  && t1 > t3*0.6 ) || (t2 > 0.05 && t2 < 0.2 && t1 > t2*0.6){
			// todo 判断是否有余额
			// 卖出
			count := 0.0
			// 最高只买100美元
			if buys[0].Price * buys[0].Count > 100.0 {
				count = 100.0 / buys[0].Price
			}
			// 卖出
			log.Warnf(">>>mark:%s, 卖出%s, 卖出价格%f， 卖出数量: %f",mark, Symbol[sy], buys[0].Price, count)
			// 等待低价买回来
		}
		time.Sleep(time.Second)
	}

}


// 高价卖出
func ListenSell(c exchange.ExClient,sy int,price,count float64,mark string)  {
	// 获取买单，获利超过百分之4就抛售
	log.Warn(">>>等待抛售：mark:", mark)
	// 买入花费
	buyP := count*price
	// 卖出赚取
	sellP := 0.0
	for count > 0 {
		// 获取买单价格
		buys := c.GetBuyP(sy)
		// 目标价格
		targetP := price + price*0.04
		if buys[0].Price > targetP {
			if count > buys[0].Count{
				sellP += buys[0].Count * buys[0].Price
			}else {
				sellP += count * buys[0].Price
			}
			count = count - buys[0].Count

			log.Warnf(">>>mark: %s, 卖出%s, price: %f, count: %f, 盈利%f",
				mark, Symbol[sy], buys[0].Price, buys[0].Count,sellP-buyP)
		}
		time.Sleep(time.Second)
	}
}


// 低价买回
func ListenBuy(c exchange.ExClient,sy int,price,count float64,mark string)  {
	// 获取卖单，获利超过百分之4就抛售

}