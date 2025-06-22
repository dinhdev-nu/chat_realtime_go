package config

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/dinhdev-nu/realtime_auth_go/global"
	"github.com/jordan-wright/email"
)

func SendOTPEmailByText(toEmail, otp string) error {
	from := global.Config.Mail.From
	password := global.Config.Mail.Password
	host := global.Config.Mail.Host
	port := global.Config.Mail.Port

	subject := "Subject: Your OTP Code\n"
	body := fmt.Sprintf("Your OTP code is: %s", otp)
	msg := []byte(subject + "\n" + body)

	addr := fmt.Sprintf("%s:%s", host, port)
	auth := smtp.PlainAuth("", from, password, host)
	err := smtp.SendMail(addr, auth, from, []string{toEmail}, msg)

	if err != nil {
		fmt.Println("::::::::::::: Failed to send email:", err)
		return err
	}

	fmt.Println("::::::::::::: Email sent successfully")
	return nil
}

// SendOTPEmailByTemplate sends an OTP email using a templatep

func RenderHTMLTTempate(filepath string, data map[string]interface{}) (string, error) {
	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, data)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
func SendOTPEmailByTemplate(toEmail, otp string) error {
	from := global.Config.Mail.From
	password := global.Config.Mail.Password
	host := global.Config.Mail.Host
	port := global.Config.Mail.Port

	subject := "Subject: Your OTP Code\n"

	// Load the HTML template

	data := map[string]interface{}{
		"otp": otp,
	}
	httpBody, err := RenderHTMLTTempate("template-email/otp-register-email.html", data)
	if err != nil {
		fmt.Println("::::::::::::: Failed to load template:", err)
		return err
	}

	mime := "MIME-version: 1.0;\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\n"

	msg := []byte(subject + mime + "\n" + httpBody)

	auth := smtp.PlainAuth("", from, password, host)
	addr := fmt.Sprintf("%s:%s", host, port)

	return smtp.SendMail(addr, auth, from, []string{toEmail}, msg)

}

// lib
func SendEmail(toEmail, otp string) error {
	data := map[string]interface{}{
		"otp": otp,
	}

	tpl, err := template.ParseFiles("template-email/otp-register-email.html")
	if err != nil {
		fmt.Println("::::::::::::: Failed to parse template:", err)
	}

	var buffer bytes.Buffer
	err = tpl.Execute(&buffer, data)
	if err != nil {
		fmt.Println("::::::::::::: Failed to execute template:", err)
	}

	// create the email
	e := email.NewEmail()
	e.From = global.Config.Mail.From
	e.To = []string{toEmail}
	e.Subject = "Your OTP Code"
	e.HTML = buffer.Bytes()

	err = e.Send(
		fmt.Sprintf("%s:%s", global.Config.Mail.Host, global.Config.Mail.Port),
		smtp.PlainAuth(
			"",
			global.Config.Mail.From,
			global.Config.Mail.Password,
			global.Config.Mail.Host,
		),
	)
	if err != nil {
		fmt.Println("::::::::::::: Failed to send email:", err)
		return err
	}

	fmt.Println("::::::::::::: Email sent successfully")
	return nil
}
