package config

import "net/http"

func SetCookieSession(w http.ResponseWriter, token string) error {
	cookie := http.Cookie{
		Name:     "JWTToken",
		Value:    token,
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   SecureCookie,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)

	return nil
}

func DeleteCookieSession(w http.ResponseWriter) {
	newCookie := http.Cookie{
		Name:     "JWTToken",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   SecureCookie,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &newCookie)
}