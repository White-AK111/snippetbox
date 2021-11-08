package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

// serverError write error to errorLog
func (a *app) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	a.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError send error to client
func (a *app) clientError(w http.ResponseWriter, status uint) {
	http.Error(w, http.StatusText(int(status)), int(status))
}

// notFound for send 404
func (a *app) notFound(w http.ResponseWriter) {
	a.clientError(w, http.StatusNotFound)
}

// render template
func (a *app) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {

	ts, ok := a.templateCache[name]

	if !ok {
		a.serverError(w, fmt.Errorf("template %s not exist", name))
		return
	}

	buf := new(bytes.Buffer)

	err := ts.Execute(buf, a.addDefaultData(td, r))
	if err != nil {
		a.serverError(w, err)
		return
	}

	buf.WriteTo(w)
}

// addDefaultData add default data into template
func (a *app) addDefaultData(td *templateData, r *http.Request) *templateData {

	if td == nil {
		td = &templateData{}
	}

	td.CurrentYear = uint(time.Now().Year())
	//td.Flash = app.session.PopString(r, "flash")
	return td
}

// Logging is middleware for logging information about each request.
func (a *app) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, req)
		a.infoLog.Printf("%s %s %s", req.Method, req.RequestURI, time.Since(start))
	})
}

// PanicRecovery is middleware for recovering from panics in `next` and
// returning a StatusInternalServerError to the client.
func (a *app) PanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				a.errorLog.Println(string(debug.Stack()))
			}
		}()
		next.ServeHTTP(w, req)
	})
}
