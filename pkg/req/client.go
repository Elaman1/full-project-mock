package req

import "net/http"

func GetClientMeta(r *http.Request) (ip string, userAgent string) {
	// Попытка получить из X-Forwarded-For (если за прокси)
	ip = r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}

	userAgent = r.Header.Get("User-Agent")
	return ip, userAgent
}
