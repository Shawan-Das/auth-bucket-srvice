package service

import (
	"fmt"
	"github.com/jordan-wright/email"
	"log"
	"net/smtp"

	"github.com/spf13/viper"
	gomail "gopkg.in/gomail.v2"
)

type SmtpService struct{}

const (
	EMAIL_DESIGN_HTML = `<head>
    <style>
    body {
            font-family: Arial, sans-serif;
            background-color: #f5f5f5;
            margin: 0;
            padding: 0;
        }
    .container {
        max-width: 600px;
        margin: 0 auto;
        padding: 20px;
        background-color: #ffffff;
        box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
    }
    .header {
        text-align: center;
        margin-bottom: 20px;
    }
    .logo {
        max-width: 100px;
        height: auto;
    }
    .content {
        font-size: 16px;
        line-height: 1.5;
    }
    .otp {
        font-size: 24px;
        font-weight: bold;
        color: #267D54;  /* updated color */
        margin-top: 10px;
        letter-spacing: 0.1rem;
    }
    .link {
        display: inline-flex;
        justify-content: center;
        align-items: center;
        gap: 8px;
        margin-top: 20px;
        padding: 10px 16px;
        border-radius: 12px;
        border: 1px solid #EFD3A6;
        background: #FDF9F4;
        box-shadow: 0 1px 1px 0 rgba(0,0,0,0.08);
        cursor: pointer;
        color: #000 !important;
        text-decoration: none;
    }
    .link:hover {
        color: #D89121 !important;
        text-decoration: underline;
    }
    .footer {
        margin-top: 20px;
        text-align: center;
        color: #999;
    }
    </style>
</head>`
)

type CustomEmail struct {
	Username string `json:"username"`
	Subject  string `json:"subject"`
	Body     string `json:"body"`
}

var customEmail CustomEmail

// SendEmail sends an email
func (s *SmtpService) SendEmail(input CustomEmail) error {

	newComposeEmail := gomail.NewMessage()

	smtpHost := viper.GetString("smtp_host")
	senderEmail := viper.GetString("senderEmail")
	senderEmailAppPass := viper.GetString("password")
	smtpPort := viper.GetInt("smtp_port")
	newComposeEmail.SetHeader("From", senderEmail)
	newComposeEmail.SetHeader("To", input.Username)
	newComposeEmail.SetHeader("Subject", input.Subject)
	newComposeEmail.SetBody("text/html", input.Body)

	a := gomail.NewDialer(smtpHost, smtpPort, senderEmail, senderEmailAppPass)

	fmt.Println("smtpHost", smtpHost)
	if err := a.DialAndSend(newComposeEmail); err != nil {
		log.Printf("Failed to send email to %s: %v", input.Username, err)
		return err
	}
	return nil
}

func (s *SmtpService) SetNewPassWordMail(username string, empname string, tempPass string) error {
	rawUrl := viper.GetViper().GetStringMapString("url")["uiurl"] + "/#/login"
	customEmail.Username = username
	customEmail.Subject = "Temporary Password from HR Management System"
	customEmail.Body = `
	<!DOCTYPE html>
	<html>
	` + EMAIL_DESIGN_HTML + `
	<body>
		<div class="container">
			<div class="content">
				<p>Hello ` + empname + `,</p>
				<p>Your Temporary Password is: <span class="otp">` + tempPass + `</span></p>
				<p>To set your new password in HR Management Module, please click the following button:</p>
				<p><a class="link" href="` + rawUrl + `" target="_blank">Redirect Link</a></p>
			</div>
			<div class="footer">
			<p>This email has sent by  <span style="color:black">system administrator.</span></p>
			</div>
		</div>
	</body>
	</html>
	`

	emailSendError := s.SendEmail(customEmail)
	if emailSendError != nil {
		log.Println("Error sending email:", emailSendError)
		return emailSendError
	}

	return nil
}

// TODO: Version 2 of mail service
type EmailService struct{}

const (
	smtpAuthAddress   = "smtp.gmail.com"
	smtpServerAddress = "smtp.gmail.com:587"
)

type GmailSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

type EmailSender interface {
	SendEmail(
		subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
		attachFiles []string,
	) error
}

func (d *GmailSender) Init(configBytes []byte) error {
	d.name = ""
	d.fromEmailAddress = ""
	d.fromEmailPassword = ""

	return nil
}

func NewGmailSender(name string, fromEmailAddress string, fromEmailPassword string) EmailSender {
	return &GmailSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
	}
}

func (sender *GmailSender) SendEmail(
	subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
	attachFiles []string,
) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddress)
	e.Subject = subject
	e.HTML = []byte(content)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc

	for _, f := range attachFiles {
		_, err := e.AttachFile(f)
		if err != nil {
			return fmt.Errorf("failed to attach file %s: %w", f, err)
		}
	}

	smtpAuth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, smtpAuthAddress)
	return e.Send(smtpServerAddress, smtpAuth)
}

func SendAccountOpeningEmail(userName, userEmail, empIdentityNummer string, tempPass string) error {
	fmt.Printf("Inside Send email... Strat|| userName:%s , User:%s", userName, userEmail)
	rawUrl := viper.GetViper().GetStringMapString("url")["uiurl"] + "/#/login"
	senderName := viper.GetString("senderName")
	senderEmail := viper.GetString("senderEmail")
	password := viper.GetString("password")
	sender := NewGmailSender(senderName, senderEmail, password)

	subject := "Temporary Password from HR Management System"
	// content := `
	// <h3>Hello ` + userId + `,</h3>
	// <p><a class="link" href="` + rawUrl + `" target="_blank">Redirect Link</a></p>
	// `
	content := `
	<!DOCTYPE html>
	<html>
	` + EMAIL_DESIGN_HTML + `
	<body>
		<div class="container">
			<div class="content">
				<p>Hello ` + userName + `,</p>
				<p>Your Temporary Password is: <span class="otp">` + tempPass + `</span></p>
				<p>To set your new password in HR Management Module, please click the following button:</p>
				<p><a class="link" href="` + rawUrl + `" target="_blank">Redirect Link</a></p>
			</div>
			<div class="footer">
			<p>This email has sent by  <span style="color:black">system administrator.</span></p>
			</div>
		</div>
	</body>
	</html>
	`

	fmt.Println("Send email url content::::::::::::::::", content)
	to := []string{userEmail}
	//attachFiles := []string{"../README.md"}

	err := sender.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}

	fmt.Printf("Send email... End User:%s", userEmail)
	return nil
}
