package main

import (
	"fmt"
	"net/http"
)

// Sets secure headers, per OWASP guidelines.
// https://owasp.org/www-project-secure-headers/
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

// Middleware to logs each HTTP request, including the requests IP, protocol,
// method, and URI.
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip       = r.RemoteAddr
			protocol = r.Proto
			method   = r.Method
			uri      = r.URL.RequestURI()
		)

		app.logger.Info("received request", "ip", ip, "protocol", protocol, "method", method, "uri", uri)

		next.ServeHTTP(w, r)
	})
}

/*
Middleware to recover from panics and return a 500 server error. This should probably be the first middleware in the chain.

Note that this middleware will only have effect within a given go routine. So if a separate goroutine is initiated, you should include code to recover from panics inside that goroutine.

# Example

	go func() {
		defer func() {
			if err := recover(); err != nil {
				app.logger.Error(fmt.Sprint(err))
			}
		}()
		doSomeBackgroundProcessing()
	}()
*/
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/* Deferred functions are always run after a panic, as Go "unwinds" the handler stack. */
		defer func() {
			// If panicing, recover() restores normal execution and returns the error.
			if err := recover(); err != nil {
				// Set a header that will automatically close the HTTP connection after
				// response is set.
				w.Header().Set("Connection", "close")

				// Return type of recover() is any, so we need to format it as an error.
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
