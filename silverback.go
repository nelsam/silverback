package silverback

import (
	"net/http"
	"path"

	"github.com/gorilla/mux"
)

// A Response is a container for typical response fields.
type Response struct {
	Status  int
	Headers http.Header
	Body    interface{}
}

// A Router is an extension of "github.com/gorilla/mux".Router.
type Router struct {
	mux.Router
}

// Route routes the methods on Handler to paths, based on the
// Handler's Path().
func (r *Router) Route(handler Handler) {
	routePath := handler.Path()
	idRoutePath := path.Join(routePath, "{id}")

	_, hasGetter := handler.(Getter)
	_, hasPutter := handler.(Putter)
	_, hasPatcher := handler.(Patcher)
	_, hasDeleter := handler.(Deleter)
	var idAllowed []string
	if hasGetter {
		idAllowed = append(idAllowed, "GET", "HEAD")
	}
	if hasPutter {
		idAllowed = append(idAllowed, "PUT")
	}
	if hasPatcher {
		idAllowed = append(idAllowed, "PATCH")
	}
	if hasDeleter {
		idAllowed = append(idAllowed, "DELETE")
	}
	if len(idAllowed) > 0 {
		idAllowed = append(idAllowed, "OPTIONS")
		route := r.Path(idRoutePath)
		subRouter := route.Subrouter()
		subRouter.Methods("OPTIONS").HandlerFunc(optionsHandler(idAllowed))
		if hasGetter {
			subRouter.Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				h := handler.New(r).(Getter)
				resp := idHandle(h, h.Get, mux.Vars(r)["id"])
				writeResponse(resp, w)
			})
			subRouter.Methods("HEAD").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				h := handler.New(r).(Getter)
				resp := idHandle(h, h.Get, mux.Vars(r)["id"])
				writeHead(resp, w)
			})
		}
		if hasPutter {
			subRouter.Methods("PUT").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				h := handler.New(r).(Putter)
				resp := idHandle(h, h.Put, mux.Vars(r)["id"])
				writeResponse(resp, w)
			})
		}
		if hasPatcher {
			subRouter.Methods("PATCH").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				h := handler.New(r).(Patcher)
				resp := idHandle(h, h.Patch, mux.Vars(r)["id"])
				writeResponse(resp, w)
			})
		}
		if hasDeleter {
			subRouter.Methods("DELETE").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				h := handler.New(r).(Deleter)
				resp := idHandle(h, h.Delete, mux.Vars(r)["id"])
				writeResponse(resp, w)
			})
		}
		route.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			writeAllowHeader(idAllowed, w)
			w.WriteHeader(http.StatusMethodNotAllowed)
		})
	}

	_, hasQuerier := handler.(Querier)
	_, hasPoster := handler.(Poster)
	var allowed []string
	if hasQuerier {
		allowed = append(allowed, "GET", "HEAD")
	}
	if hasPoster {
		allowed = append(allowed, "POST")
	}
	if len(allowed) > 0 {
		allowed = append(allowed, "OPTIONS")
		route := r.Path(routePath)
		subRouter := route.Subrouter()
		subRouter.Methods("OPTIONS").HandlerFunc(optionsHandler(allowed))
		if hasQuerier {
			subRouter.Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				h := handler.New(r).(Querier)
				resp := handle(h, h.Query)
				writeResponse(resp, w)
			})
			subRouter.Methods("HEAD").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				h := handler.New(r).(Querier)
				resp := handle(h, h.Query)
				writeHead(resp, w)
			})
		}
		if hasPoster {
			subRouter.Methods("POST").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				h := handler.New(r).(Poster)
				resp := handle(h, h.Post)
				writeResponse(resp, w)
			})
		}
		route.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, method := range allowed {
				w.Header().Add("Allow", method)
			}
			w.WriteHeader(http.StatusMethodNotAllowed)
		})
	}
}

func optionsHandler(methods []string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		writeAllowHeader(methods, w)
		w.WriteHeader(http.StatusNoContent)
	}
}

func writeAllowHeader(methods []string, w http.ResponseWriter) {
	for _, method := range methods {
		w.Header().Add("Allow", method)
	}
}

func handle(h Handler, f func() Response) Response {
	if before, ok := h.(BeforeHandler); ok {
		before.BeforeHandle()
	}
	resp := f()
	if after, ok := h.(AfterHandler); ok {
		after.AfterHandle(resp)
	}
	return resp
}

func idHandle(h Handler, f func(string) Response, id string) Response {
	if before, ok := h.(BeforeHandler); ok {
		before.BeforeHandle()
	}
	resp := f(id)
	if after, ok := h.(AfterHandler); ok {
		after.AfterHandle(resp)
	}
	return resp
}

func writeHeaders(resp Response, w http.ResponseWriter) {
	for name, values := range resp.Headers {
		for _, v := range values {
			w.Header().Add(name, v)
		}
	}
}

func writeHead(resp Response, w http.ResponseWriter) {
	writeHeaders(resp, w)
	if resp.Status < http.StatusBadRequest {
		resp.Status = http.StatusNoContent
	}
	// TODO: write Content-Length
	w.WriteHeader(resp.Status)
}

func writeResponse(resp Response, w http.ResponseWriter) {
	writeHeaders(resp, w)
	w.WriteHeader(resp.Status)
	// TODO: Write body.
}
