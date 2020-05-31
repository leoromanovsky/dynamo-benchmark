package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
	cli "github.com/urfave/cli/v2"
	"github.com/leoromanovsky/dynamo-benchmark/internal/benchmark"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

/*
This program performs various benchmarks on DynamoDB tables
	//"github.com/schollz/progressbar/v2"
 */

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

var svc *dynamodb.Client

func runBenchmark() {
	fmt.Print("\n--- \033[1;32mBENCHMARK\033[0m ----------------------------------------------------------------------------------------------------------------\n\n")
	printHeader()
	fmt.Print("+---------+----------------+------------------------------------------------+------------------------------------------------+\n\n")
}

// prints the table header for the test results
func printHeader() {
	// print the table header
	fmt.Printf("Download performance with \033[1;33m%-s\033[0m objects%s\n", "foo", "bar")
	fmt.Println("                           +-------------------------------------------------------------------------------------------------+")
	fmt.Println("                           |            Time to First Byte (ms)             |            Time to Last Byte (ms)              |")
	fmt.Println("+---------+----------------+------------------------------------------------+------------------------------------------------+")
	fmt.Println("|       # |     Throughput |  avg   min   p25   p50   p75   p90   p99   max |  avg   min   p25   p50   p75   p90   p99   max |")
	fmt.Println("+---------+----------------+------------------------------------------------+------------------------------------------------+")
}

func setupDynamoDBClient() {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}

	// Set the AWS Region that the service clients should use
	cfg.Region = "us-east-1"

	// Using the Config value, create the DynamoDB client
	svc = dynamodb.New(cfg)

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
}

func createDynamoTable(tableName string, svc *dynamodb.Client) {
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
		TableName: aws.String(tableName),
	}

	req := svc.CreateTableRequest(input)
	_, err := req.Send(context.Background())
	if err != nil {
		fmt.Println("Got error calling CreateTable:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	time.Sleep(10 * time.Second)

	fmt.Println("Created the table", tableName)
}

func insertRow(tableName string, svc *dynamodb.Client) {
	item := benchmark.Item{
		Year:   2015,
		Title:  randSeq(20),
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
		TableName: aws.String(tableName),
	}

	req := svc.PutItemRequest(input)
	_, err = req.Send(context.Background())
	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	year := strconv.Itoa(item.Year)

	fmt.Println("Successfully added '" + item.Title + "' (" + year + ") to table " + tableName)
}

func deleteTable(tableName string, svc *dynamodb.Client) {
	input := &dynamodb.DeleteTableInput{TableName: aws.String(tableName)}
	req := svc.DeleteTableRequest(input)
	_, err := req.Send(context.Background())
	if err != nil {
		fmt.Println("Got error calling DeleteTable:")
		fmt.Println(err.Error())
	}
}

const tableNameBase = "Movies"

func tableName() string {
	return fmt.Sprintf("%s-dev", tableNameBase)
}

// flag to just cleanup test artifacts
var cleanupOnly bool
// flag to create table if it does not exist
var createTable bool = true

// program entry point
func main() {
	rand.Seed(time.Now().UnixNano())

	// create client
	setupDynamoDBClient()

	// parse flags
	app := &cli.App{
		Name: "dynamo-benchmark",
		Version:     "0.1.0",
		Description: "This is how we describe greet the app",
		Usage: "make an explosive entrance",
		Flags: []cli.Flag {
			&cli.BoolFlag{
				Name:  "cleanup-only",
				Usage: "Cleanup the test objects and exit",
				Destination: &cleanupOnly,
			},
			&cli.BoolFlag{
				Name:  "create-table",
				Usage: "Create the table if it does not exist",
				Value: true,
				Destination: &createTable,
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "generate-test",
				Usage:   "generate test file",
				Action: func(c *cli.Context) error {
					// create random data and write it out to a file
					fmt.Println("added task: ", c.Args().First())
					return nil
				},
			},
			{
				Name:    "run-test",
				Usage:   "add a task to the list",
				Action: func(c *cli.Context) error {
					fmt.Println("boom! I say!")

					if cleanupOnly {
						cleanup()
						return nil
					}

					runTest()

					return nil
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func cleanup() {
	deleteTable(tableName(), svc)
}

func setup() {
	// create table
	if createTable {
		createDynamoTable(tableName(), svc)
	}

	// populate table
	for i := 1; i <= 10; i++ {
		insertRow(tableName(), svc)
	}
}

func runTest() {
	// bootstrap the dynamodb table and upload test data
	setup()

	// run the test
	runBenchmark()

	// cleanup
	cleanup()
}
