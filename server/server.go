package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/goraft/raft"
	"github.com/gorilla/mux"
)

type server struct {
	host string
	port int
	r    *mux.Router
	db   *DB
	rs   raft.Server
	name string
	path string // raft path
}

func NewServer(host string, port int, path string) *server {
	db := NewDB()
	Server := server{
		host: host,
		port: port,
		db:   db,
		path: path,
	}
	return &Server
}

func (s *server) ListenAndLeave(header string) error {
	// register command
	raft.RegisterCommand(&WriteCommand{})
	// create raft name
	rand.Seed(time.Now().Unix())
	s.r = mux.NewRouter()

	s.CreateRaftName()

	// create raft server
	transporter := raft.NewHTTPTransporter("/raft", 200*time.Millisecond)
	s.rs, _ = raft.NewServer(s.name, s.path, transporter, nil, s.db, s.ConnectString())
	transporter.Install(s.rs, s)
	// start raft server
	starterr := s.rs.Start()
	if starterr != nil {
		log.Printf("raft server start error", starterr)
	}

	// config router
	s.r.HandleFunc("/db/{key}", s.GetHandler).Methods("Get")
	s.r.HandleFunc("/db/{key}", s.WriteHandler).Methods("Post")
	s.r.HandleFunc("/join", s.JoinHandler).Methods("Post")
	s.r.HandleFunc("/peer", s.PeerHandler).Methods("Get")

	// create http server
	addr := s.host + ":" + strconv.Itoa(s.port)
	httpServer := &http.Server{
		Addr:    addr,
		Handler: s.r,
	}

	// try to join header
	if header != "" {
		// try to join hader
		joinCommand := raft.DefaultJoinCommand{
			Name:             s.rs.Name(),
			ConnectionString: s.ConnectString(),
		}
		var buffer bytes.Buffer
		json.NewEncoder(&buffer).Encode(joinCommand)
		resp, err := http.Post(fmt.Sprintf("http://%s/join", header), "application/json", &buffer)
		if err != nil {
			log.Printf("http post failed")
		}
		resp.Body.Close()
	} else if s.rs.IsLogEmpty() {
		log.Printf("init: " + s.ConnectString() + " " + s.rs.Name())
		_, err := s.rs.Do(&raft.DefaultJoinCommand{
			Name:             s.rs.Name(),
			ConnectionString: s.ConnectString(),
		})
		if err != nil {
			log.Printf("dsqlite raft start header failed %s", err)
		}
	} else {
		log.Printf("retrive from snapshot")
	}

	log.Printf("listen http server" + addr)
	return httpServer.ListenAndServe()
}

func (s *server) ConnectString() string {
	connectString := fmt.Sprintf("http://%s:%d", s.host, s.port)
	return connectString
}

func (s *server) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.r.HandleFunc(pattern, handler)
}

func (s *server) CreateRaftName() {
	raftPathName := filepath.Join(s.path, "name")
	if data, err := ioutil.ReadFile(raftPathName); err == nil {
		name := string(data)
		s.name = name
	} else {
		s.name = fmt.Sprintf("%07x", rand.Int())[0:7]
		log.Printf("raft name is %v", s.name)
		if err := ioutil.WriteFile(raftPathName, []byte(s.name), 0766); err != nil {
			log.Printf("raftName write to ")
		}
	}
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

	body, _ := ioutil.ReadAll(req.Body)
	var value string = string(body)

	// just save to local db.
	// s.db.Put(key, value)

	// put save command to raft server.
	_, err := s.rs.Do(NewWriteCommand(key, value))
	if err != nil {
		log.Printf("raft writeHandler failed")
	}
}

func (s *server) JoinHandler(w http.ResponseWriter, req *http.Request) {
	joinCommand := &raft.DefaultJoinCommand{}
	err := json.NewDecoder(req.Body).Decode(joinCommand)
	if err != nil {
		log.Printf("dsqlite joinhandler error")
	}
	_, err = s.rs.Do(joinCommand)
	log.Printf("server connection is comming %s", joinCommand.ConnectionString)
	if err != nil {
		log.Fatalf("dsqlite raft join failed, %s", err)
	}
}

func (s *server) PeerHandler(w http.ResponseWriter, req *http.Request) {
	peers := s.rs.Peers()
	value := fmt.Sprintf("the leader is %s, ", s.rs.Leader())
	for _, peer := range peers {
		value += peer.ConnectionString + "/"
	}
	w.Write([]byte(value))
}
