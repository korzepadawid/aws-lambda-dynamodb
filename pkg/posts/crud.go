package posts

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/google/uuid"
	"github.com/korzepadawid/aws-lambda-dynamo/pkg/util"
)

const (
	TableName              = "Posts"
	DynamoDBDefaultTimeout = 10 * time.Second
)

type Post struct {
	ID     string `json:"id,omitempty" dynamodbav:"id"`
	Title  string `json:"title,omitempty" dynamodbav:"title"`
	Body   string `json:"body,omitempty" dynamodbav:"body"`
	UserID int    `json:"userId,omitempty" dynamodbav:"userId"`
}

type PostService struct {
	Request  events.APIGatewayProxyRequest
	DynamoDB *dynamodb.Client
}

func NewPostService(req events.APIGatewayProxyRequest, dyndb *dynamodb.Client) *PostService {
	return &PostService{
		DynamoDB: dyndb,
		Request:  req,
	}
}

func (s *PostService) Get() (events.APIGatewayProxyResponse, error) {
	return util.ResponseWithBody(http.StatusOK, Post{Title: "not implemented"}), nil
}

func (s *PostService) Create() (events.APIGatewayProxyResponse, error) {
	post, err := deserialize(s.Request)
	if err != nil {
		return util.ResponseWithError(http.StatusBadRequest, err), nil
	}
	post.ID = uuid.New().String()

	pv, err := attributevalue.MarshalMap(post)
	if err != nil {
		return util.ResponseWithError(http.StatusInternalServerError, fmt.Errorf("error when marshalling new post item: %w", err)), nil
	}

	input := &dynamodb.PutItemInput{
		Item:      pv,
		TableName: aws.String(TableName),
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), DynamoDBDefaultTimeout)
	defer cancelFn()

	if _, err = s.DynamoDB.PutItem(ctx, input); err != nil {
		return util.ResponseWithError(http.StatusInternalServerError, fmt.Errorf("error when creating new post item: %w", err)), nil
	}

	return util.ResponseWithBody(http.StatusCreated, post), nil
}

func (s *PostService) Update() (events.APIGatewayProxyResponse, error) {
	return util.ResponseWithBody(http.StatusOK, Post{Title: "not implemented"}), nil
}

func (s *PostService) Delete() (events.APIGatewayProxyResponse, error) {
	return util.ResponseWithStatusCode(http.StatusNoContent), nil
}
