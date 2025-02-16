package posts

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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
	requestedID := s.Request.PathParameters["id"]
	key := map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberS{Value: requestedID},
	}

	getItemInput := &dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key:       key,
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), DynamoDBDefaultTimeout)
	defer cancelFn()
	result, err := s.DynamoDB.GetItem(ctx, getItemInput)
	if err != nil {
		return util.ResponseWithError(http.StatusInternalServerError, fmt.Errorf("error when getting post item: %w", err)), nil
	}

	if result.Item == nil {
		return util.ResponseWithError(http.StatusNotFound, fmt.Errorf("post item not found")), nil
	}

	var post Post
	if err = attributevalue.UnmarshalMap(result.Item, &post); err != nil {
		return util.ResponseWithError(http.StatusInternalServerError, fmt.Errorf("error when unmarshalling post item: %w", err)), nil
	}

	return util.ResponseWithBody(http.StatusOK, post), nil
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
	updateID := s.Request.PathParameters["id"]

	inputPost, err := deserialize(s.Request)
	if err != nil {
		return util.ResponseWithError(http.StatusBadRequest, err), nil
	}

	key := map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberS{
			Value: updateID,
		},
	}

	getItemInput := &dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key:       key,
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), DynamoDBDefaultTimeout)
	defer cancelFn()

	out, err := s.DynamoDB.GetItem(ctx, getItemInput)
	if err != nil {
		return util.ResponseWithError(
			http.StatusInternalServerError,
			fmt.Errorf("error when getting post item: %w", err),
		), nil
	}

	if out.Item == nil {
		return util.ResponseWithError(http.StatusNotFound, fmt.Errorf("post item not found")), nil
	}

	var post Post
	if err = attributevalue.UnmarshalMap(out.Item, &post); err != nil {
		return util.ResponseWithError(http.StatusInternalServerError, fmt.Errorf("error when unmarshalling post item: %w", err)), nil
	}

	updatePost(&post, inputPost)

	pv, err := attributevalue.MarshalMap(post)
	if err != nil {
		return util.ResponseWithError(http.StatusInternalServerError, fmt.Errorf("error when marshalling new post item: %w", err)), nil
	}

	input := &dynamodb.PutItemInput{
		Item:      pv,
		TableName: aws.String(TableName),
	}

	if _, err = s.DynamoDB.PutItem(ctx, input); err != nil {
		return util.ResponseWithError(
			http.StatusInternalServerError,
			fmt.Errorf("error when creating new post item: %w", err),
		), nil
	}

	return util.ResponseWithBody(http.StatusOK, post), nil
}

func (s *PostService) Delete() (events.APIGatewayProxyResponse, error) {
	deleteID := s.Request.PathParameters["id"]

	key := map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberS{Value: deleteID},
	}

	dynamodbDeleteInput := &dynamodb.DeleteItemInput{
		TableName: aws.String(TableName),
		Key:       key,
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), DynamoDBDefaultTimeout)
	defer cancelFn()

	if _, err := s.DynamoDB.DeleteItem(ctx, dynamodbDeleteInput); err != nil {
		return util.ResponseWithError(
			http.StatusInternalServerError,
			fmt.Errorf("error when deleting post item: %w", err),
		), nil
	}

	return util.ResponseWithStatusCode(http.StatusNoContent), nil
}
