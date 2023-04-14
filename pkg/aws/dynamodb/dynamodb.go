package dynamodb

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Client struct {
	client *dynamodb.Client
}

func (svc Client) Get(ctx context.Context, tableName, keyCol, valCol, key string, out interface{}) error {
	result, err := svc.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			keyCol: &types.AttributeValueMemberS{Value: key},
		},
		ProjectionExpression: aws.String(fmt.Sprintf("%s,%s", keyCol, valCol)),
	})
	if err != nil {
		return fmt.Errorf("failed to get item from dynamodb: %w", err)
	}
	return attributevalue.NewDecoder().Decode(result.Item[valCol], out)
}

func (svc Client) BatchGet(ctx context.Context, tableName, keyCol, valCol string, keys []string, out interface{}) error {
	ks := []map[string]types.AttributeValue{}
	for _, key := range keys {
		ks = append(ks, map[string]types.AttributeValue{
			keyCol: &types.AttributeValueMemberS{Value: key},
		})
	}
	reqItems := map[string]types.KeysAndAttributes{
		tableName: {
			Keys:                 ks,
			ProjectionExpression: aws.String(fmt.Sprintf("%s,%s", keyCol, valCol)),
		},
	}
	res := make(map[string]types.AttributeValue)
	for len(reqItems) > 0 {
		result, err := svc.client.BatchGetItem(ctx, &dynamodb.BatchGetItemInput{
			RequestItems: reqItems,
		})
		if err != nil {
			return fmt.Errorf("failed to get items from dynamodb: %w", err)
		}
		items := result.Responses[tableName]
		for _, item := range items {
			var key string
			err := attributevalue.NewDecoder().Decode(item[keyCol], &key)
			if err != nil {
				return fmt.Errorf("unable to decode key from response for %v: %w", item, err)
			}
			res[key] = item[valCol]
		}
		reqItems = result.UnprocessedKeys
	}
	return attributevalue.UnmarshalMap(res, out)
}

func New(cfg aws.Config) *Client {
	return &Client{dynamodb.NewFromConfig(cfg)}
}
