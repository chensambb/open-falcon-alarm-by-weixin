package proc

import (
	"sync/atomic"
)

var smsCount, mailCount, weixinCount uint32

func GetSmsCount() uint32 {
	return atomic.LoadUint32(&smsCount)
}

func GetMailCount() uint32 {
	return atomic.LoadUint32(&mailCount)
}

func GetWeixinCount() uint32 {
	return atomic.LoadUint32(&weixinCount)
}

func IncreSmsCount() {
	atomic.AddUint32(&smsCount, 1)
}

func IncreMailCount() {
	atomic.AddUint32(&mailCount, 1)
}

func IncreWeixinCount() {
	atomic.AddUint32(&weixinCount, 1)
}
