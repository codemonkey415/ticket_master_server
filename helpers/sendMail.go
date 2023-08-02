package helper

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"
)

type Mail struct {
	Sender  string
	To      []string
	Subject string
	Body    string
}

func SendResetPasswordLink(email string, token string) (bool, error) {

	sender := "foryoursuccess0415@gmail.com"

	to := []string{
		email,
	}

	user := "foryoursuccess0415@gmail.com"
	password := "ukubrwzsfcxwjoih"

	subject := "Reset Password"
	body := `Click on the following link to reset your password: <a href="https://smt-client-mui.vercel.app/reset-password/"` + token + `">https://smt-client-mui.vercel.app/reset-password/` + token + `</a>`

	request := Mail{
		Sender:  sender,
		To:      to,
		Subject: subject,
		Body:    body,
	}

	addr := "smtp.gmail.com:587"
	host := "smtp.gmail.com"

	msg := BuildMessage(request)
	auth := smtp.PlainAuth("", user, password, host)
	err := smtp.SendMail(addr, auth, sender, to, []byte(msg))

	fmt.Println(body)

	if err != nil {
		fmt.Println(err)
		return false, err
	}

	return true, nil
}

func SendMail(email string, data string) {

	sender := "foryoursuccess0415@gmail.com"

	to := []string{
		email,
	}

	user := "foryoursuccess0415@gmail.com"
	password := "ukubrwzsfcxwjoih"

	subject := "Simple HTML mail"
	body := `<pre>` + data + `</pre>`

	request := Mail{
		Sender:  sender,
		To:      to,
		Subject: subject,
		Body:    body,
	}

	addr := "smtp.gmail.com:587"
	// host := "smtp.mailtrap.io"
	host := "smtp.gmail.com"

	msg := BuildMessage(request)
	auth := smtp.PlainAuth("", user, password, host)
	err := smtp.SendMail(addr, auth, sender, to, []byte(msg))

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Email sent successfully")
}

func BuildMessage(mail Mail) string {
	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg += fmt.Sprintf("From: %s\r\n", mail.Sender)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";"))
	msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	msg += fmt.Sprintf("\r\n%s\r\n", mail.Body)

	return msg
}
