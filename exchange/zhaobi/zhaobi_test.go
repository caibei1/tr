package zhaobi

import (
	"testing"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

func TestGetAccount(t *testing.T) {
	resp := GetAccount()
	fmt.Printf("%+v",resp)
}

func TestBtyYcc(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	zbc := NewZBClient()
	zbc.Init()
	// 等待初始化完成
	time.Sleep(time.Second*5)



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


	// 获取余额
	acc := GetAccount()
	usdt := 0.5


	PlaceOrder(usdt/btySells[0].Price, "BTY","USDT", btySells[0].Price,"BUY")


	//bty买入ycc
	//allBty := float64(acc.Data.List.BTY.Active)
	allBty := 2.0
	PlaceOrder(allBty*(1-0.002)/YCCBTYSells[0].Price,"YCC","BTY",YCCBTYSells[0].Price,"BUY")

	// 卖出YCC
	allYCC := float64(acc.Data.List.YCC.Active)
	PlaceOrder(allYCC,"YCC","USDT",yccBuys[0].Price,"SELL")
}