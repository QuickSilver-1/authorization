package db

import (
	"auth/internal/config"
	e "auth/internal/errors"
	"auth/internal/logger"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type ConnectDatabase struct {
	Conn    *sql.DB
}

// NewDB создаёт новое подключение к базе данных
func NewDB() (*ConnectDatabase, error) {
    logger.Log.Debug("Подключение к базе данных")
    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
    config.AppConfig.PgHost, config.AppConfig.PgPort, config.AppConfig.PgUser, config.AppConfig.PgPass, config.AppConfig.PgName)

    conn, err := sql.Open("postgres", psqlInfo)

    if err != nil {
        logger.Log.Error(fmt.Sprintf("Проблема с подключением к базе данных: %v", err))
        return nil, &e.DbConnectionError{
			Err: fmt.Sprintf("Проблема с подключением к базе данных: %s", err),
		}
    }

    // Настройки подключения
    conn.SetMaxOpenConns(10)
    conn.SetMaxIdleConns(5)
    conn.SetConnMaxLifetime(time.Second * 3)

    DB := &ConnectDatabase{
        Conn: conn,
    }

    return DB, nil
}

func (c *ConnectDatabase) SetRefresh(token []byte, id int, ip string, errchan chan error) {
    defer close(errchan)

    query, err := c.Conn.Begin()

    if err != nil {
        query.Rollback()
        errchan<- &e.DbQueryError{
            Err: fmt.Sprintf("Creating transaction error: %v", err),
        }

        return
    }

    _, err = query.Exec(` DELETE FROM tokens WHERE "user_id" = $1 `, id)

    if err != nil {
        query.Rollback()
        errchan<- &e.DbQueryError{
            Err: fmt.Sprintf("Deleting old refresh token error: %v", err),
        }

        return
    }

    _, err = query.Exec(` INSERT INTO tokens ("user_id", "value", "expires", "ip") VALUES ($1, $2, $3, $4) `, id, token, time.Now().AddDate(0, 0, 30), ip)

    if err != nil {
        query.Rollback()
        errchan<- &e.DbQueryError{
            Err: fmt.Sprintf("Setting refresh token error: %v", err),
        }

        return
    }

    err = query.Commit()


    if err != nil {
        query.Rollback()
        errchan<- &e.DbQueryError{
            Err: fmt.Sprintf("Commiting error: %v", err),
        }

        return
    }

    errchan<- nil
}

func (c *ConnectDatabase) DelRefresh(id int, ip string, out chan string) {
    defer close(out)

    var token string
    var expires time.Time
    c.Conn.QueryRow(` DELETE FROM tokens WHERE "user_id" = $1 AND "ip" = $2 RETURNING "value", "expires" `, id, ip).Scan(&token, &expires)

    if expires.IsZero() {
        out<- "0"

        return
    }

    if expires.Before(time.Now()) {
        out<- "-1"

        return
    }

    out<- token
}

func (c *ConnectDatabase) GetIp(id int, out chan string, errchan chan error) {
    defer close(errchan)
    defer close(out)

    rows, err := c.Conn.Query(` SELECT "value" FROM ips WHERE "user_id" = $1 `, id)

    if err != nil {
        errchan<- &e.DbQueryError{
            Err: fmt.Sprintf("Getting ip error: %v", err),
        }
        return
    }

    var existIp string
    if rows.Next() {
        rows.Scan(&existIp)
        out<- existIp
    }

    errchan <-nil
} 

func (c *ConnectDatabase) SetIp(id int, ip string, errchan chan error) {
    defer close(errchan)

    _, err := c.Conn.Exec(` INSERT INTO ips ("value", "user_id") VALUES ($1, $2) `, ip, id)

    if err != nil {
        errchan<- &e.DbQueryError{
            Err: fmt.Sprintf("Setting ip error: %v", err),
        }
        return
    }
}