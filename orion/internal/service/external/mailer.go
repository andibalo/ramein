package external

import (
	"context"
	"github.com/andibalo/ramein/orion/internal/config"
	sendinblue "github.com/sendinblue/APIv3-go-library/v2/lib"
	"go.uber.org/zap"
)

type Mailer interface {
	SendEmail(data SendEmailReq) error
}

type sendInBlueWrapper struct {
	cfg config.Config
	sib *sendinblue.APIClient
}

type SendEmailReq struct {
	SenderName     string
	SenderEmail    string
	Subject        string
	RecipientName  string
	RecipientEmail string
	HtmlContent    string
	TextContent    string
}

func NewSendInBlueService(appCfg config.Config) *sendInBlueWrapper {
	cfg := sendinblue.NewConfiguration()

	cfg.AddDefaultHeader("api-key", appCfg.SendInBlueApiKey())

	cfg.AddDefaultHeader("partner-key", appCfg.SendInBlueApiKey())

	sib := sendinblue.NewAPIClient(cfg)

	return &sendInBlueWrapper{
		sib: sib,
		cfg: appCfg,
	}
}

func (s *sendInBlueWrapper) SendEmail(data SendEmailReq) error {

	sender := &sendinblue.SendSmtpEmailSender{
		Name:  s.cfg.DefaultSenderName(),
		Email: s.cfg.DefaultSenderEmail(),
	}

	if data.SenderEmail != "" {
		sender.Email = data.SenderEmail
	}

	if data.SenderName != "" {
		sender.Name = data.SenderName
	}

	req := sendinblue.SendSmtpEmail{
		Sender: sender,
		To: []sendinblue.SendSmtpEmailTo{
			{
				Email: data.RecipientEmail,
				Name:  data.RecipientName,
			},
		},
		Subject: data.Subject,
	}

	if data.HtmlContent != "" {
		req.HtmlContent = data.HtmlContent
	} else {
		req.TextContent = data.TextContent
	}

	_, _, err := s.sib.TransactionalEmailsApi.SendTransacEmail(context.Background(), req)

	if err != nil {
		s.cfg.Logger().Error("[SendInBlue] Error sending email", zap.String("recipient_email", data.RecipientEmail), zap.Error(err))
		return err
	}

	return nil
}
