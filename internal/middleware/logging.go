package middleware

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"swagtask/internal/utils"
	"time"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}
func (w *wrappedWriter) Write(b []byte) (int, error) {
	i, err := w.ResponseWriter.Write(b)
	return i, err
}
func (w *wrappedWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := w.ResponseWriter.(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, fmt.Errorf("wrappedWriter: underlying ResponseWriter doesn't support hijacking")
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		showHttpDumps := false
		// ANSI escape code for bold text
		bold := "\033[1m"
		reset := "\033[0m"
		start := time.Now()
		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		dumpReq, errReq := httputil.DumpRequest(r, false)
		if errReq != nil {
			utils.LogError("Failed to dump request body:", errReq)
		}
		next.ServeHTTP(wrapped, r)

		if showHttpDumps {
			log.Println(
				"-------------------------------------------------------\n",
				bold+"HTTP 1.1 REQUEST representation:"+reset+"\n",
				indent(string(dumpReq), "\t"), "\n",

				bold+"RESPONSE:"+reset, "status", wrapped.statusCode, "|", time.Since(start),
				"\n----------------------------------------------------------------------------",
			)
		} else {
			log.Println(
				"-------------------------------------------------------\n",
				"request:", r.Method, r.URL.Path, "\n",
				"response:", wrapped.statusCode, time.Since(start),
				"\n----------------------------------------------------------------------------",
			)
		}
	})
}

func indent(text string, prefix string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = prefix + line
	}
	return strings.Join(lines, "\n")
}
