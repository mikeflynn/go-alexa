package chatter

import (
	"errors"
	"math/rand"
	"time"
)

var phrases map[string][]string = map[string][]string{
	"yes":       []string{"Yes", "Yup", "Absolutely", "Affirmative", "Okay", "Yeah"},
	"no":        []string{"No", "I don't think so.", "No way.", "Nope", "Negative", "Um...no"},
	"incorrect": []string{"Incorrect", "Erroneous", "False", "Inaccurate"},
	"correct":   []string{"Correct", "Accurate", "Right", "Valid"},
	"hello":     []string{"Hello", "Hi", "Hey", "Howdy", "Yo", "Good Day", "What's Up?", "Bonjour", "Greetings"},
	"goodbye":   []string{"Goodbye", "Bye Bye", "Later", "See ya", "Peace"},
}

func Get(section string) (string, error) {
	rand.Seed(time.Now().Unix())

	if v, ok := phrases[section]; ok {
		return phrases[rand.Intn(len(phrases))]
	}

	return "", errors.New("Invalid section")
}

func Set(section string, phrase string) (string, error) {
	if _, ok := phrases[section]; ok {
		phrases[section] = append(phrases[section], phrase)
	} else {
		phrases[section] = []string{phrase}
	}
}
