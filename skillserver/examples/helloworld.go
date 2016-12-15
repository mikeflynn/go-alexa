package main

import (
	alexa "github.com/mikeflynn/go-alexa/skillserver"
	"log"
	"net/http"
)

func main() {
	echoApp := alexa.NewSkillHandler("xxxxxxx")
	echoApp.OnIntent = EchoIntentHandler
	echoApp.OnLaunch = EchoIntentHandler
	http.Handle("/echo/helloworld", echoApp)
	log.Fatalf("Stopped listening: %+v", http.ListenAndServe(":8080", nil))
}

func EchoIntentHandler(echoReq *alexa.EchoRequest, echoResp *alexa.EchoResponse) {
	echoResp.OutputSpeech("Hello world from my new Echo test app!").Card("Hello World", "This is a test card.")
}
