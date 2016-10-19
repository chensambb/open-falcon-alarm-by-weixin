package cron

import (
	"github.com/open-falcon/sender/g"
)

var (
	SmsWorkerChan    chan int
	MailWorkerChan   chan int
	WeixinWorkerChan chan int
)

func InitWorker() {
	workerConfig := g.Config().Worker
	SmsWorkerChan = make(chan int, workerConfig.Sms)
	MailWorkerChan = make(chan int, workerConfig.Mail)
	WeixinWorkerChan = make(chan int, workerConfig.Weixin)
}
