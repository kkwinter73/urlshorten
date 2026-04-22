package main

import (
	"log"
	"net/http"
	"time"
)

// statusRecorder: ResponseWriter を包んで、書き込まれたステータスコードを記録する
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// loggingMiddleware: リクエストごとにメソッド・パス・ステータス・処理時間をログに出す
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rec := &statusRecorder{
			ResponseWriter: w,
			status:         http.StatusOK, // WriteHeader 呼ばれなければ 200 扱い
		}

		next.ServeHTTP(rec, r)

		log.Printf("%s %s %d %s",
			r.Method,
			r.URL.Path,
			rec.status,
			time.Since(start),
		)
	})
}
