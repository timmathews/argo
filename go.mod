module github.com/asbjorn/argo

// module github.com/timmathews/argo

go 1.15

require (
	github.com/burntsushi/toml v0.2.0
	github.com/deckarep/golang-set v0.0.0-20170202203032-fc8930a5e645
	github.com/eclipse/paho.mqtt.golang v1.1.0
	github.com/gorilla/mux v1.4.0
	github.com/gorilla/websocket v1.2.0
	github.com/imdario/mergo v0.0.0-20160216103600-3e95a51e0639
	github.com/jacobsa/go-serial v0.0.0-20160401030210-6f298b3ae27e
	github.com/op/go-logging v0.0.0-20160211212156-b2cb9fa56473
	github.com/satori/go.uuid v1.1.0
	github.com/timmathews/argo v0.0.0-00010101000000-000000000000
	github.com/wsxiaoys/terminal v0.0.0-20160513160801-0940f3fc43a0
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	gopkg.in/vmihailenco/msgpack.v2 v2.9.1
)

replace github.com/timmathews/argo => github.com/asbjorn/argo v1.1.2-rc1
