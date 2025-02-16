package util

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type Error struct {
	Message string `json:"message,omitempty"`
	Status  int    `json:"status,omitempty"`
}

func ResponseWithBody(status int, body interface{}) events.APIGatewayProxyResponse {
	b, err := json.Marshal(body)

	if err != nil {
		return ResponseWithStatusCode(http.StatusInternalServerError)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       string(b),
	}
}

func ResponseWithStatusCode(status int) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
	}
}

func ResponseWithError(status int, err error) events.APIGatewayProxyResponse {
	return ResponseWithBody(status, Error{
		Status:  status,
		Message: err.Error()},
	)
}
