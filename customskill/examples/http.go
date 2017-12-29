package main

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/mikeflynn/go-alexa/customskill"
	"github.com/mikeflynn/go-alexa/customskill/request"
	"github.com/mikeflynn/go-alexa/customskill/response"
)

var (
	s customskill.Skill
)

func handler(w http.ResponseWriter, r *http.Request) {
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, r.Body); err != nil {
		log.Panicf("failed to copy request body to buffer: %v", err)
	}
	defer r.Body.Close()
	if err := s.Handle(w, buf.Bytes()); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("Got an error: %v", err)
	}
}

func main() {
	s = customskill.Skill{
		OnLaunch: onLaunch,
	}

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func onLaunch(launchRequest *request.LaunchRequest) (*response.Response, map[string]interface{}, error) {
	resp := response.New()
	resp.SetEndSession(response.Bool(true))
	sessAttrs := make(map[string]interface{})
	sessAttrs["hello"] = "world"
	return resp, sessAttrs, nil
}
