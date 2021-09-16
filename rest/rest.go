package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type LogRecord struct {
	TimeStamp  time.Time `json:"time_stamp"`
	AccessData string    `json:"access_data"`
}

type LogServer struct {
	accessLog []LogRecord
	logMutex  sync.Mutex
	sink      chan string
	srv       *http.Server
}

func restPortSetup() uint16 {
	port := os.Getenv("PORT")
	result, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		result = 8080
	}
	return uint16(result)
}

//Serve rest server
func (s *LogServer) Serve(pipe chan string) {
	port := restPortSetup()

	s.sink = pipe
	s.accessLog = make([]LogRecord, 0)

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.indexHandler)

	s.srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	go func() {
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen: ", err)
		}
	}()

	log.Printf("Listening on port %d", port)
	log.Printf("Open http://localhost:%d in the browser", port)

	go s.accessLogHandle()
}

func (s *LogServer) accessLogHandle() {

	for logString := range s.sink {
		s.logMutex.Lock()
		s.accessLog = append(s.accessLog,
			LogRecord{
				TimeStamp:  time.Now(),
				AccessData: logString})
		s.logMutex.Unlock()
	}

	s.shutdown(time.Second * 5)
}

func (s *LogServer) shutdown(timeout time.Duration) {
	if s.srv == nil {
		return
	}

	ctxShutDown, cancel := context.WithTimeout(context.Background(), timeout)
	defer func() {
		cancel()
	}()

	err := s.srv.Shutdown(ctxShutDown)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("server exited properly")
}

func (s *LogServer) indexHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	s.logMutex.Lock()
	body, err := json.Marshal(s.accessLog)
	s.logMutex.Unlock()
	if err != nil {
		return
	}

	_, err = w.Write(body)
	if err != nil {
		return
	}
}
