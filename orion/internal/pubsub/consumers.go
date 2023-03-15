package pubsub

import (
	"bytes"
	"context"
	"encoding/json"
	pubsubCommons "github.com/andibalo/ramein/commons/pubsub"
	"github.com/andibalo/ramein/commons/rabbitmq"
	"github.com/andibalo/ramein/orion/internal/constants"
	"github.com/andibalo/ramein/orion/internal/service/external"
	"go.uber.org/zap"
	"html/template"
)

func (p *pubsub) CoreNewUserRegisteredHandler(c context.Context, message rabbitmq.Message) error {
	p.LogPayload(CORE_NEW_USER_REGISTERED, message)

	payload := message.Payload

	var data pubsubCommons.CoreNewRegisteredUserPayload

	jsonData, err := json.Marshal(payload)
	if err != nil {
		p.Config.Logger().Error("Error marshaling payload to json", zap.Error(err))
		return err
	}

	json.Unmarshal(jsonData, &data)

	tmpl, err := p.templateRepo.GetByTemplateName(constants.CORE_VERIFY_EMAIL_V1)

	if err != nil {
		p.Config.Logger().Error("Error getting email template", zap.Error(err))
		return err
	}

	t, err := template.New("verify_email").Parse(tmpl.Template)
	if err != nil {
		p.Config.Logger().Error("Error parsing template", zap.Error(err))
		return err
	}

	buf := new(bytes.Buffer)

	if err = t.Execute(buf, data); err != nil {
		p.Config.Logger().Error("Error binding data to template", zap.Error(err))
		return err
	}

	emailBody := buf.String()

	p.Config.Logger().Info("Email", zap.String("email", emailBody))

	sendEmailReq := external.SendEmailReq{
		Subject:        "Welcome to Ramein!",
		RecipientName:  data.FirstName + " " + data.LastName,
		RecipientEmail: data.Email,
		HtmlContent:    emailBody,
	}

	err = p.mailer.SendEmail(sendEmailReq)
	if err != nil {
		p.Config.Logger().Error("[CoreNewUserRegisteredHandler] Failed to send email", zap.Error(err))
		return err
	}

	return nil
}
