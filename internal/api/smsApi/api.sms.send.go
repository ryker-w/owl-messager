package smsApi

import (
	"encoding/json"
	"github.com/kataras/iris/v12"
	"github.com/lishimeng/app-starter"
	"github.com/lishimeng/go-log"
	"github.com/lishimeng/owl/internal/api/common"
	"github.com/lishimeng/owl/internal/db/repo"
	"github.com/lishimeng/owl/internal/db/service"
)

type Req struct {
	Template      string      `json:"template,omitempty"` // template of this mail
	TemplateParam interface{} `json:"params,omitempty"`   // template params
	Sender        string      `json:"sender,omitempty"`   // mail send account on the platform
	Receiver      string      `json:"receiver,omitempty"` // receiver list(with comma if multi)
	Cc            string      `json:"cc,omitempty"`       // cc list(with comma if multi)
}

type Resp struct {
	app.Response
	MessageId int `json:"messageId,omitempty"`
}

// SendMail
/**
@Summary send a email

@Router /api/send/mail [post]

@Example
http://localhost/api/send/mail

{
	"template":"b7411049bbfe8068",
	"params":{"content":"O35A0001"},
	"subject":"电量低超提醒",
	"sender":"e949ae24481a9527",
	"receiver":"xxxx@qq.com"
}

*/
func SendMail(ctx iris.Context) {
	log.Debug("send mail api")
	var req Req
	var resp Resp
	err := ctx.ReadJSON(&req)
	if err != nil {
		log.Info("read req fail")
		log.Info(err)
		resp.Code = -1
		resp.Message = "req error"
		common.ResponseJSON(ctx, resp)
		return
	}

	// check params
	log.Debug("check params")

	if len(req.Sender) == 0 {
		log.Debug("param sender code nil")
		resp.Code = -1
		resp.Message = "sender nil"
		common.ResponseJSON(ctx, resp)
		return
	}
	sender, err := repo.GetMailSenderByCode(req.Sender)
	if err != nil {
		log.Debug("param sender not exist")
		resp.Code = -1
		resp.Message = "sender not exist"
		common.ResponseJSON(ctx, resp)
		return
	}

	if len(req.Template) == 0 {
		log.Debug("param template code nil")
		resp.Code = -1
		resp.Message = "template nil"
		common.ResponseJSON(ctx, resp)
		return
	}
	tpl, err := repo.GetMailTemplateByCode(req.Template)
	if err != nil {
		log.Debug("param template not exist")
		resp.Code = -1
		resp.Message = "template not exist"
		common.ResponseJSON(ctx, resp)
		return
	}

	var templateParams string
	switch req.TemplateParam.(type) {
	case string:
		templateParams = (req.TemplateParam).(string)
	default:
		bs, e := json.Marshal(req.TemplateParam)
		if e == nil {
			templateParams = string(bs)
		}
	}

	m, err := service.CreateMailMessage(
		sender,
		tpl,
		templateParams,
		req.Receiver, req.Cc)
	if err != nil {
		log.Info("can't create mail")
		log.Info(err)
		resp.Code = -1
		resp.Message = "create message failed"
		common.ResponseJSON(ctx, resp)
		return
	}

	log.Debug("create message success, id:%d", m.Id)
	resp.MessageId = m.Id

	resp.Code = common.RespCodeSuccess
	common.ResponseJSON(ctx, resp)
}
