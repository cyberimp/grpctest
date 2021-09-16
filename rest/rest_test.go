package rest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestLogServer_Serve(t *testing.T) {
	testServer := &LogServer{}
	logPipe := make(chan string)
	port := restPortSetup()
	testURL:=fmt.Sprintf("http://localhost:%d", port)
	testServer.Serve(logPipe)


	t.Run("Should open server", func(t *testing.T) {
		_, err := http.Get(testURL)
		if err != nil{
			t.Errorf("Did not created server, error: %q", err)
		}
	})



	tests := []struct {
		name string
		write string
	}{
		{"Should accept first log string", "test"},
		{"Should accept second log string", "test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logPipe <- tt.write
			_, err := http.Get(testURL)
			if err != nil{
				t.Errorf("Server closed unexpectedly")
			}
		})
	}

	t.Run("Should shutdown on pipe close", func(t *testing.T) {
		close(logPipe)
		time.Sleep(1*time.Second)
		_, err := http.Get(testURL)
		if err == nil{
			t.Errorf("Server did not shutdown")
		}
	})
}

func TestLogServer_accessLogHandle(t *testing.T) {

	var (
		tests = []struct {
			name	string
			want	int
		}{
			{name: "should have 0 records on start", want: 0},
			{name: "should have 1 records after writing to chan", want: 1},

		}
	)

	testSink := make(chan string)

	testServer := &LogServer{
		accessLog: []LogRecord{},
		sink:      testSink,
	}

	go testServer.accessLogHandle()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testResult := len(testServer.accessLog)
			if testResult != tt.want{
				t.Errorf("Wrong log size, got %v, wanted:%v",
					testResult, tt.want)
			}
			testSink <- "Test"
		})
	}
	close(testSink)

}


func TestLogServer_indexHandler(t *testing.T) {
	var (
		tests = []struct {
			name   	string
			log 	[]LogRecord
			want 	string
		}{
			{
				name: "should return [] by default",
				log: []LogRecord{},
				want: "[]",
			},
			{
				name: "should return json with written data",
				log: []LogRecord{{TimeStamp: time.Unix(0,0), AccessData: "User Name"}},
				want: "[{\"time_stamp\":\"1970-01-01T03:00:00+03:00\",\"access_data\":\"User Name\"}]",
			},
		}
	)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &LogServer{
				accessLog: tt.log,
				sink:      nil,
				srv:       nil,
			}
			rr := httptest.NewRecorder()
			s.indexHandler(rr, nil)
			result := rr.Body.String()
			if result != tt.want{
				t.Errorf("got wrong result, wanted:%v, got: %v", tt.want, result)
			}
		})
	}
}


func trySetPort (value string) {
	err := os.Setenv("PORT", value)
	if err != nil {
		return
	}
}

func Test_restPortSetup(t *testing.T) {
	tests := []struct {
		name string
		port string
		want uint16
	}{
		{name: "Should have 8080 by default", port: "", want: 8080},
		{name: "Should have 1234 when set to 1234", port: "1234", want: 1234},
	}

	reset := os.Getenv("port")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trySetPort(tt.port)
			if got := restPortSetup(); got != tt.want {
				t.Errorf("restPortSetup() = %v, want %v", got, tt.want)
			}
		})
	}

	trySetPort(reset)
}