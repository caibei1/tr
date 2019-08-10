package zhaobi

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
	"strconv"
	"sync"
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



		if usdtCount > 5.0*(1+0.004) {
			log.Warn("=====in usdt买bty  bty换ycc  然后卖掉ycc换usdt =====")
			log.Warnf("usdtCount: %f",usdtCount)

			// 获取每个挂单的金额
			if btySells[0].Price*btyBuys[0].Count < 2.0 {
				time.Sleep(time.Second)
				continue
			}

			if YCCBTYSells[0].Price* YCCBTYSells[0].Count < 2.0 {
				time.Sleep(time.Second)
				continue
			}

			if yccBuys[0].Price*yccBuys[0].Count < 2.0 {
				time.Sleep(time.Second)
				continue
			}

			// 获取余额
			acc := GetAccount()
			usdt := 4.0
			wg := sync.WaitGroup{}
			wg.Add(3)

			// 买入bty
			// 获取usdt
			allUsdt := float64(acc.Data.List.USDT.Active)
			if allUsdt < 4.0 {

				time.Sleep(time.Second*5)
				continue
				usdt = allUsdt
			}

			go PlaceOrder(wg,usdt/btySells[0].Price, "BTY","USDT", btySells[0].Price,"BUY")

			// bty买入ycc
			// 获取所有bty
			//var allBty float64
			//allBty = float64(GetAccount().Data.List.BTY.Active)
			//if allBty < 1.0 {
			//	time.Sleep(time.Second*2)
			//	allBty = float64(GetAccount().Data.List.BTY.Active)
			//}
			//if allBty < 1.0 {
			//	log.Error("1 allBty < 1.0")
			//	continue
			//}

			allBty := usdt/btySells[0].Price

			go PlaceOrder(wg,allBty*(1-0.001)/YCCBTYSells[0].Price,"YCC","BTY",YCCBTYSells[0].Price,"BUY")

			// 卖出YCC
			var allYCC float64

			//allYCC = float64(GetAccount().Data.List.YCC.Active)
			//if allYCC < 10.0 {
			//	time.Sleep(time.Second*2)
			//	allYCC = float64(GetAccount().Data.List.YCC.Active)
			//}
			//if allYCC < 10.0 {
			//	log.Error("1 allYCC < 10.0")
			//	continue
			//}

			allYCC = allBty*(1-0.001)/YCCBTYSells[0].Price

			go PlaceOrder(wg,allYCC,"YCC","USDT",yccBuys[0].Price,"SELL")

			wg.Wait()
			end := float64(GetAccount().Data.List.USDT.Active)
			if end - allUsdt < 1 {
				time.Sleep(time.Second*2)
			}
			log.Warnf("原始usdt：%f,成交后usdt：%f, 盈利：%f",allUsdt, end, end - allUsdt)


		}
		time.Sleep(time.Second*1)
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



		if usdtCount > 5.0*(1+0.004) {
			log.Warn("=====in  usdt买ycc  ycc换bty  然后卖掉bty换usdt=====")
			log.Warnf("usdtCount: %f",usdtCount)
			// 获取每个挂单的金额
			if btySells[0].Price*btyBuys[0].Count < 2.0 {
				time.Sleep(time.Second)
				continue
			}

			if YCCBTYSells[0].Price* YCCBTYSells[0].Count < 2.0 {
				time.Sleep(time.Second)
				continue
			}

			if yccBuys[0].Price*yccBuys[0].Count < 2.0 {
				time.Sleep(time.Second)
				continue
			}


			// 获取余额
			acc := GetAccount()
			usdt := 4.0

			wg := sync.WaitGroup{}
			wg.Add(3)

			// 获取usdt
			allUsdt := float64(acc.Data.List.USDT.Active)
			if allUsdt < 4.0 {
				time.Sleep(time.Second*2)
				continue
				usdt = allUsdt
			}

			// 买入ycc
			go PlaceOrder(wg,usdt/yccSells[0].Price, "YCC","USDT", yccSells[0].Price,"BUY")

			// 卖出ycc 获得bty
			var allYCC float64
			//allYCC = float64(GetAccount().Data.List.YCC.Active)
			//if allYCC < 10.0 {
			//	time.Sleep(time.Second*2)
			//	allYCC = float64(GetAccount().Data.List.YCC.Active)
			//}
			//if allYCC < 10.0 {
			//	log.Error("2 allYCC < 10.0")
			//	time.Sleep(time.Millisecond*500)
			//	continue
			//}
			allYCC = usdt/yccSells[0].Price
			go PlaceOrder(wg,allYCC,"YCC","BTY",YCCBTYBuys[0].Price,"SELL")

			// 卖出bty
			var allBTY = allYCC*YCCBTYBuys[0].Price
			//allBTY = float64(GetAccount().Data.List.BTY.Active)
			//if allBTY < 1.0 {
			//	time.Sleep(time.Second)
			//	allBTY = float64(GetAccount().Data.List.BTY.Active)
			//}
			//if allBTY < 1.0 {
			//	log.Error("2 allBTY < 1.0")
			//	time.Sleep(time.Millisecond*500)
			//	continue
			//}
			go PlaceOrder(wg,allBTY,"BTY","USDT",btyBuys[0].Price,"SELL")
			wg.Wait()
			end := float64(GetAccount().Data.List.USDT.Active)
			log.Warnf("原始usdt：%f,成交后usdt：%f, 盈利：%f",allUsdt, end, end - allUsdt)
			if end - allUsdt < -1 {
				time.Sleep(time.Second*2)
			}
			time.Sleep(time.Millisecond*300)

		}
		time.Sleep(time.Millisecond*500)
	}

}

// 0 买 1卖
func PlaceOrder(wg sync.WaitGroup,amout float64, currency, currency2 string, price float64, ty string) bool {
	log.Warnf("amout: %f, currency： %s, currency2:  %s, price: %f, ty: %s", amout,currency,currency2,price,ty)
	defer wg.Done()
	a := FloatToString(price)[:8]
	if currency == "BTY" {
		a = FloatToString(price)[:6]
	}
	fmt.Println("价格: ",a)

	b := ""
	if currency == "YCC" {
		b = strings.Split(FloatToString(amout),".")[0]
	}

	if currency == "BTY" {
		b = strings.Split(FloatToString(amout),".")[0] + "." + string(strings.Split(FloatToString(amout),".")[1][1])
		//b = FloatToString(amout)[:4]
	}
	fmt.Println("数量：",b)
	PostBill(b,currency,currency2,a,ty)

	return true
}

func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}
