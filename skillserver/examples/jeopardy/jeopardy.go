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

	"github.com/kennygrant/sanitize"
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
	rand.Seed(time.Now().UTC().UnixNano())
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
	"Whatever. I'm just sitting here I guess.",
	"No. I'm busy. Just kidding.",
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
	"You got it! ",
	"Nice. ",
	"Bingo. ",
	"Nailed it! ",
	"Correct. ",
	"That is right. ",
	"Holy crap you got it! ",
}

var JeopardyWrongAnswer = []string{
	"Nope.",
	"Sorry, that's incorrect.",
	"Wow, no, not even close.",
	"Yikes. No.",
	"Awww, too bad.",
}

// #1: Greeting and ask for category.
// #2: Deliver question with dollar amount.
// #3: Verify answer (update session total)
// #4: Give current score and list categories.

func EchoJeopardy(w http.ResponseWriter, r *http.Request) {
	echoReq := alexa.GetEchoRequest(r)

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
		log.Println(echoReq.GetIntentName())

		session := getJeopardySession(col, echoReq.GetSessionID())

		var echoResp *alexa.EchoResponse

		switch echoReq.GetIntentName() {
		case "StartJeopardy":
			echoResp, session = jeopardyStart(echoReq, session)
		case "ListCategories":
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
		case "RepeatQuestion":
			if session.CurrentQuestion.Question != "" {
				echoResp = alexa.NewEchoResponse().OutputSpeech(session.CurrentQuestion.Question).EndSession(false)
			} else {
				echoResp = alexa.NewEchoResponse().OutputSpeech("I didn't ask a question. Please pick a category first.").EndSession(false)
			}
		case "QuitGame":
			noun := "questions"
			if session.NumQuestions == 1 {
				noun = "question"
			}
			echoResp = alexa.NewEchoResponse().OutputSpeech("You ended with " + strconv.Itoa(session.Dollars) + " after " + strconv.Itoa(session.NumQuestions) + " " + noun + " .").EndSession(true)
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

	msg := ""
	if echoReq.GetIntentName() == "StartJeopardy" {
		msg = JeopardyGreetings[rand.Intn(len(JeopardyGreetings))] + " Please pick one of the following categories: "
	} else {
		msg = "The categories are "
	}
	msg += strings.Join(catNames, ", ")
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

	msg += "From " + category + " for " + strconv.Itoa(clue.Value) + ". " + clue.Question + ". I need your answer in the form of a question."

	session.CurrentQuestion.Category = category
	session.CurrentQuestion.Answer = sanitize.HTML(clue.Answer)
	session.CurrentQuestion.Question = clue.Question
	session.CurrentQuestion.Value = clue.Value

	log.Println(session.CurrentQuestion.Question)
	log.Println(session.CurrentQuestion.Answer)

	echoResp.OutputSpeech(msg).Card("Question", msg).Reprompt("Times up. I need your answer in the form of a question.").EndSession(false)

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
	log.Println(answer)
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
	log.Println("http://jservice.io/api/clues?category=" + strconv.Itoa(catID) + "&offset=" + strconv.Itoa(offset))
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
