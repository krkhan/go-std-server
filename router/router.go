package router

import (
	"context"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type ContextKey struct{}

func NewRoute(name, method, pattern string, handler http.HandlerFunc) Route {
	return Route{
		Name: name,
		Method:  method,
		Regex:   regexp.MustCompile("^" + pattern + "$"),
		Handler: handler,
		Stats: &RouteStats{
			TotalRequests: 0,
			TotalTime:     0,
			StatsLock:     &sync.RWMutex{},
		},
	}
}

type RouteStats struct {
	TotalRequests uint64
	TotalTime     uint64
	StatsLock     *sync.RWMutex
}

type Route struct {
	Name    string
	Method  string
	Regex   *regexp.Regexp
	Handler http.HandlerFunc
	Stats   *RouteStats
}

func Serve(routes []Route, w http.ResponseWriter, r *http.Request) {
	var disallowedMethods []string
	for _, route := range routes {
		matches := route.Regex.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
			if r.Method != route.Method {
				disallowedMethods = append(disallowedMethods, route.Method)
				continue
			}
			ctx := context.WithValue(r.Context(), ContextKey{}, matches[1:])
			handlerStart := time.Now()
			route.Handler(w, r.WithContext(ctx))
			handlerEnd := time.Now()
			elapsedMicroseconds := handlerEnd.Sub(handlerStart).Microseconds()
			go func() {
				// Probably not needed for the time being (stuff we are tracking can be atomically increased)
				// But can be helpful for tracking more complex stats
				route.Stats.StatsLock.Lock()
				defer route.Stats.StatsLock.Unlock()
				atomic.AddUint64(&(route.Stats.TotalRequests), 1)
				atomic.AddUint64(&(route.Stats.TotalTime), uint64(elapsedMicroseconds))
			}()
			return
		}
	}
	if len(disallowedMethods) > 0 {
		w.Header().Set("Allow", strings.Join(disallowedMethods, ", "))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.NotFound(w, r)
}

func GetParam(r *http.Request, index int) string {
	fields := r.Context().Value(ContextKey{}).([]string)
	return fields[index]
}
