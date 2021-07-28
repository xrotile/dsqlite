module dsqlite

go 1.15

replace github.com/xrotile/dsqlite/server => ./server

require (
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/goraft/raft v0.0.0-20150509002813-0061b6c82526
	github.com/gorilla/mux v1.8.0
)
