package main

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

var first bool = true

func SendVerifyMail(target string, id, code string, name string) {
	// for the offline situations

	fmt.Println("")

	if first {

		fmt.Println(`
	İnternet bağlantısı yoksa kodu öğrenciye siz söylemelisiniz.`)

		first = false
	}

	fmt.Printf(`
	💯 %v • %v: ✅ %v ✅`, id, name, code)
	//smpt values
	from := MAIL
	password := PASS
	to := []string{target}
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	//the content of the email
	templateData, err := os.ReadFile("./verify_en.html")
	if err != nil {
		panic(err)
	}
	template := string(templateData)
	template = strings.ReplaceAll(template, "var_code", code)
	template = strings.ReplaceAll(template, "var_name", name)

	// Creating MIME headers
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = strings.Join(to, ",")
	// headers["Subject"] = "Ödevini onaylamalısın!"
	headers["MIME-version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"UTF-8\""

	// Creating the email body
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + template

	// Creating for authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending the email
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(message))
	if err != nil {
		panic(err)
	}
}
