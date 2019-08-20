package wechat

import (
	"fmt"
	"testing"
	"time"
)

func Test(t *testing.T) {
	backoff := &Backoff{Max: time.Second * 5}
	for backoff.Wait() {
		fmt.Println(backoff)
	}
}
