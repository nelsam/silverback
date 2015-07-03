package silverback

import (
	"fmt"
	"net/http"
	"path"
	"strconv"

	"github.com/gorilla/handlers"
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
	h := make(handlers.MethodHandler, 5)
	if _, hasGetter := handler.(Getter); hasGetter {
		h["GET"] = http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			h := handler.New(req).(Getter)
			resp := idHandle(h, h.Get, mux.Vars(req)["id"])
			WriteResponse(writer, resp, r.codecs)
		})
		h["HEAD"] = http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			h := handler.New(req).(Getter)
			resp := idHandle(h, h.Get, mux.Vars(req)["id"])
			WriteHead(writer, resp, r.codecs)
		})
	}
	if _, hasPutter := handler.(Putter); hasPutter {
		h["PUT"] = http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			h := handler.New(req).(Putter)
			resp := idHandle(h, h.Put, mux.Vars(req)["id"])
			WriteResponse(writer, resp, r.codecs)
		})
	}
	if _, hasPatcher := handler.(Patcher); hasPatcher {
		h["PATCH"] = http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			h := handler.New(req).(Patcher)
			resp := idHandle(h, h.Patch, mux.Vars(req)["id"])
			WriteResponse(writer, resp, r.codecs)
		})
	}
	if _, hasDeleter := handler.(Deleter); hasDeleter {
		h["DELETE"] = http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			h := handler.New(req).(Deleter)
			resp := idHandle(h, h.Delete, mux.Vars(req)["id"])
			WriteResponse(writer, resp, r.codecs)
		})
	}
	if len(h) > 0 {
		idRoutePath := path.Join(handler.Path(), "{id}")
		r.Path(idRoutePath).Handler(h)
	}
}

func (r *Router) setupNonIDPaths(handler Handler) {
	h := make(handlers.MethodHandler, 3)
	if _, hasQuerier := handler.(Querier); hasQuerier {
		h["GET"] = http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			h := handler.New(req).(Querier)
			resp := handle(h, h.Query)
			WriteResponse(writer, resp, r.codecs)
		})
		h["HEAD"] = http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			h := handler.New(req).(Querier)
			resp := handle(h, h.Query)
			WriteHead(writer, resp, r.codecs)
		})
	}
	if _, hasPoster := handler.(Poster); hasPoster {
		h["POST"] = http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			h := handler.New(req).(Poster)
			resp := handle(h, h.Post)
			WriteResponse(writer, resp, r.codecs)
		})
	}
	if len(h) > 0 {
		r.Path(handler.Path()).Handler(h)
	}
}

// AddCodec registers a codec with this router.  Any codecs added in
// this way will be supplied to any *Response value that has not had
// its codecs set (via NewResponseForCodecs).
func (r *Router) AddCodec(codec Codec) {
	r.codecs = append(r.codecs, codec)
}

// Route routes the methods on handler to paths, based on handler's
// Path().
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
