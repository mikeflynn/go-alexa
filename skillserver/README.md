## `go-alexa/skillserver`

A simple Go framework to quickly create an Amazon Alexa Skills web service.

### What?

After beta testing the Amazon Echo (and it's voice assistant Alexa) for several months, Amazon has released the product to the public and created an SDK for developers to add new "Alexa Skills" to the product.

You can see the SDK documentation here: [developer.amazon.com/public/solutions/alexa/alexa-skills-kit](https://developer.amazon.com/public/solutions/alexa/alexa-skills-kit) but in short, a developer can make a web service that allows a user to say: _Alexa, ask [your service] to [some action your service provides]_

### Requirements

Amazon has a list of requirements to get a new Skill up and running

1. Creating your new Skill on their Development Dashboard populating it with details and example phrases. That process is documented here: [developer.amazon.com/appsandservices/solutions/alexa/alexa-skills-kit/docs/defining-the-voice-interface](https://developer.amazon.com/appsandservices/solutions/alexa/alexa-skills-kit/docs/defining-the-voice-interface)
2. A lengthy request validation proces. Documented here: 
3. A formatted JSON response.
4. SSL connection required, even for development.

### How `skillserver` Helps

The `go-alexa/skillserver` takes care of #2 and #3 for you so you can concentrate on #1 and coding your app. (#4 is what it is. See the section on SSL below.)

### An Example App



### The SSL Requirement
