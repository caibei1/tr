package exchange

import "sync"

// 公共接口
type ExClient interface {
	// 获取最新购买价
	GetLastBuyPrice(symbol int) float64

	// 获取最新卖出价
	GetLastSellPrice(symbol int) float64

	// 获取最新人民币成交价
	//GetLastSuccessRMBPrice(symbol int) float64
}


// 公共结构体
// 包含币的各种信息，所需要的所有信息从此处拿
type Quotation struct {
	Lastrmb float64 // 最新成交价
	Range float64 // 涨幅
	Open float64
	Buys []Tr // 买单
	Sells []Tr // 卖单
	Histroy []Tr // 历史成交
	MyHistroy []Tr //我的成交历史
	MyBuys []Tr  // 我的买单
	MySell []Tr // 我的卖单
	sync.RWMutex
	IsBuy bool //是否有买单
	IsSell bool //是否有卖单
	Symbol int  // 交易对
	SymbolStr string  // 交易对
}

// 挂单 价格数量
type Tr struct {
	Price float64
	Count float64
	Point int  // 当前指针
	Date string // 时间
	Success  int // 是否成功，在我的历史里的字段
}



const YCC  = 0
const BTY  = 1
const BTC  = 2
const BCH  = 3
const ETH  = 4
const ETC  = 5
const ZEC  = 6
const LTC  = 7