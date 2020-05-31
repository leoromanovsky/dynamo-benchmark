package dynamo

import (
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func NewDynamoDB(tableName string) *DynamoDB {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}

	// Set the AWS Region that the service clients should use
	cfg.Region = "us-east-1"

	// Using the Config value, create the DynamoDB client
	svc := dynamodb.New(cfg)

	/*
		req := svc.DescribeTableRequest(&dynamodb.DescribeTableInput{
			TableName: aws.String("Movies"),
		})

		// Send the request, and get the response or error back
		resp, err := req.Send(context.Background())
		if err != nil {
			panic("failed to describe table, "+err.Error())
		}
		fmt.Println("Response", resp)
	*/
	return &DynamoDB{
		tableName: tableName,
		client: svc,
	}
}

type DynamoDB struct {
	tableName string
	client *dynamodb.Client
}
