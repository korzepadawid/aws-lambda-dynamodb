package posts

import (
	"encoding/base64"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

func deserialize(request events.APIGatewayProxyRequest) (*Post, error) {
	body := request.Body

	if request.IsBase64Encoded {
		base64Decoded, err := base64.StdEncoding.DecodeString(body)
		if err != nil {
			return nil, err
		}
		body = string(base64Decoded)
	}

	p := &Post{}
	if err := json.Unmarshal([]byte(body), p); err != nil {
		return nil, err
	}

	return p, nil
}

func updatePost(p *Post, changes *Post) {
	if changes.Title != "" {
		p.Title = changes.Title
	}
	if changes.Body != "" {
		p.Body = changes.Body
	}
	if changes.UserID != 0 {
		p.UserID = changes.UserID
	}
}
