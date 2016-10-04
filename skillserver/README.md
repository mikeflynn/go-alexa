## `go-alexa/skillserver`

A simple Go framework to quickly create an Amazon Alexa Skills web service.

### Updates
10/3/16: Go 1.7 is required now as `go-alexa` uses the new core context library. It's not ideal to require 1.7, but with Go's no breaking changes promise it should be an easy upgrade for the vast majority of projects out there and it's better to keep up with the current release. If this change causes any issues, please reach out with an issue.

4/5/16: After taking a few good addtions from the community recently, I also just added new hooks that make it even easier to get going since you don't have to write a full `net/http` handler (see the new Hello World below)!

### What?

After beta testing the Amazon Echo (and it's voice assistant Alexa) for several months, Amazon has released the product to the public and created an SDK for developers to add new "Alexa Skills" to the product.

You can see the SDK documentation here: [developer.amazon.com/public/solutions/alexa/alexa-skills-kit](https://developer.amazon.com/public/solutions/alexa/alexa-skills-kit) but in short, a developer can make a web service that allows a user to say: _Alexa, ask [your service] to [some action your service provides]_

### Requirements

Amazon has a list of requirements to get a new Skill up and running

1. Creating your new Skill on their Development Dashboard populating it with details and example phrases. That process is documented here: [developer.amazon.com/appsandservices/solutions/alexa/alexa-skills-kit/docs/defining-the-voice-interface](https://developer.amazon.com/appsandservices/solutions/alexa/alexa-skills-kit/docs/defining-the-voice-interface)
2. A lengthy request validation proces. Documented here: [developer.amazon.com/public/solutions/alexa/alexa-skills-kit/docs/developing-an-alexa-skill-as-a-web-service](https://developer.amazon.com/public/solutions/alexa/alexa-skills-kit/docs/developing-an-alexa-skill-as-a-web-service) 
3. A formatted JSON response.
4. SSL connection required, even for development.

### How `skillserver` Helps

The `go-alexa/skillserver` takes care of #2 and #3 for you so you can concentrate on #1 and coding your app. (#4 is what it is. See the section on SSL below.)

### An Example App

Creating an Alexa Skill web service is easy with `go-alexa/skillserver`. Simply import the project as any other Go project, define your app, and write your endpoint. All the web service, security checks, and assistance in creating the response objects are done for you.

Here's a simple, but complete web service example:

```go
package main

import (
	alexa "github.com/mikeflynn/go-alexa/skillserver"
)

var Applications = map[string]interface{}{
	"/echo/helloworld": alexa.EchoApplication{ // Route
		AppID:    "xxxxxxxx", // Echo App ID from Amazon Dashboard
		OnIntent: EchoIntentHandler,
		OnLaunch: EchoIntentHandler,
	},
}

func main() {
	alexa.Run(Applications, "3000")
}

func EchoIntentHandler(echoReq *alexa.EchoRequest, echoResp *alexa.EchoResponse) {
	echoResp.OutputSpeech("Hello world from my new Echo test app!").Card("Hello World", "This is a test card.")
}
```

Details:
* You define your endpoints by creating a `map[string]interface{}` and loading it with `EchoApplication` types that specify the Application ID and handler function.
* All Skill endpoints must start with `/echo/` as that's the route grouping that has the security middleware.
* The easiest way to get started is define handler functions by using `OnIntent`, `OnLaunch`, or `OnSessionEnded` that take an EchoRequest and an EchoResponse.
* ...but if you want full control you can still use the `EchoApplication.Handler` hook to write a regular `net/http` handler so you have full access to the request and ResponseWriter.
* The JSON from the Echo request is already parsed for you. Grab it by calling `skillserver.GetEchoRequest(r *http.Request)`.
* You generate the Echo Response by using the EchoResponse struct that has methods to generate each part and that's it! ...unless you use the `EchoApplication.Handler` hook. In that case you need to write your JSON to the string with the `EchoResponse.toString()` method.

### The SSL Requirement

Amazon requires an SSL connection for all steps in the Skill process, even local development (which still gets requests from the Echo web service). Amazon is pushing their AWS Lamda service that takes care of SSL for you, but Go isn't an option on Lamda. What I've done personally is put Nginx in front of my Go app and let Nginx handle the SSL (a self-signed cert for development and a real cert when pushing to production). More information here on  [nginx.com](https://www.nginx.com/blog/nginx-ssl/).

### Contributors

Mike Flynn ([@thatmikeflynn](https://twitter.com/thatmikeflynn))
