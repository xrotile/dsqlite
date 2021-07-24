package server

import (
	"db"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type server struct {
	host string
	port int
	r    *mux.Router
	db   *db
}

func NewServer(host string, port int) *server {
	db := db.NewDB()
	Server := server{
		host: host,
		port: port,
		db:   db,
	}
	return &Server
}

func (s *server) ListenAndLeave() error {

	s.r = mux.NewRouter()
	// config router
	s.r.HandleFunc("/db/{key}", s.GetHandler).Methods("Get")
	s.r.HandleFunc("/db/{key}", s.WriteHandler).Methods("Post")
	s.r.HandleFunc("/join", s.JoinHandler).Methods("Post")

	// create http server
	addr := s.host + ":" + strconv.Itoa(s.port)
	httpSever := &http.Server{
		Addr:    addr,
		Handler: s.r,
	}

	log.Printf("listen http server" + addr)
	return httpSever.ListenAndServe()
}

func (s *server) GetHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	log.Printf("dsqlite: get db value")
	key := vars["key"]
	value, err := s.db.Get(key)
	if err {
		w.Write([]byte(value))
	} else {
		errp := "The key is invalid"
		w.Write([]byte(errp))
	}
}

func (s *server) WriteHandler(w http.ResponseWriter, req *http.Request) {
	// single node write
	// get key from url path, get value from body.
	vars := mux.Vars(req)
	key := vars["key"]

	body := ioutils.ReadAll(req.Body)
	var value string = string(body)

	// just save to local db.
	s.db.Put(key, value)

	// put save command to raft server.
}

func (s *server) JoinHandler(w http.ResponseWriter, req *http.Request) {

}
