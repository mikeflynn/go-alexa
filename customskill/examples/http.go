package main

import (
	"net/http"

	"github.com/mikeflynn/go-alexa/customskill"
	"github.com/mikeflynn/go-alexa/customskill/request"
	"github.com/mikeflynn/go-alexa/customskill/response"
)

var (
	s customskill.Skill
)

func main() {
	s = customskill.Skill{
		OnLaunch: onLaunch,
	}

	http.HandleFunc("/", s.DefaultHTTPHandler)
	http.ListenAndServe(":8080", nil)
}

func onLaunch(launchRequest *request.LaunchRequest) (*response.Response, map[string]interface{}, error) {
	resp := response.New()
	resp.SetEndSession(response.Bool(true))
	sessAttrs := make(map[string]interface{})
	sessAttrs["hello"] = "world"
	return resp, sessAttrs, nil
}
