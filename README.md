## `go-alexa`: A Go toolset for creating Amazon Alexa Skills

The Amazon Echo, with it's voice assitant Alexa, is a surprisingly amazing tool. Having the power of voice recognition tied to the web ready at any time is quite powerful and now that Amazon has opened up a developer platform it's even more exciting!

Amazon has supplied packages for Java and Node.js (tied to the AWS Lamda platform) but I wanted to develop my skills in Go. As I moved through the process making my app work with Amazon's spec, a simple web framework that took care all the heavy lifting on security and crafting the response object formed. I'm looking forward to more Go-based tools getting created and living in this `go-alexa` bucket but for now the `skillserver` is the first tool.

Mike Flynn gave a talk about this library an conversational applications in general at the 2016 Strange Loop Conference: ["Exploring Conversational Interfaces with Amazon Alexa and Go"](https://www.youtube.com/watch?v=pDdE3PKy6mo)

### Tools

* [`skillserver`](skillserver/) - A framework to quickly create a skill web service that handles all of the Amazon requirements.
  * Example: [Jeopardy](skillserver/examples/jeopardy)

### Future Proposed Tools

* An Amazon Echo request simulator
* A library for Alexa responses

### Original Author

Mike Flynn ([@thatmikeflynn](http://twitter.com/thatmikeflynn))
