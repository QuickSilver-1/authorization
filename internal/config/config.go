package config

import (
	e "auth/internal/errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var AppConfig *Config

// Config представляет конфигурацию приложения
type Config struct {
	HttpPort  int    // Порт для запуска приложения на HTTP
	PgHost    string // Хост для базы данных PostgreSQL
	PgPort    int    // Порт для базы данных PostgreSQL
	PgName    string // Имя базы данных PostgreSQL
	PgUser    string // Имя пользователя для базы данных PostgreSQL
	PgPass    string // Пароль для базы данных PostgreSQL
	SecretKey string // Секретный код для JWT
	SmtpHost  string // Хост для отправки сообщений через почтовый сервис
	SmtpPort  int    // Порт для отправки сообщений через почтовый сервис
	MailUser  string // Корпоративная почта
	MailPass  string // Пароль к корпоративной почте
}

func NewConfig() error {
	// Загрузка переменных окружения из .env файла
	err := godotenv.Load("config.env")

	if err != nil {
		return &e.ConfigError{
			Err : fmt.Sprintf(`Error loading configuration file: %v`, err),
		}
	}

	x := os.Getenv("DB_PORT")
	port, err := strconv.Atoi(x)

	if err != nil {
		return &e.ConfigError{
			Err: `Invalid "DB_PORT"`,
		}
	}

	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))

	if err != nil {
		return &e.ConfigError{
			Err: `Invalid "SMTP_PORT"`,
		}
	}

	httpPort, err := strconv.Atoi(os.Getenv("HTTP_PORT"))

	if err != nil {
		return &e.ConfigError{
			Err: `Invalid "HTTP_PORT"`,
		}
	}

	config := &Config{
		HttpPort:  httpPort,
		PgHost:    os.Getenv("DB_HOST"),
		PgPort:    port,
		PgName:    os.Getenv("DB_NAME"),
		PgUser:    os.Getenv("DB_USER"),
		PgPass:    os.Getenv("DB_PASSWORD"),
		SecretKey: os.Getenv("SECRET_KEY"),
		SmtpHost:  os.Getenv("SMTP_HOST"),
		SmtpPort:  smtpPort,
		MailUser:  os.Getenv("MAIL_USER"),
		MailPass:  os.Getenv("MAIL_PASS"),
	}

	AppConfig = config

	return nil
}
