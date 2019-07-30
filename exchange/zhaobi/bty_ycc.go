package zhaobi

import (
	log "github.com/sirupsen/logrus"
	"time"
	"strconv"
)

// 比特元ycc交易对
func BtyYcc(zbc *ZBClient)  {
	//account := GetAccount()

	log.Warn("start")

	for true {
		// 先获取ycc-usdt  买和卖价格，数量
		yccBuys := zbc.GetBuyP(Symbol1["YCCUSDT"])
		yccSells := zbc.GetSellP(Symbol1["YCCUSDT"])

		log.Debugf("ycc buy %+v", yccBuys)
		log.Debugf("ycc sell %+v", yccSells)

		// 获取bty-usdt
		btyBuys := zbc.GetBuyP(Symbol1["BTYUSDT"])
		btySells := zbc.GetSellP(Symbol1["BTYUSDT"])

		log.Debugf("bty buy %+v", btyBuys)
		log.Debugf("bty sell %+v", btySells)

		// 获ycc-bty
		YCCBTYBuys := zbc.GetBuyP(Symbol1["YCCBTY"])
		YCCBTYSells := zbc.GetSellP(Symbol1["YCCBTY"])

		log.Debugf("ycc-bty buy %+v", YCCBTYBuys)
		log.Debugf("ycc-bty sell %+v", YCCBTYSells)

		// usdt买bty  bty换ycc  然后卖掉ycc换usdt
		// 初始5.0usdt
		btyCount := 5.0/btySells[0].Price
		yccCount := btyCount/YCCBTYSells[0].Price
		usdtCount := yccCount*yccBuys[0].Price

		log.Warnf("usdtCount: %f",usdtCount)

		if usdtCount > 5.0*(1+0.005) {
			log.Warn("=====in=====")
			// 获取每个挂单的金额
			if btySells[0].Price*btyBuys[0].Count < 2.0 {
				continue
			}

			if YCCBTYSells[0].Price* YCCBTYSells[0].Count < 2.0 {
				continue
			}

			if yccBuys[0].Price*yccBuys[0].Count < 2.0 {
				continue
			}

			log.Warn("1111111")
			// 获取余额
			acc := GetAccount()
			usdt := 2.0


			// 买入bty
			// 获取usdt
			allUsdt := float64(acc.Data.List.USDT.Active)
			if allUsdt < 2.0 {
				usdt = allUsdt
			}

			PlaceOrder(usdt/btySells[0].Price, "BTY","USDT", btySells[0].Price,"BUY")

			// bty买入ycc
			allBty := float64(acc.Data.List.BTY.Active)
			PlaceOrder(allBty*(1-0.002)/YCCBTYSells[0].Price,"YCC","BTY",YCCBTYSells[0].Price,"BUY")

			// 卖出YCC
			allYCC := float64(acc.Data.List.YCC.Active)
			PlaceOrder(allYCC,"YCC","USDT",yccBuys[0].Price,"SELL")

			continue
		}
		time.Sleep(time.Millisecond*300)
	}

}


func YCCBTY(zbc *ZBClient)  {
	//account := GetAccount()

	for true {
		// 先获取ycc-usdt  买和卖价格，数量
		yccBuys := zbc.GetBuyP(Symbol1["YCCUSDT"])
		yccSells := zbc.GetSellP(Symbol1["YCCUSDT"])

		log.Debugf("ycc buy %+v", yccBuys)
		log.Debugf("ycc sell %+v", yccSells)

		// 获取bty-usdt
		btyBuys := zbc.GetBuyP(Symbol1["BTYUSDT"])
		btySells := zbc.GetSellP(Symbol1["BTYUSDT"])

		log.Debugf("bty buy %+v", btyBuys)
		log.Debugf("bty sell %+v", btySells)

		// 获ycc-bty
		YCCBTYBuys := zbc.GetBuyP(Symbol1["YCCBTY"])
		YCCBTYSells := zbc.GetSellP(Symbol1["YCCBTY"])

		log.Debugf("ycc-bty buy %+v", YCCBTYBuys)
		log.Debugf("ycc-bty sell %+v", YCCBTYSells)

		// usdt买ycc  ycc换bty  然后卖掉bty换usdt
		// 初始5.0usdt
		yccCount := 5.0/yccSells[0].Price
		// ycc 卖掉
		btyCount := yccCount*YCCBTYBuys[0].Price
		usdtCount := btyCount*btyBuys[0].Price

		log.Warnf("usdtCount: %f",usdtCount)

		if usdtCount > 5.0*(1+0.005) {
			log.Warn("=====in=====")
			// 获取每个挂单的金额
			if btySells[0].Price*btyBuys[0].Count < 2.0 {
				continue
			}

			if YCCBTYSells[0].Price* YCCBTYSells[0].Count < 2.0 {
				continue
			}

			if yccBuys[0].Price*yccBuys[0].Count < 2.0 {
				continue
			}


			// 获取余额
			acc := GetAccount()
			usdt := 2.0



			// 获取usdt
			allUsdt := float64(acc.Data.List.USDT.Active)
			if allUsdt < 2.0 {
				usdt = allUsdt
			}

			log.Warn("222222222")
			// 买入ycc
			PlaceOrder(usdt/yccSells[0].Price, "YCC","USDT", yccSells[0].Price,"BUY")

			// 卖出ycc 获得bty
			allYCC := float64(acc.Data.List.YCC.Active)
			PlaceOrder(allYCC,"YCC","BTY",YCCBTYBuys[0].Price,"SELL")

			// 卖出bty
			allBTY := float64(acc.Data.List.BTY.Active)
			PlaceOrder(allBTY,"BTY","USDT",btyBuys[0].Price,"SELL")
			continue

		}
		time.Sleep(time.Millisecond*300)
	}

}

// 0 买 1卖
func PlaceOrder(amout float64, currency, currency2 string, price float64, ty string) bool {
	log.Warnf("amout: %f, currency： %s, currency2:  %s, price: %f, ty: %s", amout,currency,currency2,price,ty)

	//PostBill(FloatToString(amout),currency,currency2,FloatToString(price),ty)

	return true
}

func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}
