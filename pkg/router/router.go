package router

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/korzepadawid/aws-lambda-dynamo/pkg/posts"
	"github.com/korzepadawid/aws-lambda-dynamo/pkg/util"
)

type Router struct {
	DynDB *dynamodb.Client
}

func NewRouter(dynDB *dynamodb.Client) *Router {
	return &Router{
		DynDB: dynDB,
	}
}

func (r *Router) Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	postSvc := posts.NewPostService(request, r.DynDB)
	switch request.HTTPMethod {
	case http.MethodGet:
		return postSvc.Get()
	case http.MethodPost:
		return postSvc.Create()
	case http.MethodPut:
		return postSvc.Update()
	case http.MethodDelete:
		return postSvc.Delete()
	default:
		return util.ResponseWithError(
			http.StatusMethodNotAllowed,
			fmt.Errorf("method %s not allowed", request.HTTPMethod),
		), nil
	}
}
