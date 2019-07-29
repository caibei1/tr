package zhaobi

import (
	"testing"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

func TestGetAccount(t *testing.T) {
	resp := GetAccount()
	fmt.Println(resp)
}

func TestBtyYcc(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	zbc := NewZBClient()
	zbc.Init()
	// 等待初始化完成
	time.Sleep(time.Second*2)
	BtyYcc(zbc)
}