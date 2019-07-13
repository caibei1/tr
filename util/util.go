package util

import log "github.com/sirupsen/logrus"

// 判断价格差是否达到1%
// 买方是大平台 卖方是小平台
func IsOnePercent(buy,sell float64) bool {
	percent := (buy - sell)/sell
	//log.Warnf("大平台buy: %f, 小平台sell: %f, diff %f", buy ,sell, percent*100.0)
	//
	if percent > 0.0015 && percent < 0.2{
		log.Warnf(">>>大平台buy: %f, 小平台sell: %f, diff %f", buy ,sell, percent*100.0)
		return true
	}
	return false
}
