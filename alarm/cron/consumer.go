package cron

import (
	"encoding/json"
	"github.com/open-falcon/alarm/api"
	"github.com/open-falcon/alarm/g"
	"github.com/open-falcon/alarm/redis"
	"github.com/open-falcon/common/model"
	"log"
)

func consume(event *model.Event, isHigh bool) {
	actionId := event.ActionId()
	if actionId <= 0 {
		return
	}

	action := api.GetAction(actionId)
	if action == nil {
		return
	}

	if action.Callback == 1 {
		HandleCallback(event, action)
		return
	}

	if isHigh {
		consumeHighEvents(event, action)
	} else {
		consumeLowEvents(event, action)
	}
}

// 高优先级的不做报警合并
func consumeHighEvents(event *model.Event, action *api.Action) {
	if action.Uic == "" {
		return
	}

	phones, mails, weixins := api.ParseTeams(action.Uic)

	smsContent := GenerateSmsContent(event)
	mailContent := GenerateMailContent(event)
	weixinContent := GenerateWeixinContent(event)

	if event.Priority() < 3 {
		redis.WriteSms(phones, smsContent)
	}

	redis.WriteMail(mails, smsContent, mailContent)
	redis.WriteWeixin(weixins, weixinContent)
}

// 低优先级的做报警合并
func consumeLowEvents(event *model.Event, action *api.Action) {
	if action.Uic == "" {
		return
	}

	if event.Priority() < 3 {
		ParseUserSms(event, action)
	}

	ParseUserMail(event, action)
	ParseUserWeixin(event, action)
}

func ParseUserSms(event *model.Event, action *api.Action) {
	userMap := api.GetUsers(action.Uic)

	content := GenerateSmsContent(event)
	metric := event.Metric()
	status := event.Status
	priority := event.Priority()

	queue := g.Config().Redis.UserSmsQueue

	rc := g.RedisConnPool.Get()
	defer rc.Close()

	for _, user := range userMap {
		dto := SmsDto{
			Priority: priority,
			Metric:   metric,
			Content:  content,
			Phone:    user.Phone,
			Status:   status,
		}
		bs, err := json.Marshal(dto)
		if err != nil {
			log.Println("json marshal SmsDto fail:", err)
			continue
		}

		_, err = rc.Do("LPUSH", queue, string(bs))
		if err != nil {
			log.Println("LPUSH redis", queue, "fail:", err, "dto:", string(bs))
		}
	}
}

func ParseUserWeixin(event *model.Event, action *api.Action) {
	userMap := api.GetUsers(action.Uic)

	content := GenerateWeixinContent(event)
	metric := event.Metric()
	status := event.Status
	priority := event.Priority()

	queue := g.Config().Redis.UserWeixinQueue

	rc := g.RedisConnPool.Get()
	defer rc.Close()

	for _, user := range userMap {
		dto := WeixinDto{
			Priority: priority,
			Metric:   metric,
			Content:  content,
			Im:       user.Im,
			Status:   status,
		}
		bs, err := json.Marshal(dto)
		if err != nil {
			log.Println("json marshal WeixinDto fail:", err)
			continue
		}

		_, err = rc.Do("LPUSH", queue, string(bs))
		if err != nil {
			log.Println("LPUSH redis", queue, "fail:", err, "dto:", string(bs))
		}
	}
}

func ParseUserMail(event *model.Event, action *api.Action) {
	userMap := api.GetUsers(action.Uic)

	metric := event.Metric()
	subject := GenerateSmsContent(event)
	content := GenerateMailContent(event)
	status := event.Status
	priority := event.Priority()

	queue := g.Config().Redis.UserMailQueue

	rc := g.RedisConnPool.Get()
	defer rc.Close()

	for _, user := range userMap {
		dto := MailDto{
			Priority: priority,
			Metric:   metric,
			Subject:  subject,
			Content:  content,
			Email:    user.Email,
			Status:   status,
		}
		bs, err := json.Marshal(dto)
		if err != nil {
			log.Println("json marshal MailDto fail:", err)
			continue
		}

		_, err = rc.Do("LPUSH", queue, string(bs))
		if err != nil {
			log.Println("LPUSH redis", queue, "fail:", err, "dto:", string(bs))
		}
	}
}
