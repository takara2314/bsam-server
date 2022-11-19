module bsam-server

// +heroku goVersion go1.18
go 1.16

require (
	github.com/gin-contrib/cors v1.4.0
	github.com/gin-gonic/gin v1.8.1
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/gorilla/websocket v1.5.0
	github.com/lib/pq v1.10.6
	github.com/xo/dburl v0.2.0
	gopkg.in/yaml.v3 v3.0.1
)
