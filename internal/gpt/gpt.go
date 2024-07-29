package gpt

import (
	"fmt"

	"github.com/felipeek/brasileirao-simulation/internal/util"
)

type messageCategory string

const (
	GPT_CONTEXT_MESSAGE = "You are being used in the simulation of Brazilian Soccer Championship (Brasileirao).\n" +
		"You are being invoked after a tournament round and your job is to generate a random event that will:\n" +
		"Impact the team [%s] by adding the value [%c%f] to the attribute [%s].\n" +
		"Attribute description: [%s]\n" +
		"If the received value is POSITIVE, the event MUST BE POSITIVE, otherwise it MUST BE NEGATIVE.\n" +
		"The category of the message is [%s]. (Do not explicitly mention the category in the response)\n" +
		"Also make sure that the message does not conflict with the real characteristics of these teams (they are real teams).\n" +
		"Also do not assume match results, nor the standings of the tournament, as you don't know that.\n" +
		"Note that a small value (i.e. close to 0) means a not so significant event, whereas a big value (e.g. close to 5) means a very significant event!\n" +
		"\n" +
		"The response must be at maximum 3 sentences. Keep it short.\n"
)

const (
	MESSAGE_CATEGORY_FUNNY         messageCategory = "FUNNY"
	MESSAGE_CATEGORY_CONTROVERSIAL messageCategory = "CONTROVERSIAL"
	MESSAGE_CATEGORY_INJURY        messageCategory = "MEDICAL_DEPARTMENT"
)

func GptRetrieveMessage(apiKey string, teamName string, attributeName string, attributeDescription string, valueDiff float64) (string, error) {
	messageCategory := util.RandomChoice(MESSAGE_CATEGORY_CONTROVERSIAL, MESSAGE_CATEGORY_FUNNY, MESSAGE_CATEGORY_INJURY).(messageCategory)

	signal := '+'
	if valueDiff < 0 {
		signal = '-'
	}
	fullMessage := fmt.Sprintf(GPT_CONTEXT_MESSAGE, teamName, signal, valueDiff, attributeName, attributeDescription, messageCategory)
	gptMessage, err := GptApiCall(apiKey, fullMessage)

	return gptMessage, err
}
