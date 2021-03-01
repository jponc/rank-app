package rankings

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jponc/rank-app/api"
	log "github.com/sirupsen/logrus"
)

// ProcessKeyword processes one keyword
func (s *service) ProcessKeyword(ctx context.Context, snsEvent events.SNSEvent) {
	snsMsg := snsEvent.Records[0].SNS.Message

	var processKeywordMsg api.ProcessKeywordMessage

	err := json.Unmarshal([]byte(snsMsg), &processKeywordMsg)
	if err != nil {
		log.Fatalf("unable to unarmarshal message: %v", err)
	}

	res, err := s.zenserpClient.Search(ctx, processKeywordMsg.Keyword, processKeywordMsg.Count)
	if err != nil {
		log.Fatalf("Unable to to query data from zenserp using keyword: %s", processKeywordMsg.Keyword)
	}

	log.Info("Result: %v", res)
}
