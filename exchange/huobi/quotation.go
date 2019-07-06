package huobi

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
	pub "tr/exchange"
)

// 文档地址
// https://github.com/huobiapi/API_Docs/wiki/REST_api_reference


type HBClient struct {

}

func NewHBClient() *HBClient {
	return &HBClient{}
}

func (*HBClient)Init()  {
	// 初始化交易对
	initSymbol()
	// 初始化数据
	InitHBQuotation()
	// 同步数据
	go func() {
		for  {
			SyncMarketInfo()
			time.Sleep(time.Second)
		}

	}()

}

var HBQuotationMap map[int]*pub.Quotation

func InitHBQuotation()  {
	HBQuotationMap =  make(map[int]*pub.Quotation,len(Symbol))
	for k,v := range Symbol {
		q :=  &pub.Quotation{}
		q.Lock()
		HBQuotationMap[k] = q
		q.Symbol = k
		q.SymbolStr = v
		q.Unlock()
	}
}



// 交易对
var Symbol map[int]string
// 反交易对
var Symbol1 map[string]int
func initSymbol()  {
	Symbol = make(map[int]string,10)
	//Symbol[pub.YCC] = "YCCUSDT"
	//Symbol[pub.BTY] = "BTYUSDT"
	Symbol[pub.BTC] = "btcusdt"
	Symbol[pub.BCH] = "bchusdt"
	Symbol[pub.ETH] = "ethusdt"
	Symbol[pub.ETC] = "etcusdt"
	Symbol[pub.ZEC] = "zecusdt"
	Symbol[pub.LTC] = "ltcusdt"

	Symbol1 = make(map[string]int,10)
	//Symbol1["YCCUSDT"] = pub.YCC
	//Symbol1["BTYUSDT"] = pub.BTY
	Symbol1["btcusdt"] = pub.BTC
	Symbol1[ "bchusdt"] = pub.BCH
	Symbol1["ethusdt"] = pub.ETH
	Symbol1["etcusdt"] = pub.ETC
	Symbol1["zecusdt"] = pub.ZEC
	Symbol1["ltcusdt"] = pub.LTC

}

//"tick": {
//    "id": K线id（时间戳）,
//    "amount": 成交量,
//    "count": 成交笔数,
//    "open": 开盘价,
//    "close": 收盘价,当K线为最晚的一根时，是最新成交价
//    "low": 最低价,
//    "high": 最高价,
//    "vol": 成交额, 即 sum(每一笔成交价 * 该笔的成交量)
//    "bid": [买1价,买1量],
//    "ask": [卖1价,卖1量]
//  }
func GetHBInfo(sy int) *HBInfoResp {
	url := "https://api.huobi.br.com/market/detail/merged?symbol="+Symbol[sy]
	resp,err := http.Get(url)
	if err != nil {
		log.Error("GetIndexInfo Get 获取火币信息market/detail/merged失败：",err)
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("GetIndexInfo ReadAll 获取火币信息market/detail/merged失败：",err)
		return nil
	}
	var info = &HBInfoResp{}
	err = json.Unmarshal(body,info)
	if err != nil {
		log.Error("GetIndexInfo ReadAll 获取火币信息market/detail/merged失败：",err)
		return nil
	}
	if len(info.Tick.Ask) < 2 {
		log.Info(info)
		log.Info("获取火币信息market/detail/merged失败")
		return nil
	}
	log.Info("获取火币信息market/detail/merged成功")
	return info
}

type HBInfoResp	 struct {
	Tick HBTick `json:"tick"`
}

type HBTick struct {
	Id int64 `json:"id"`
	Amount float64 `json:"amount"`
	Count int64 `json:"count"`
	Open float64 `json:"open"`
	Close float64 `json:"close"`
	Low float64 `json:"low"`
	High float64 `json:"high"`
	Vol float64`json:"vol"`
	Bid []float64 `json:"bid"`
	Ask []float64 `json:"ask"`
}



func SyncMarketInfo()  {
	for k,_ := range Symbol  {
		resp := GetHBInfo(k)
		if resp != nil{
			q := HBQuotationMap[k]
			q.Lock()
			// todo
			// 买
			q .Buys = make([]pub.Tr,1,20)
			q.Buys[0] = pub.Tr{Price:resp.Tick.Bid[0],Count:resp.Tick.Bid[1]}
			// 卖
			q .Sells = make([]pub.Tr,1,20)
			q.Sells[0] = pub.Tr{Price:resp.Tick.Ask[0],Count:resp.Tick.Ask[1]}
			q.Unlock()
		}
	}
}

func (*HBClient)GetLastBuyPrice(symbol int) float64  {
	q := HBQuotationMap[symbol]
	if len(q.Buys) == 0 {
		return 0
	}
	return q.Buys[0].Price
}

// 获取最新卖出价
func (*HBClient)GetLastSellPrice(symbol int) float64{
	q := HBQuotationMap[symbol]
	return q.Sells[0].Price
}


// todo
func (*HBClient) GetBuyP(symbol int) []pub.Tr{
	return nil
}

func (*HBClient) GetSellP(symbol int) []pub.Tr{
	return nil
}