package sns

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	awsSns "github.com/aws/aws-sdk-go/service/sns"
	log "github.com/sirupsen/logrus"
)

// Client interface
type Client interface {
	// Publish publishes a new message to a SNS topic
	Publish(topic string, message interface{}) error
}

type client struct {
	awsSnsClient *awsSns.SNS
	snsPrefix    string
}

// NewClient instantiates a SNS client
func NewClient(awsRegion, snsPrefix string) (Client, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})

	if err != nil {
		return nil, fmt.Errorf("cannot create aws session: %v", err)
	}

	awsSnsClient := sns.New(sess)

	c := &client{
		awsSnsClient: awsSnsClient,
		snsPrefix:    snsPrefix,
	}

	return c, nil
}

func (c *client) Publish(topic string, message interface{}) error {
	msg, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal sns message: %v", err)
	}

	input := &awsSns.PublishInput{
		Message:  aws.String(string(msg)),
		TopicArn: aws.String(fmt.Sprintf("%s-%s", c.snsPrefix, topic)),
	}

	result, err := c.awsSnsClient.Publish(input)
	if err != nil {
		return fmt.Errorf("failed to publish to sns: %v", err)
	}

	log.Info("successfully processed sns message: (%s) for topic (%s)", result, topic)
	return nil
}
