package main

import (
	"crypto/tls"
	"time"

	log "github.com/sirupsen/logrus"
	mail "github.com/xhit/go-simple-mail/v2"
)

var AuthStringToType map[string]mail.AuthType = map[string]mail.AuthType{
	"none":    mail.AuthNone,
	"crammd5": mail.AuthCRAMMD5,
	"login":   mail.AuthLogin,
	"plain":   mail.AuthPlain,
}

var EncryptionStringToType map[string]mail.Encryption = map[string]mail.Encryption{
	"ssltls":   mail.EncryptionSSLTLS,
	"starttls": mail.EncryptionSTARTTLS,
	"none":     mail.EncryptionNone,
}

type EmailSettings struct {
	SMTPHost           string   `yaml:"smtp_host"`
	SMTPPort           int      `yaml:"smtp_port"`
	SMTPUsername       string   `yaml:"smtp_username"`
	SMTPPassword       string   `yaml:"smtp_password"`
	SenderEmail        string   `yaml:"sender_email"`
	SenderName         string   `yaml:"sender_name"`
	ReceiverEmails     []string `yaml:"receiver_emails"`
	EmailSubject       string   `yaml:"email_subject"`
	EncryptionType     string   `yaml:"encryption_type"`
	AuthType           string   `yaml:"auth_type"`
	InsecureSkipVerify bool     `yaml:"tls_skip_verify"`
}

type Email struct {
	Server   *mail.SMTPServer `yaml:"-"`
	Client   *mail.SMTPClient `yaml:"-"`
	Settings *EmailSettings   `yaml:"settings,omitempty"`
}

func NewEmailClient(c *CslsConfig) *Email {
	if c.EmailClient == nil {
		return nil
	}
	S := Email{
		Server: &mail.SMTPServer{
			Host:           c.EmailSettings.SMTPHost,
			Port:           c.EmailSettings.SMTPPort,
			Username:       c.EmailSettings.SMTPUsername,
			Password:       c.EmailSettings.SMTPPassword,
			Encryption:     EncryptionStringToType[c.EmailSettings.EncryptionType],
			KeepAlive:      false,
			ConnectTimeout: 10 * time.Second,
			SendTimeout:    10 * time.Second,
			TLSConfig:      &tls.Config{InsecureSkipVerify: c.EmailSettings.InsecureSkipVerify},
			Authentication: AuthStringToType[c.EmailSettings.AuthType],
		},
		Settings: c.EmailSettings,
	}
	smtpC, err := S.Server.Connect()
	if err != nil {
		log.Debugf("Error connecting to smtp %v", err.Error())
	}
	S.Client = smtpC
	return &S
}
