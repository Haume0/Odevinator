package main

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

func SendVerifyMail(target string, code string, name string) {
	fmt.Println("Sending mail to "+target, " with code: "+code, " and name: "+name)
	from := "haumewastaken@gmail.com"
	password := "zduh hkqc cxwe zmjx"
	to := []string{target}
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	templateData, err := os.ReadFile("./verify_en.html")
	if err != nil {
		panic(err)
	}
	template := string(templateData)
	template = strings.ReplaceAll(template, "var_code", code)
	template = strings.ReplaceAll(template, "var_name", name)

	// Create MIME headers
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = strings.Join(to, ",")
	// headers["Subject"] = "Ödevini onaylamalısın!"
	headers["MIME-version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"UTF-8\""

	// Create the email body
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + template

	// Create authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Send actual message
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(message))
	if err != nil {
		panic(err)
	}
}
