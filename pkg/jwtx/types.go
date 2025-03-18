package jwtx

import "net/http"

type Handler interface {
	ExtractToken(r *http.Request) string
	SetJWTToken(w http.ResponseWriter, r *http.Request, userId int64, ssid string) error
	ClearToken(w http.ResponseWriter, r *http.Request) error
	SetLoginToken(w http.ResponseWriter, r *http.Request, userId int64) error
	CheckSession(r *http.Request, ssid string) error
}
