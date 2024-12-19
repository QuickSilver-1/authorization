package server

import (
	"auth/internal/config"
	"auth/internal/db"
	"auth/internal/logger"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Answer struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	if idStr == "" {
		http.Error(w, "Empty id", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ip := strings.Split(r.RemoteAddr, ":")[0]

	DB, err := db.NewDB()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	generateTokens(w, ip, id, DB)
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	if idStr == "" {
		http.Error(w, "Empty id", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ip := strings.Split(r.RemoteAddr, ":")[0]

	DB, err := db.NewDB()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out := make(chan string)
	go DB.DelRefresh(id, ip, out)

	ipOut := make(chan string)
	ipErrorchan := make(chan error)
	doneChan := make(chan bool)
	go func() {
		exist := false

		for result := range ipOut {

			if ip == result {
				exist = true
				break
			}
		}

		if !exist {
			doneChan <- true
		}

		doneChan <- false
	}()

	go DB.GetIp(id, ipOut, ipErrorchan)

	cookie, err := r.Cookie("RefreshToken")

	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			http.Error(w, "Token not present. Need authorisation", http.StatusUnauthorized)
			return
		}

		http.Error(w, "Bad cookie", http.StatusBadRequest)
		return
	}

	refresh, err := base64.URLEncoding.DecodeString(cookie.Value)

	if err != nil {
		http.Error(w, "Base64 decoding error", http.StatusInternalServerError)
		return
	}

	existRefresh := <-out

	if existRefresh == "-1" {
		http.Error(w, "Refresh token expires. Need autorisation", http.StatusUnauthorized)
		return
	}

	if existRefresh == "0" {
		http.Error(w, "Refresh for this access not exist. Need authorisation", http.StatusBadRequest)
		return
	}

	err = <-ipErrorchan

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	needAuth := <-doneChan

	if needAuth {
		http.Error(w, "New ip. Need authorisation", http.StatusUnauthorized)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(existRefresh), refresh) != nil {
		http.Error(w, "Invalid refresh token", http.StatusBadRequest)
		return
	}

	generateTokens(w, ip, id, DB)
}

// AnswerHandler обрабатывает и отправляет ответ клиенту
func answerHandler(w http.ResponseWriter, code int, value *Answer) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(value)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Log.Error(fmt.Sprintf("Ошибка кодирования ответа: %v", err))
	}
}

func generateTokens(w http.ResponseWriter, ip string, id int, DB *db.ConnectDatabase) {
	defer DB.Conn.Close()

	out := make(chan string)
	ipErrchan := make(chan error)

	var settingErrchan chan error
	go func() {
		exist := false

		for result := range out {
			if ip == result {
				exist = true
				break
			}
		}

		if !exist {
			settingErrchan = make(chan error)
			go DB.SetIp(id, ip, settingErrchan)
		}
	}()

	go DB.GetIp(id, out, ipErrchan)

	token := NewToken(id, ip)
	refresh, err := token.CreateRefresh()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	refreshBcrypt, err := bcrypt.GenerateFromPassword(refresh, bcrypt.DefaultCost)

	if err != nil {
		http.Error(w, "Bcrypt code generate error", http.StatusInternalServerError)
		return
	}

	refreshErrchan := make(chan error)
	go DB.SetRefresh(refreshBcrypt, id, ip, refreshErrchan)

	access, err := token.MakeJWT(config.AppConfig.SecretKey)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = <-refreshErrchan

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = <-ipErrchan

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = <-refreshErrchan

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	refreshString := base64.URLEncoding.EncodeToString(refresh)

	cookie := &http.Cookie{
		Name:     "RefreshToken",
		Value:    refreshString,
		HttpOnly: true,
		Path:     "/auth",
		Expires:  time.Now().AddDate(0, 0, 30),
	}

	http.SetCookie(w, cookie)

	answerHandler(w, http.StatusOK, &Answer{
		Access:  access,
		Refresh: refreshString,
	})
}
