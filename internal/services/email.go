package services

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

func SendEmailWarning(userID, oldIP, newIP string) {
	from := os.Getenv("SMTP_EMAIL") 
	password := os.Getenv("SMTP_PASSWORD") 
	to := []string{"user-email@example.com"} 
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	subject := "IP Change Warning"
	body := fmt.Sprintf("Warning: user %s changed IP from %s to %s\n", userID, oldIP, newIP)
	msg := "From: " + from + "\n" +
		"To: " + to[0] + "\n" +
		"Subject: " + subject + "\n\n" +
		body
	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(msg))
	if err != nil {
		log.Printf("Error sending email: %v", err)
	} else {
		log.Printf("Email sent successfully to %s", to[0])
	}
}
