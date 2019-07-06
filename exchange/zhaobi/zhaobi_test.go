package zhaobi

import (
	"testing"
	"fmt"
)

func TestGetAccount(t *testing.T) {
	resp := GetAccount()
	fmt.Println(resp)
}