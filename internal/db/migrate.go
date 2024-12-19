package db

import (
	"auth/internal/logger"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// CreateSchema выполняет миграции базы данных для создания схемы
func CreateSchema() error {
    logger.Log.Debug("Начало миграции")

	db, err := NewDB()

	if err != nil {
		return err
	}

    // Создаем экземпляр драйвера для PostgreSQL
    driver, err := postgres.WithInstance(db.Conn, &postgres.Config{})
    if err != nil {
        logger.Log.Error(fmt.Sprintf("Ошибка создания драйвера PostgreSQL: %v", err))
		return err
	}

    // Создаем мигратор с указанным источником миграций и базой данных
    m, err := migrate.NewWithDatabaseInstance("file://../../internal/migrations", "postgres", driver)
    if err != nil {
		logger.Log.Error(fmt.Sprintf("Ошибка создания мигратора: %v", err))
		return err
    }

    // Применяем миграции к базе данных
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Log.Error(fmt.Sprintf("Ошибка применения миграций: %v", err))
		return err
    }
    
    logger.Log.Debug("Миграции успешно применены!")

	return nil
}