package main

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/context"
	alexa "github.com/mikeflynn/go-alexa/skillserver"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var Applications = map[string]interface{}{
	"/echo/jeopardy": alexa.EchoApplication{
		AppID:   os.Getenv("JEOPARDY_APP_ID"),
		Handler: EchoJeopardy,
	},
}

func main() {
	alexa.Run(Applications, "3000")
}

type JeopardySession struct {
	AWSID           string
	Dollars         int
	NumQuestions    int
	CurrentQuestion JeopardyQuestion
	UpdatedAt       int64
}

type JeopardyQuestion struct {
	Category string
	Question string
	Answer   string
	Value    int
}

var JeopardyCategories = map[string]int{
	"food":       49,
	"hodgepodge": 227,
	"history":    114,
	"sports":     42,
	"science":    25,
	"television": 67,
	"people":     442,
	"rhyme time": 561,
	"pop music":  770,
	"quotations": 1420,
}

var JeopardyGreetings = []string{
	"Sure.",
	"Lets do it!",
	"Hell to the yeah!",
	"Whatever. I'm just sitting here I guess.",
	"No. I'm busy...Just kidding.",
	"Coolio.",
	"Lets play Jeopardy!",
}

var JeopardyCatSelect = []string{
	"Don't worry your pretty little head about it. Lets go with ",
	"I'll pick ",
	"Lets go with ",
	"You should have picked. I'm going with ",
	"You've already failed, but lets keep going with ",
}

var JeopardyRightAnswer = []string{
	"You got it!",
	"Nice.",
	"Bingo.",
	"Nailed it!",
	"Correct",
	"That is right.",
	"Holy crap you got it!",
}

var JeopardyWrongAnswer = []string{
	"Nope.",
	"Sorry, that's incorrect.",
	"Wow...no...not even close.",
	"Yikes. No.",
	"Awww, too bad.",
}

var JeopardyThinking = []string{
	"I'm going to talk a little while and give you a chance to think about the answer. If I stop talking I need answer right away so I'll just jabber on about whatever...until...now.",
	"While you think I'm going to sing the Who's the Boss theme song. There's a time for love and a time for living. Take a chance and face the wind! Open road and a road that's taken. A brand new life around the bend! Ok.",
	"How long does Jeopardy give you to answer a question? Not very long but long enough for Alec to say We need an answer and then that beep beep comes in. Well, guess what?",
}

// #1: Greeting and ask for category.
// #2: Deliver question with dollar amount.
// #3: Verify answer (update session total)
// #4: Give current score and list categories.

func EchoJeopardy(w http.ResponseWriter, r *http.Request) {
	echoReq := context.Get(r, "echoRequest").(*alexa.EchoRequest)

	// Start up Mongo!
	mongodb, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}

	col := mongodb.DB("echo").C("jeopardy")
	defer mongodb.Close()

	log.Println(echoReq.GetRequestType())
	log.Println(echoReq.GetSessionID())

	if echoReq.GetRequestType() == "LaunchRequest" {
		session := getJeopardySession(col, echoReq.GetSessionID())

		echoResp, session := jeopardyStart(echoReq, session)

		json, _ := echoResp.String()
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.Write(json)
	} else if echoReq.GetRequestType() == "IntentRequest" {
		session := getJeopardySession(col, echoReq.GetSessionID())

		log.Println(echoReq.GetIntentName())

		var echoResp *alexa.EchoResponse

		switch echoReq.GetIntentName() {
		case "StartJeopardy":
			echoResp, session = jeopardyStart(echoReq, session)
		case "PickCategory":
			if session.CurrentQuestion.Category == "" {
				echoResp, session = jeopardyCategory(echoReq, session)
			} else {
				echoResp, session = jeopardyAnswer(echoReq, session)
			}

			session.Update(col)
		case "AnswerQuestion":
			echoResp, session = jeopardyAnswer(echoReq, session)
			session.Update(col)
		case "QuitGame":
			echoResp = alexa.NewEchoResponse().OutputSpeech("Ok. You ended with " + strconv.Itoa(session.Dollars) + " after " + strconv.Itoa(session.NumQuestions) + " questions.").EndSession(true)
		default:
			echoResp = alexa.NewEchoResponse().OutputSpeech("I'm sorry, I didn't get that. Can you say that again?").EndSession(false)
		}

		json, _ := echoResp.String()
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.Write(json)
	} else if echoReq.GetRequestType() == "SessionEndedRequest" {
		//session.Delete(col)
	}
}

func jeopardyStart(echoReq *alexa.EchoRequest, session *JeopardySession) (*alexa.EchoResponse, *JeopardySession) {
	catNames := []string{}
	for k, _ := range JeopardyCategories {
		catNames = append(catNames, k)
	}

	msg := JeopardyGreetings[rand.Intn(len(JeopardyGreetings))] + " Please pick one of the following categories: " + strings.Join(catNames, ", ")
	echoResp := alexa.NewEchoResponse().OutputSpeech(msg).EndSession(false)

	return echoResp, session
}

func jeopardyCategory(echoReq *alexa.EchoRequest, session *JeopardySession) (*alexa.EchoResponse, *JeopardySession) {
	msg := ""
	echoResp := alexa.NewEchoResponse()

	// Declare the category
	category, err := echoReq.GetSlotValue("Category")
	_, catExists := JeopardyCategories[category]
	if err != nil || !catExists {
		catNames := []string{}
		for k, _ := range JeopardyCategories {
			catNames = append(catNames, k)
		}

		category = getRandom(catNames)

		msg = msg + getRandom(JeopardyCatSelect) + category + ". "
	} else {
		category = strings.ToLower(category)
	}

	clue, err := getJServiceClue(JeopardyCategories[category])
	if err != nil {
		clue, err = getJServiceClue(JeopardyCategories[category])
		if err != nil {
			echoResp := alexa.NewEchoResponse().OutputSpeech("I'm sorry, but I can't seem to get a question right now.").EndSession(true)
			return echoResp, session
		}
	}

	msg += "From " + category + " for " + strconv.Itoa(clue.Value) + ". " + clue.Question + "." + getRandom(JeopardyThinking) + " I need your answer in the form of a question!"

	session.CurrentQuestion.Category = category
	session.CurrentQuestion.Answer = clue.Answer
	session.CurrentQuestion.Question = clue.Question
	session.CurrentQuestion.Value = clue.Value

	echoResp.OutputSpeech(msg).Card("Question", msg).EndSession(false)

	return echoResp, session
}

func jeopardyAnswer(echoReq *alexa.EchoRequest, session *JeopardySession) (*alexa.EchoResponse, *JeopardySession) {
	msg := ""
	echoResp := alexa.NewEchoResponse()

	if session.CurrentQuestion.Answer == "" {
		echoResp.OutputSpeech("You need to get a question before answering!").EndSession(false)
		return echoResp, session
	}

	answer, err := echoReq.GetSlotValue("Answer")
	if err != nil {
		echoResp.OutputSpeech("We need an answer!").EndSession(false)
		return echoResp, session
	}

	if strings.ToLower(answer) == strings.ToLower(session.CurrentQuestion.Answer) {
		msg += getRandom(JeopardyRightAnswer)
		session.Dollars = session.Dollars + session.CurrentQuestion.Value
	} else {
		msg += getRandom(JeopardyWrongAnswer) + " The correct answer was " + session.CurrentQuestion.Answer + ". "
	}

	session.NumQuestions++
	session.CurrentQuestion = JeopardyQuestion{}

	msg += "You're at " + strconv.Itoa(session.Dollars) + " after " + strconv.Itoa(session.NumQuestions) + " questions. Please pick another category."

	echoResp.OutputSpeech(msg).Card("Answer", session.CurrentQuestion.Answer).EndSession(false)

	return echoResp, session
}

func getRandom(list []string) string {
	return list[rand.Intn(len(list))]
}

func getJeopardySession(col *mgo.Collection, sessid string) *JeopardySession {
	user := &JeopardySession{}
	err := col.Find(bson.M{"awsid": sessid}).One(&user)
	if err != nil || user.AWSID == "" {
		log.Println("Creating new session.")
		user.AWSID = sessid
		user.NumQuestions = 0
		user.Dollars = 0
		user.UpdatedAt = time.Now().Unix()
		err = col.Insert(&user)
		if err != nil {
			panic(err)
		}
	}

	return user
}

func (this *JeopardySession) Update(col *mgo.Collection) error {
	this.UpdatedAt = time.Now().Unix()
	err := col.Update(bson.M{"awsid": this.AWSID}, this)
	if err != nil {
		return err
	}

	return nil
}

func (this *JeopardySession) Delete(col *mgo.Collection) error {
	err := col.Remove(bson.M{"awsid": this.AWSID})
	if err != nil {
		return err
	}

	return nil
}

// jservice

type JServiceClue struct {
	ID           int    `json:"id"`
	Answer       string `json:"answer"`
	Question     string `json:"question"`
	Value        int    `json:"value"`
	AirDate      string `json:"airdate"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	CategoryID   int    `json:"category_id"`
	GameID       int    `json:"game_id"`
	InvalidCount int    `json:"invalid_count,omitempty"`
	Category     struct {
		ID         int    `json:"id"`
		Title      string `json:"title"`
		CreatedAt  string `json:"created_at"`
		UpdatedAt  string `json:"updated_at"`
		CluesCount int    `json:"clues_count"`
	} `json:"category"`
}

func getJServiceClue(catID int) (JServiceClue, error) {
	offset := rand.Intn(100)
	response, err := http.Get("http://jservice.io/api/clues?category=" + strconv.Itoa(catID) + "&offset=" + strconv.Itoa(offset))
	if err != nil {
		log.Println(err.Error())
		return JServiceClue{}, err
	}

	var clues []JServiceClue
	err = json.NewDecoder(response.Body).Decode(&clues)
	if err != nil {
		log.Println(err.Error())
		return JServiceClue{}, err
	}

	if len(clues) > 0 {
		return clues[0], nil
	}

	return JServiceClue{}, errors.New("No clues returned.")
}
