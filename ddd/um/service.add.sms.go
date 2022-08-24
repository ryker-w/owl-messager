package um

import (
	"github.com/lishimeng/owl/internal/db/model"
	"github.com/lishimeng/owl/internal/db/repo"
	"github.com/lishimeng/owl/internal/db/service"
)

func serviceAddSms(smsTemplate, tplParams, receiver string) (m model.MessageInfo, err error) {

	tpl, err := repo.GetSmsTemplateByCode(smsTemplate)
	if err != nil {
		return
	}
	m, err = service.CreateSmsMessage(
		tpl,
		tplParams,
		receiver,
	)
	return
}
