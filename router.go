package silverback

import (
	"fmt"
	"net/http"
	"path"
	"strconv"

	"github.com/gorilla/mux"
)

// A Router is an extension of "github.com/gorilla/mux".Router.
type Router struct {
	mux.Router

	codecs []Codec
}

func NewRouter() *Router {
	return &Router{
		Router: *mux.NewRouter(),
	}
}

func (r *Router) setupIDPaths(handler Handler) {
	idRoutePath := path.Join(handler.Path(), "{id}")

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
			subRouter.Methods("GET").HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
				h := handler.New(req).(Getter)
				resp := idHandle(h, h.Get, mux.Vars(req)["id"])
				WriteResponse(writer, resp, r.codecs)
			})
			subRouter.Methods("HEAD").HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
				h := handler.New(req).(Getter)
				resp := idHandle(h, h.Get, mux.Vars(req)["id"])
				WriteHead(writer, resp, r.codecs)
			})
		}
		if hasPutter {
			subRouter.Methods("PUT").HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
				h := handler.New(req).(Putter)
				resp := idHandle(h, h.Put, mux.Vars(req)["id"])
				WriteResponse(writer, resp, r.codecs)
			})
		}
		if hasPatcher {
			subRouter.Methods("PATCH").HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
				h := handler.New(req).(Patcher)
				resp := idHandle(h, h.Patch, mux.Vars(req)["id"])
				WriteResponse(writer, resp, r.codecs)
			})
		}
		if hasDeleter {
			subRouter.Methods("DELETE").HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
				h := handler.New(req).(Deleter)
				resp := idHandle(h, h.Delete, mux.Vars(req)["id"])
				WriteResponse(writer, resp, r.codecs)
			})
		}
		route.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			writeAllowHeader(idAllowed, writer)
			writer.WriteHeader(http.StatusMethodNotAllowed)
		})
	}
}

func (r *Router) setupNonIDPaths(handler Handler) {
	routePath := handler.Path()

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
			subRouter.Methods("GET").HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
				h := handler.New(req).(Querier)
				resp := handle(h, h.Query)
				WriteResponse(writer, resp, r.codecs)
			})
			subRouter.Methods("HEAD").HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
				h := handler.New(req).(Querier)
				resp := handle(h, h.Query)
				WriteHead(writer, resp, r.codecs)
			})
		}
		if hasPoster {
			subRouter.Methods("POST").HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
				h := handler.New(req).(Poster)
				resp := handle(h, h.Post)
				WriteResponse(writer, resp, r.codecs)
			})
		}
		route.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			writeAllowHeader(allowed, writer)
			writer.WriteHeader(http.StatusMethodNotAllowed)
		})
	}
}

// AddCodec registers a codec with this router.  Any codecs added in
// this way will be supplied to any *Response value that has not had
// its codecs set (via NewResponseForCodecs).
func (r *Router) AddCodec(codec Codec) {
	r.codecs = append(r.codecs, codec)
}

// Route routes the methods on Handler to paths, based on the
// Handler's Path().
func (r *Router) Route(handler Handler) {
	r.setupIDPaths(handler)
	r.setupNonIDPaths(handler)
}

func optionsHandler(methods []string) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, r *http.Request) {
		writeAllowHeader(methods, writer)
		writer.WriteHeader(http.StatusNoContent)
	}
}

func writeAllowHeader(methods []string, writer http.ResponseWriter) {
	for _, method := range methods {
		writer.Header().Add("Allow", method)
	}
}

func handle(h Handler, f func() *Response) *Response {
	if before, ok := h.(BeforeHandler); ok {
		before.BeforeHandle()
	}
	resp := f()
	if after, ok := h.(AfterHandler); ok {
		after.AfterHandle(resp)
	}
	return resp
}

func idHandle(h Handler, f func(string) *Response, id string) *Response {
	if before, ok := h.(BeforeHandler); ok {
		before.BeforeHandle()
	}
	resp := f(id)
	if after, ok := h.(AfterHandler); ok {
		after.AfterHandle(resp)
	}
	return resp
}

func WriteHeaders(writer http.ResponseWriter, resp *Response) {
	for name, values := range resp.Headers {
		for _, v := range values {
			writer.Header().Add(name, v)
		}
	}
}

func WriteHead(writer http.ResponseWriter, resp *Response, codecs []Codec) (body []byte) {
	WriteHeaders(writer, resp)
	writer.WriteHeader(resp.Status)
	if resp.codecs == nil {
		resp.codecs = codecs
	}
	body, err := resp.Codec().Marshal(resp.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		msg := fmt.Sprintf("Error marshalling data: %v", err)
		writer.Write([]byte(msg))
		return
	}
	writer.Header().Set("Content-Length", strconv.Itoa(len(body)))
	if resp.Status == http.StatusOK {
		// For HEAD requests, the norm is StatusNoContent; however, we
		// don't want to overwrite error states, redirects, or other
		// statuses.  To avoid that, we only overwrite StatusOK.
		resp.Status = http.StatusNoContent
	}
	writer.WriteHeader(resp.Status)
	return body
}

func WriteResponse(writer http.ResponseWriter, resp *Response, codecs []Codec) {
	body := WriteHead(writer, resp, codecs)
	writer.Write(body)
}
