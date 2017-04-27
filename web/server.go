package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mijia/sweb/log"
	"github.com/mijia/sweb/render"
	"github.com/mijia/sweb/server"
	"golang.org/x/net/context"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Server struct {
	*server.Server
	render         *render.Render
	isDebug        bool
	assetsDomain   string
	muxControllers []MuxController
}

func (s *Server) addMuxController(mcs ...MuxController) {
	s.muxControllers = append(s.muxControllers, mcs...)
}

func (s *Server) ListenAndServe(addr string) error {
	if s.isDebug && s.assetsDomain != "" {
		log.Infof("Use AssetsPrefix %s", s.assetsDomain)
		s.EnableAssetsPrefix(s.assetsDomain)
	}
	s.addMuxController(&UserController{})
	s.addMuxController(&ReservationController{})
	s.render = s.initRender()

	ignoredUrls := []string{"/bundles/", "/fonts/", "/debug/vars", "/favicon", "/robots"}
	s.Middleware(server.NewRecoveryWare(s.isDebug))
	s.Middleware(server.NewStatWare(ignoredUrls...))
	s.Middleware(server.NewRuntimeWare(ignoredUrls, true, 15*time.Minute))
	s.Middleware(s.wareWebpackAssets("webpack-assets.json", "bundles"))

	s.Get("/debug/vars", "RuntimeStat", s.getRuntimeStat)
	s.Files("/static/*filepath", http.Dir("static"))
	s.Files("/assets/*filepath", http.Dir("public"))
	for _, mc := range s.muxControllers {
		mc.SetResponseRenderer(s)
		mc.SetUrlReverser(s)
		mc.MuxHandlers(s)
	}

	return s.Run(addr)
}

func (s *Server) getRuntimeStat(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	http.DefaultServeMux.ServeHTTP(w, r)
	return ctx
}

func (s *Server) wareWebpackAssets(webpackAssetsFile string, subPath string) server.MiddleFn {
	var loadOnce sync.Once
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, next server.Handler) context.Context {
		if s.isDebug {
			s.loadWebpackAssets(webpackAssetsFile, subPath)
		} else {
			loadOnce.Do(func() {
				s.loadWebpackAssets(webpackAssetsFile, subPath)
			})
		}
		return next(ctx, w, r)
	}
}

func (s *Server) loadWebpackAssets(webpackAssetsFile string, subPath string) {
	start := time.Now()
	if data, err := ioutil.ReadFile(webpackAssetsFile); err == nil {
		var packMappings map[string]map[string]string
		if err := json.Unmarshal(data, &packMappings); err == nil {
			newMappings := make(map[string]string)
			for entry, types := range packMappings {
				for ty, target := range types {
					if subPath != "" {
						target = fmt.Sprintf("%s/%s", subPath, target)
					}
					newMappings[fmt.Sprintf("%s.%s", entry, ty)] = target
				}
			}
			s.EnableExtraAssetsMapping(newMappings)
			log.Infof("[Server] Loaded webpack assets from %s, duration=%s", webpackAssetsFile, time.Now().Sub(start))
		} else {
			log.Errorf("[Server] Failed to decode the web pack assets, %s", err)
		}
	} else {
		log.Errorf("[Server] Failed to load webpack assets from %s, %s", webpackAssetsFile, err)
	}
}

func (s *Server) initRender() *render.Render {
	tSets := []*render.TemplateSet{}
	for _, mc := range s.muxControllers {
		mcTSets := mc.GetTemplates()
		tSets = append(tSets, mcTSets...)
	}
	r := render.New(render.Options{
		Directory:     "templates",
		Funcs:         s.renderFuncMaps(),
		Delims:        render.Delims{"{{", "}}"},
		IndentJson:    true,
		UseBufPool:    true,
		IsDevelopment: s.isDebug,
	}, tSets)
	log.Info("Templates loaded ...")
	return r
}

func formatTime(tm time.Time, layout string) string {
	return tm.Format(layout)
}

func (s *Server) renderFuncMaps() []template.FuncMap {
	funcs := make([]template.FuncMap, 0)
	funcs = append(funcs, s.DefaultRouteFuncs())
	funcs = append(funcs, template.FuncMap{
		"add": func(input interface{}, toAdd int) float64 {
			switch t := input.(type) {
			case int:
				return float64(t) + float64(toAdd)
			case int64:
				return float64(t) + float64(toAdd)
			case int32:
				return float64(t) + float64(toAdd)
			case float32:
				return float64(t) + float64(toAdd)
			case float64:
				return t + float64(toAdd)
			default:
				return float64(toAdd)
			}
		},
		"formatTime": formatTime,
	})
	return funcs
}

func (s *Server) RenderJsonOr500(w http.ResponseWriter, status int, v interface{}) {
	s.renderJsonOr500(w, status, v)
}

func (s *Server) renderJsonOr500(w http.ResponseWriter, status int, v interface{}) {
	if err := s.render.Json(w, status, v); err != nil {
		log.Errorf("Server got a json rendering error, %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) RenderHtmlOr500(w http.ResponseWriter, status int, name string, binding interface{}) {
	s.renderHtmlOr500(w, status, name, binding)
}

func (s *Server) RenderRawHtml(w http.ResponseWriter, status int, htmlString string) {
	s.renderString(w, status, htmlString)
}

func (s *Server) renderHtmlOr500(w http.ResponseWriter, status int, name string, binding interface{}) {
	w.Header().Set("Cache-Control", "no-store, no-cache")
	if err := s.render.Html(w, status, name, binding); err != nil {
		log.Errorf("Server got a rendering error, %q", err)
		if s.isDebug {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			//渲染一个500错误页面
			s.RenderError500(w, err)
		}
	}
}

func (s *Server) renderString(w http.ResponseWriter, status int, data string) {
	out := new(bytes.Buffer)
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(status)
	out.Write([]byte(data))
	io.Copy(w, out)
}

func (s *Server) Get(path, name string, handle server.Handler) {
	newHandle := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (newCtx context.Context) {
		newCtx = ctx
		defer func() {
			if err := recover(); err != nil {
				stack := make([]byte, 1024*8)
				stack = stack[:runtime.Stack(stack, s.isDebug)]
				msg := fmt.Sprintf("Request: %s \r\n PANIC: %s\n%s", r.URL.String(), err, stack)
				log.Error(msg)
				s.RenderError500(w, errors.New(msg))
			}
		}()
		newCtx = handle(ctx, w, r)
		return
	}
	s.Server.Get(path, name, newHandle)
}

func (s *Server) GetJson(path string, name string, handle JsonHandler) {
	s.Get(path, name, s.makeJsonHandler(handle))
}

func (s *Server) PostJson(path string, name string, handle JsonHandler) {
	s.Post(path, name, s.makeJsonHandler(handle))
}

func (s *Server) PutJson(path string, name string, handle JsonHandler) {
	s.Put(path, name, s.makeJsonHandler(handle))
}

func (s *Server) PatchJson(path string, name string, handle JsonHandler) {
	s.Patch(path, name, s.makeJsonHandler(handle))
}

func (s *Server) DeleteJson(path string, name string, handle JsonHandler) {
	s.Delete(path, name, s.makeJsonHandler(handle))
}

func (s *Server) makeJsonHandler(handle JsonHandler) server.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
		// request log
		switch r.Method {
		case http.MethodGet:
			if r.URL.RawQuery != "" {
				log.Infof("%s %s?%s", strings.ToUpper(r.Method), r.URL.Path, r.URL.RawQuery)
			} else {
				log.Infof("%s %s", strings.ToUpper(r.Method), r.URL.Path)
			}
		case http.MethodPost:
			r.ParseForm()
			log.Infof("%s %s %+v", strings.ToUpper(r.Method), r.URL.Path, r.PostForm)
		}

		statusCode, data := handle(ctx, w, r)

		// response log
		if r.Method != http.MethodGet {
			log.Infof("RESPONSE %d: %s %+v", statusCode, r.URL.Path, data)
		}

		s.renderJsonOr500(w, statusCode, data)
		return ctx
	}
}

func (s *Server) SetAssetDomain(domain string) {
	s.assetsDomain = domain
}

func (s *Server) Cleanup() {
}

func NewServer(isDebug bool) *Server {
	if isDebug {
		log.EnableDebug()
	}

	ctx := context.Background()

	srv := &Server{
		Server:         server.New(ctx, isDebug),
		muxControllers: []MuxController{},
		isDebug:        isDebug,
	}
	return srv
}
