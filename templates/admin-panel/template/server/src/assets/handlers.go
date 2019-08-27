package assets

import (
	"github.com/NYTimes/gziphandler"
	"net/http"
	"path/filepath"
	"strings"
)

type Handler struct{}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// strip leading "/"
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	// check whether a file exists at the given path
	cannonicalName := strings.Replace(path, "\\", "/", -1)
	if _, ok := _bindata[cannonicalName]; !ok {

		// file does not exist, serve index.html
		w.Header().Set("Content-Type", "text/html")
		f, err := Asset("index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(f)
	} else {

		// otherwise, use http.FileServer to serve the static dir
		http.FileServer(AssetFile()).ServeHTTP(w, r)
	}
}

func NewHandler() http.Handler {
	return gziphandler.GzipHandler(Handler{})
}
