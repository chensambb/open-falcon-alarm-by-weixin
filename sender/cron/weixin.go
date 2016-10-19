package cron

import (
	"github.com/open-falcon/sender/g"
	"github.com/open-falcon/sender/model"
	"github.com/open-falcon/sender/proc"
	"github.com/open-falcon/sender/redis"
	"github.com/toolkits/net/httplib"
	"log"
	"time"
)

func ConsumeWeixin() {
	queue := g.Config().Queue.Weixin
	for {
		L := redis.PopAllWeixin(queue)
		if len(L) == 0 {
			time.Sleep(time.Millisecond * 200)
			continue
		}
		SendWeixinList(L)
	}
}

func SendWeixinList(L []*model.Weixin) {
	for _, Weixin := range L {
		WeixinWorkerChan <- 1
		go SendWeixin(Weixin)
	}
}

func SendWeixin(Weixin *model.Weixin) {
	defer func() {
		<-WeixinWorkerChan
	}()

	url := g.Config().Api.Weixin
	r := httplib.Post(url).SetTimeout(5*time.Second, 2*time.Minute)
	r.Param("tos", Weixin.Tos)
	r.Param("content", Weixin.Content)
	resp, err := r.String()
	if err != nil {
		log.Println(err)
	}

	proc.IncreWeixinCount()

	if g.Config().Debug {
		log.Println("==Weixin==>>>>", Weixin)
		log.Println("<<<<==Weixin==", resp)
	}

}
