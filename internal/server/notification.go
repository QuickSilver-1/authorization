package server

import (
	e "auth/internal/errors"
	"auth/internal/logger"
	"fmt"
	"net/smtp"
)

type NotificationService struct {
    smtpHost string
    smtpPort int 
    username string
    password string
}

func NewNotificationService(smtpHost string, smtpPort int, username, password string) *NotificationService {
    return &NotificationService{
        smtpHost: smtpHost,
        smtpPort: smtpPort,
        username: username,
        password: password,
    }
}

// Note уведомляет о входе с нового ip
func (n *NotificationService) Note(email, ip string, out chan error) {
	message := []byte("Вход с нового места\nВ ваш аккаунт был произведен вход с нового адреса, если это были не вы, то срочно поменяйте пароль")

	// Авторизация
	auth := smtp.PlainAuth("", n.username, n.password, n.smtpHost)

	// Отправка почты
	err := smtp.SendMail(fmt.Sprintf("%s:%d", n.smtpHost, n.smtpPort), auth, n.username, []string{email}, message)
	
	if err != nil {
		logger.Log.Error("Message sendding error")
		out<- &e.EmailError{
			Err: fmt.Sprintf("Message sendding error: %v", err),
		}
	}

	logger.Log.Debug("Message was send")
	out<- nil
}