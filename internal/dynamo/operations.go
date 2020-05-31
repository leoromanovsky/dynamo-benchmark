package dynamo

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
	"github.com/leoromanovsky/dynamo-benchmark/internal/benchmark"
	"github.com/leoromanovsky/dynamo-benchmark/internal/utils"
	"os"
	"strconv"
	"time"
)

func (d *DynamoDB) CreateTable() {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Year"),
				AttributeType: dynamodb.ScalarAttributeTypeN,
			},
			{
				AttributeName: aws.String("Title"),
				AttributeType: dynamodb.ScalarAttributeTypeS,
			},
		},
		KeySchema: []dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Year"),
				KeyType:     dynamodb.KeyTypeHash,
			},
			{
				AttributeName: aws.String("Title"),
				KeyType:      dynamodb.KeyTypeRange,
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(d.tableName),
	}

	req := d.client.CreateTableRequest(input)
	_, err := req.Send(context.Background())
	if err != nil {
		fmt.Println("Got error calling CreateTable:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	time.Sleep(10 * time.Second)

	fmt.Println("Created the table", d.tableName)
}

func (d *DynamoDB) InsertRow() {
	item := benchmark.Item{
		Year:   utils.RandNumber(5),
		Title:  utils.RandSeq(20),
		Plot:   "Nothing happens at all.",
		Rating: 0.0,
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		fmt.Println("Got error marshalling new movie item:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(d.tableName),
	}

	req := d.client.PutItemRequest(input)
	_, err = req.Send(context.Background())
	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	year := strconv.Itoa(item.Year)

	fmt.Println("Successfully added '" + item.Title + "' (" + year + ") to table " + d.tableName)
}

func (d *DynamoDB) DeleteTable() {
	input := &dynamodb.DeleteTableInput{TableName: aws.String(d.tableName)}
	req := d.client.DeleteTableRequest(input)
	_, err := req.Send(context.Background())
	if err != nil {
		fmt.Println("Got error calling DeleteTable:")
		fmt.Println(err.Error())
	}
}
