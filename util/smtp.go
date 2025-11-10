package util

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/smtp"

	"github.com/sirupsen/logrus"
)

var _smppLogger = logrus.New()

// SMTPUtil sends mail using SMTP protocol
type SMTPUtil struct {
	tlsConfig   *tls.Config
	auth        smtp.Auth
	host        string
	server      string
	fromEmail   string
	testEmailID string
}
type smtpConfig struct {
	Server    string `json:"server"`
	UID       string `json:"smtpUserId"`
	Secret    string `json:"smtpPassword"`
	TestEmail string `json:"testEmailID"`
}

// NewSMTPUtil returns an initialized version of SMTPUtil
func NewSMTPUtil(config []byte) *SMTPUtil {
	var smtpConf smtpConfig
	err := json.Unmarshal(config, &smtpConf)
	if err != nil {
		_smppLogger.Errorf("Error in SMTP configuration file ", err)
		return nil
	}
	smtpUtil := new(SMTPUtil)
	if smtpUtil.Init(smtpConf.Server, smtpConf.UID, smtpConf.Secret, smtpConf.TestEmail) != nil {
		_smppLogger.Errorf("Error in intialization")
		return nil
	}

	return smtpUtil
}

// Init initializes the util class
func (s *SMTPUtil) Init(server, uid, secret, testEmailID string) error {
	s.host, _, _ = net.SplitHostPort(server)
	s.server = server
	s.fromEmail = uid
	s.auth = smtp.PlainAuth("", uid, secret, s.host)
	// TLS config
	s.tlsConfig = &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         s.host,
	}
	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", s.server, s.tlsConfig)
	if err != nil {
		_smppLogger.Errorf("Unable to connect to %s %v", s.server, err)
		return err
	}

	client, err := smtp.NewClient(conn, s.host)
	if err != nil {
		_smppLogger.Errorf("Unable create client  %v", s.server, err)
		return err
	}
	if err = client.Auth(s.auth); err != nil {
		_smppLogger.Errorf("Authentication error %v  %v", s.auth, err)
		return err
	}
	client.Quit()
	if len(testEmailID) > 0 {
		s.testEmailID = testEmailID
	}
	_smppLogger.Info("Successfully authenticated SMPP server")
	return nil
}

// SendEmail sends an email
func (s *SMTPUtil) SendEmail(toAddress, subject string, mailBody []byte) error {
	// We assume init has been called already
	conn, err := tls.Dial("tcp", s.server, s.tlsConfig)
	if err != nil {
		_smppLogger.Errorf("Unable to connect to %s %v", s.server, err)
		return err
	}

	client, err := smtp.NewClient(conn, s.host)
	if err != nil {
		_smppLogger.Errorf("Unable create client  %v", s.server, err)
		return err
	}
	defer client.Quit()
	if err = client.Auth(s.auth); err != nil {
		_smppLogger.Errorf("Authentication error %v  %v", s.auth, err)
		return err
	}
	// To && From

	if err = client.Mail(s.fromEmail); err != nil {
		return err
	}
	actualToAddr := toAddress
	if len(s.testEmailID) > 0 {
		actualToAddr = s.testEmailID
	}
	if err = client.Rcpt(actualToAddr); err != nil {
		return err
	}
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("From:%s\r\n", s.fromEmail))
	buf.WriteString(fmt.Sprintf("To:%s\r\n", actualToAddr))
	buf.WriteString(fmt.Sprintf("Subject:%s\r\n", subject))
	buf.WriteString("\r\n")
	buf.Write(mailBody)

	// Data
	writer, err := client.Data()
	if err != nil {
		return err
	}
	defer writer.Close()
	_, err = writer.Write(buf.Bytes())
	if err != nil {
		return err
	}
	_smppLogger.Infof("Mail sent to %s successfully ", toAddress)
	return nil
}
