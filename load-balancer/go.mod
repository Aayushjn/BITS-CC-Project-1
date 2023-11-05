module github.com/aayushjn/load-balancer

go 1.21.1

require (
	github.com/go-chi/chi/v5 v5.0.10
	github.com/orcaman/concurrent-map/v2 v2.0.1
	golang.org/x/sync v0.4.0
)

replace (
	github.com/aayushjn/load-balancer/balancer/backend => ./balancer/backend
	github.com/aayushjn/load-balancer/balancer/strategy => ./balancer/strategy
	github.com/aayushjn/load-balancer/balancer => ./balancer
	github.com/aayushjn/load-balancer/errors => ./errors
	github.com/aayushjn/load-balancer/server => ./server
	github.com/aayushjn/load-balancer/util => ./util
)
