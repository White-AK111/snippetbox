package main

import (
	"net/http"
)

func (a *app) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", a.home)
	mux.HandleFunc("/snippet", a.showSnippet)
	mux.HandleFunc("/snippet/create", a.createSnippet)

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("../../ui/static/")})
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	handler := a.Logging(mux)
	handler = a.PanicRecovery(handler)

	return handler
}

// neuteredFileSystem for custom file system
type neuteredFileSystem struct {
	fs http.FileSystem
}

// Open static file
func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() && err != nil {
		index := "index.html"
		if _, err := nfs.fs.Open(index); err == nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
		} else {
			return nil, err
		}
	}

	return f, nil
}
