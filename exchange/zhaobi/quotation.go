package zhaobi

import (
	"net/http"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"encoding/json"
	"sync"
	"strconv"
	"time"
)

// 获取找币行情  只做usdt的

type ZBClient struct {

}

func NewZBClient() *ZBClient {
	return &ZBClient{}
}

var ZBQuotationMap map[int]*ZBQuotation

const YCC  = 0
const BTY  = 1
const BTC  = 2
const BCH  = 3
const ETH  = 4
const ETC  = 5
const ZEC  = 6
const LTC  = 7


// 包含币的各种信息，所需要的所有信息从此处拿
type ZBQuotation struct {
	Lastrmb float64 // 最新成交价
	Range float64 // 涨幅
	Open float64
	Buys []ZBTr // 买单
	Sells []ZBTr // 卖单
	Histroy []ZBTr // 历史成交
	MyHistroy []ZBTr //我的成交历史
	MyBuys []ZBTr  // 我的买单
	MySell []ZBTr // 我的卖单
	sync.RWMutex
	IsBuy bool //是否有买单
	IsSell bool //是否有卖单
	Symbol int  // 交易对
	SymbolStr string  // 交易对
}


// 挂单 价格数量
type ZBTr struct {
	Price float64
	Count float64
	Point int  // 当前指针
	Date string // 时间
	Success  int // 是否成功，在我的历史里的字段
}


// 加载币的数据
func Init()  {
	// 初始化交易对
	InitSymbol()
	
	// 初始化数据
	InitZBQuotation()

	// 同步数据
	go func() {
		for   {
			SyncIndexInfo()
			time.Sleep(time.Second)
		}

	}()

	// 同步数据
	go func() {
		for   {
			SyncMarketInfo()
			time.Sleep(time.Second)
		}

	}()


}

// 交易对
var Symbol map[int]string
// 反交易对
var Symbol1 map[string]int
// 初始化交易对
func InitSymbol()  {
	Symbol = make(map[int]string,10)
	Symbol[YCC] = "YCCUSDT"
	Symbol[BTY] = "BTYUSDT"
	Symbol[BTC] = "BTCUSDT"
	Symbol[BCH] = "BCHUSDT"
	Symbol[ETH] = "ETHUSDT"
	Symbol[ETC] = "ETCUSDT"
	Symbol[ZEC] = "ZECUSDT"
	Symbol[LTC] = "LTCUSDT"

	Symbol1 = make(map[string]int,10)
	Symbol1["YCCUSDT"] = YCC
	Symbol1["BTYUSDT"] = BTY
	Symbol1["BTCUSDT"] = BTC
	Symbol1[ "BCHUSDT"] = BCH
	Symbol1["ETHUSDT"] = ETH
	Symbol1["ETCUSDT"] = ETC
	Symbol1["ZECUSDT"] = ZEC
	Symbol1["LTCUSDT"] = LTC

}

func GetSymnol(i int) string {
	return Symbol[i]
}

type ZBIndexInfo struct {
	Data ZBData `json:"data"`
}

type ZBData struct {
	USDT []USDT `json:"USDT"`
}

type USDT struct {
	Buy     FNumber   `json:"buy"`
	Sell    FNumber   `json:"sell"`
	Open    FNumber   `json:"open"`
	Lastrmb FNumber   `json:"lastrmb"` // 最新成交价，人民币
	High    FNumber   `json:"high"` // 今日最高价
	Low     FNumber   `json:"low"`  // 今日最低价
	Vol     FNumber   `json:"vol"`
	Range   Range   `json:"range"`  // 涨幅
	Symbol  string   `json:"symbol"` // 交易对 BTYBTC
}

// 首页接口
//"code":200,
//    "ecode":"200",
//    "error":"OK",
//    "message":"OK",
//    "data":{
//        "BTC":[
//            {
//                "buy":"0.0000414",
//                "sell":"0.0000421",
//                "open":"0.0000424",
//                "last":"0.0000419",
//                "lastrmb":"3.6032",
//                "high":"0.0000432",
//                "low":"0.0000416",
//                "vol":"6393571.5",
//                "range":"-1.18%",
//                "symbol":"BTYBTC",
//                "plat":"fxee",
//                "date":"2019-06-30 11:43:32"
//            },
func GetZBIndexInfo() *ZBIndexInfo {
	resp,err := http.Get(`https://api.biqianbao.top/api/data/Ticker?sort=cname`)
	if err != nil {
		log.Error("GetIndexInfo Get 获取首页信息失败：",err)
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("GetIndexInfo ReadAll 获取首页信息失败：",err)
		return nil
	}
	var info = &ZBIndexInfo{}
	err = json.Unmarshal(body,info)
	if err != nil {
		log.Error("GetIndexInfo ReadAll 获取首页信息失败：",err)
		return nil
	}
	if len(info.Data.USDT) < 2 || info.Data.USDT[0].Sell <= 0 {
		log.Info(info)
		log.Info("获取首页信息失败")
		return nil
	}
	log.Info("获取首页信息成功")
	return info
}

// 1卖0买
// 获取币的详情信息 买单卖单成交
// symbol交易对
func GetZBMarketInfo(num int,symbol int) *ZBMarketInfoResp {

	symbolStr := GetSymnol(symbol)
	url := `https://api.biqianbao.top/api/data/market?num=`+strconv.Itoa(num)+`&format=&symbol=`+symbolStr
	resp,err := http.Get(url)
	log.Debug(url)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("GetIndexInfo ReadAll 获取市场信息失败：",err)
		return nil
	}
	var info = &ZBMarketInfoResp{}
	err = json.Unmarshal(body,info)
	if err != nil {
		log.Error("GetIndexInfo ReadAll 获取市场信息失败：",err)
		return nil
	}
	log.Debug("获取市场信息成功")
	if len(info.ZBMarketInfo.MarketData.Buys) == 0 || len(info.ZBMarketInfo.MarketData.Sell) == 0{
		log.Info(info.ZBMarketInfo.MarketData.Buys)
		log.Info(info.ZBMarketInfo.MarketData.Sell)
		log.Error("买，卖为空， 获取市场信息失败：",err)
		return nil
	}
	log.Debug(info)
	log.Info("获取市场信息成功")
	return info

}

type ZBMarketInfoResp struct {
	ZBMarketInfo ZBMarketInfo `json:"data"`
} 

type ZBMarketInfo struct {
	Trade []Trade `json:"trade"`
	MarketData MarketData `json:"marketdata"`
}

//"price":"0.5105",
// "am":"500.0",
// "time":"13:17:27",
// "type":"1"
// 1卖0买
// 历史成交信息
type Trade struct {
	Price FNumber `json:"price"`
	Am FNumber `json:"am"`
	Time string `json:"time"`
	Type INumber `json:"type"`
}

type MarketData struct {
	Buys []Trade `json:"buy"`
	Sell []Trade `json:"sell"`
}


type FNumber float64
type INumber int
type Range float64


func (t *FNumber) UnmarshalJSON(b []byte) error {
	v, err := strconv.ParseFloat(string(b[1:len(b)-1]),64)
	if err != nil {
		return err
	}
	*t = FNumber(v)
	return nil
}

func (t *INumber) UnmarshalJSON(b []byte) error {
	v, err := strconv.Atoi(string(b[1:len(b)-1]))
	if err != nil {
		return err
	}
	*t = INumber(v)
	return nil
}

func (t *Range) UnmarshalJSON(b []byte) error {
	v, err := strconv.ParseFloat(string(b[1:len(b)-2]),64)
	if err != nil {
		return err
	}
	*t = Range(v)
	return nil
}


// 同步首页信息
func SyncIndexInfo()  {
	info := GetZBIndexInfo()
	if info != nil {
		datas := info.Data.USDT
		for _,d := range datas {
			if i,ok := Symbol1[d.Symbol]; ok {
				q := ZBQuotationMap[i]
				q.Lock()
				q.Lastrmb = float64(d.Lastrmb)
				q.Range = float64(d.Range)
				q.Open = float64(d.Open)
				// todo
				q.Unlock()
			}
		}
	}

}

func SyncMarketInfo()  {
	for k,_ := range Symbol  {
		resp := GetZBMarketInfo(20,k)
		if resp != nil{
			q := ZBQuotationMap[k]
			q.Lock()
			// todo
			// 买
			q .Buys = make([]ZBTr,0,20)
			for _,v := range resp.ZBMarketInfo.MarketData.Buys{
				zBTr := ZBTr{
					Price:float64(v.Price),
					Count:float64(v.Am),
				}
				q .Buys = append(q.Buys,zBTr)
			}
			// 卖
			q .Sells = make([]ZBTr,0,20)
			for _,v := range resp.ZBMarketInfo.MarketData.Sell{
				zBTr := ZBTr{
					Price:float64(v.Price),
					Count:float64(v.Am),
				}
				q .Sells = append(q.Sells,zBTr)
			}
			q.Unlock()
		}
	}
}

func InitZBQuotation()  {
	ZBQuotationMap =  make(map[int]*ZBQuotation,len(Symbol))
	for k,v := range Symbol {
		q :=  &ZBQuotation{}
		q.Lock()
		ZBQuotationMap[k] = q
		q.Symbol = k
		q.SymbolStr = v
		q.Unlock()
	}
}

func GetZBQuotation(symbol int) *ZBQuotation {
	return ZBQuotationMap[symbol]
}

// 获取最新购买价
func (*ZBClient) GetLastBuyPrice(symbol int) float64 {
	q := ZBQuotationMap[symbol]
	q.RLock()
	defer q.RUnlock()
	buys := q.Buys
	var max = 0.0
	for _,v := range buys {
		if v.Price > max {
			max = v.Price
		}
	}
	return max
}

// 获取最新卖出价
func (*ZBClient) GetLastSellPrice(symbol int) float64 {
	q := ZBQuotationMap[symbol]
	q.RLock()
	defer q.RUnlock()
	sells := q.Sells
	var min = 10000000.0
	for _,v := range sells {
		if v.Price < min {
			min = v.Price
		}
	}
	return min
}

// 或者最新人民币成交价
func (*ZBClient) GetLastSuccessRMBPrice(symbol int) float64 {
	q := ZBQuotationMap[symbol]
	return q.Lastrmb
}


func (*ZBClient) GetOpen(symbol int) float64 {
	q := ZBQuotationMap[symbol]
	return q.Open
}