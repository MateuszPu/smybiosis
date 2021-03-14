package mail

import (
	"bytes"
	"fmt"
	"github.com/jordan-wright/email"
	"html/template"
	"net/smtp"
)

type Service struct {
	Host     string
	Port     string
	Email    string
	Password string
}

type PaymentLinkData struct {
	Link string
	Name string
}

func (service *Service) SendEmailWithPaymentLink(to string, link string) {
	go func(to string, link string) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in f", r)
			}
		}()
		t := template.Must(template.ParseFiles("templates/email/payment-link.html"))
		var tpl bytes.Buffer
		t.Execute(&tpl, PaymentLinkData{Link: link})

		e := email.NewEmail()
		e.From = fmt.Sprintf("No-Reply <%s>", service.Email)
		e.To = []string{to}
		e.Subject = "Awesome Subject"
		e.HTML = tpl.Bytes()
		err := e.Send(fmt.Sprintf("%s:%s", service.Host, service.Port), smtp.PlainAuth("", service.Email, service.Password, service.Host))
		//err := e.SendWithTLS(fmt.Sprintf("%s:%s", service.Host, service.Port),
		//	smtp.PlainAuth("", service.Email, service.Password, service.Host),
		//	&tls.Config{ServerName: service.Host})
		println(err) //TODO: think about error and retry?
	}(to, link)
}
