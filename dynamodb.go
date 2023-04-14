package props

import (
	"github.com/SumoLogic-Labs/props/pkg/aws/dynamodb"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/magiconair/properties"
)

type DynamoDBGetterArgs struct {
	Table     string
	KeyCol    string
	ValCol    string
	WatchKeys []string
}

type DynamoDBGetter struct {
	*dynamodb.Client
	args DynamoDBGetterArgs
}

func (s DynamoDBGetter) Poll(ctx context.Context) (*properties.Properties, error) {
	m := make(map[string]string)
	err := s.Client.BatchGet(ctx, s.args.Table, s.args.KeyCol, s.args.ValCol, s.args.WatchKeys, &m)
	if err != nil {
		return nil, fmt.Errorf("unable to get items from table: %w", err)
	}
	return properties.LoadMap(m), nil
}

func NewDynamoDBGetterSource(cfg aws.Config, args DynamoDBGetterArgs) *DynamoDBGetter {
	return &DynamoDBGetter{
		Client: dynamodb.New(cfg),
		args:   args,
	}
}
