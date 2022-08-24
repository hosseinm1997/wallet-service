package middlewares

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func HTTPReqPanicRecovery(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				fmt.Printf("%v", rvr)
				debug.PrintStack()
				// todo: logging
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
