package redis

import (
	"encoding/json"
	"github.com/open-falcon/alarm/g"
	"github.com/open-falcon/sender/model"
	"log"
	"strings"
)

func LPUSH(queue, message string) {
	rc := g.RedisConnPool.Get()
	defer rc.Close()
	_, err := rc.Do("LPUSH", queue, message)
	if err != nil {
		log.Println("LPUSH redis", queue, "fail:", err, "message:", message)
	}
}

func WriteSmsModel(sms *model.Sms) {
	if sms == nil {
		return
	}

	bs, err := json.Marshal(sms)
	if err != nil {
		log.Println(err)
		return
	}

	LPUSH(g.Config().Queue.Sms, string(bs))
}

func WriteMailModel(mail *model.Mail) {
	if mail == nil {
		return
	}

	bs, err := json.Marshal(mail)
	if err != nil {
		log.Println(err)
		return
	}

	LPUSH(g.Config().Queue.Mail, string(bs))
}

func WriteWeixinModel(weixin *model.Weixin) {
	if weixin == nil {
		return
	}

	bs, err := json.Marshal(weixin)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("g.Config().Queue.Weixin:", g.Config().Queue.Weixin)
	log.Println("WriteWeixinModel:", string(bs))
	LPUSH(g.Config().Queue.Weixin, string(bs))
}

func WriteSms(tos []string, content string) {
	if len(tos) == 0 {
		return
	}

	sms := &model.Sms{Tos: strings.Join(tos, ","), Content: content}
	WriteSmsModel(sms)
}

func WriteWeixin(tos []string, content string) {
	if len(tos) == 0 {
		return
	}

	weixin := &model.Weixin{Tos: strings.Join(tos, ","), Content: content}
	WriteWeixinModel(weixin)
}

func WriteMail(tos []string, subject, content string) {
	if len(tos) == 0 {
		return
	}

	mail := &model.Mail{Tos: strings.Join(tos, ","), Subject: subject, Content: content}
	WriteMailModel(mail)
}
