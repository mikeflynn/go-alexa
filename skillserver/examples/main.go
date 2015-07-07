package main

import (
	"github.com/mikeflynn/go-alexa/skillserver"
	"net/http"
)

var Applications = map[string]skillserver.EchoApplication{
	"/echo/helloworld": skillserver.EchoApplication{ // Route
		AppID:   "xxxxxxxx",     // Echo App ID from Amazon Dashboard
		Handler: EchoHelloWorld, // Handler Func
	},
}

func main() {
	skillserver.Run(Applications, "3000")
}

func EchoHelloWorld(w http.ResponseWriter, r *http.Request) {
	echoReq := skillserver.GetEchoRequest(r)

	if echoReq.GetRequestType() == "IntentRequest" || echoReq.GetRequestType() == "LaunchRequest" {
		echoResp := skillserver.NewEchoResponse().OutputSpeech("Hello world from my new Echo test app!").Card("Hello World", "This is a test card.")

		json, _ := echoResp.String()
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.Write(json)
	}
}
