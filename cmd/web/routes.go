package main

import (
	"net/http"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// mux := http.NewServeMux()
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileserver := http.FileServer(neuterFileSystem{http.Dir("./ui/static/")})

	//1. serving static file
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileserver))
	// mux.Handle("/static/", http.StripPrefix("/static", fileserver))

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.login))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.loginPost))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.signup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.signupPost))

	protected := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.logout))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(router)
}

// to disable directory listing while serving static asset
// we can also serve empty html file for this or also use middleware but below approach is considered best
type neuterFileSystem struct {
	fs http.FileSystem
}

func (ns neuterFileSystem) Open(path string) (http.File, error) {
	f, err := ns.fs.Open(path)

	if err != nil {
		return nil, err
	}

	s, err := f.Stat()

	if s.IsDir() {
		index := filepath.Join(path, "index.html")

		if _, err := ns.fs.Open(index); err != nil {
			closeErr := f.Close()

			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}
