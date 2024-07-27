package gpt

import (
	"fmt"

	"github.com/felipeek/brasileirao-simulation/internal/util"
)

type MessageCategory string

type AttributeType struct {
	Name        string
	Description string
}

const (
	GPT_CONTEXT_MESSAGE = "You are being used in the simulation of Brazilian Soccer Championship (Brasileirao).\n" +
		"You are being invoked after round [%d] and your job is to generate a random event that will:\n" +
		"Impact the team [%s] by adding the value [%f] to the attribute [%s].\n" +
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
	MORALE_DESCRIPTION             = "The morale of the squad, ranging from 0 to 10."
	PHYSICAL_CONDITION_DESCRIPTION = "The physical condition of the squad, ranging from 0 to 10."
)

const (
	MESSAGE_CATEGORY_FUNNY         MessageCategory = "FUNNY"
	MESSAGE_CATEGORY_CONTROVERSIAL MessageCategory = "CONTROVERSIAL"
	MESSAGE_CATEGORY_INJURY        MessageCategory = "MEDICAL_DEPARTMENT"
)

var (
	MORALE_ATTRIBUTE = AttributeType{
		Name:        "MORALE",
		Description: MORALE_DESCRIPTION,
	}

	PHYSICAL_CONDITION_ATTRIBUTE = AttributeType{
		Name:        "PHYSICAL_CONDITION",
		Description: PHYSICAL_CONDITION_DESCRIPTION,
	}
)

func GptRetrieveMessage(apiKey string, teamName string, roundNum int) (string, float64, string, error) {
	attributeType := util.UtilRandomChoice(MORALE_ATTRIBUTE, PHYSICAL_CONDITION_ATTRIBUTE).(AttributeType)
	messageCategory := util.UtilRandomChoice(MESSAGE_CATEGORY_CONTROVERSIAL, MESSAGE_CATEGORY_FUNNY, MESSAGE_CATEGORY_INJURY).(MessageCategory)
	valueDiff := util.SimUtilRandomValueFromNormalDistribution(0.0, 4.0)

	fullMessage := fmt.Sprintf(GPT_CONTEXT_MESSAGE, roundNum, teamName, valueDiff, attributeType.Name, attributeType.Description, messageCategory)
	gptMessage, err := GptApiCall(apiKey, fullMessage)

	return attributeType.Name, valueDiff, gptMessage, err
}
