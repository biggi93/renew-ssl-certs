package fileserver

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"simple-file-server/config"
	"strings"

	"github.com/go-chi/chi"
)

func GetCertbotHttpTestServer() *http.Server {
	config := config.Get()
	mux := chi.NewRouter()

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "/.well-known/acme-challenge"))
	getFileserverConf(mux, "/.well-known/acme-challenge", filesDir)

	mux.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Hetzner.LoadBalancerTargetPort),
		Handler: mux,
	}
}

func getFileserverConf(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

 