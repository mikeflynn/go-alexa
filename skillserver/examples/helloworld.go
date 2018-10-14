package main

import (
	alexa "github.com/mikeflynn/go-alexa/skillserver"
)

var applications = map[string]interface{}{
	"/echo/helloworld": alexa.EchoApplication{ // Route
		AppID:    "xxxxxxxx", // Echo App ID from Amazon Dashboard
		OnIntent: echoIntentHandler,
		OnLaunch: echoIntentHandler,
	},
}

func main() {
	alexa.Run(applications, "3000")
}

func echoIntentHandler(echoReq *alexa.EchoRequest, echoResp *alexa.EchoResponse) {
	echoResp.OutputSpeech("Hello world from my new Echo test app!").Card("Hello World", "This is a test card.")
}
