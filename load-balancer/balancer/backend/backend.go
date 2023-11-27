package backend

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
	"time"

	"github.com/aayushjn/load-balancer/util"
)

type Backend struct {
	URL              *url.URL `json:"url"`
	alive            atomic.Bool
	ReverseProxy     *httputil.ReverseProxy `json:"-"`
	inFlightRequests atomic.Uint64
}

func (b *Backend) SetAlive(alive bool) {
	b.alive.Store(alive)
}

func (b *Backend) IsAlive() bool {
	return b.alive.Load()
}

func (b *Backend) IncrementInFlightRequests() {
	b.inFlightRequests.Add(1)
}

func (b *Backend) DecrementInFlightRequests() {
	b.inFlightRequests.Add(^uint64(0))
}

func (b *Backend) GetInFlightRequests() uint64 {
	return b.inFlightRequests.Load()
}

func (b *Backend) Test() error {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", b.URL.Host, timeout)
	if err != nil {
		return err
	}
	_ = conn.Close()
	return nil
}

func (b *Backend) SetProxyHandler(requestHandler http.HandlerFunc) {
	b.ReverseProxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		fmt.Printf("[%s] %s\n", b.URL.Host, err.Error())
		retries := util.GetRetriesFromContext(r)
		if retries < 3 {
			time.Sleep(10 * time.Millisecond)
			ctx := context.WithValue(r.Context(), util.Retries, retries+1)
			b.ReverseProxy.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		b.SetAlive(false)

		attempts := util.GetAttemptsFromContext(r)
		fmt.Printf("%s(%s) Attempting retry %d\n", r.RemoteAddr, r.URL.Path, attempts)
		ctx := context.WithValue(r.Context(), util.Attempts, attempts+1)
		requestHandler(w, r.WithContext(ctx))
	}
}

func NewBackend(backendUrl string) (*Backend, error) {
	serverUrl, err := url.Parse(backendUrl)
	if err != nil {
		return nil, err
	}

	return &Backend{
		URL:              serverUrl,
		ReverseProxy:     httputil.NewSingleHostReverseProxy(serverUrl),
		alive:            atomic.Bool{},
		inFlightRequests: atomic.Uint64{},
	}, nil
}
