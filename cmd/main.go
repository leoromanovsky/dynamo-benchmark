package main

import (
	"fmt"
	"github.com/leoromanovsky/dynamo-benchmark/internal/dynamo"
	cli "github.com/urfave/cli/v2"
	"log"
	"math/rand"
	"os"
	"time"
)

/*
This program performs various benchmarks on DynamoDB tables
	//"github.com/schollz/progressbar/v2"
 */

var dynamoService *dynamo.DynamoDB

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
	dynamoService = dynamo.NewDynamoDB(tableName())

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
					fmt.Println("Generating test data.")
					return nil
				},
			},
			{
				Name:    "run-test",
				Usage:   "add a task to the list",
				Action: func(c *cli.Context) error {
					fmt.Println("Starting benchmark test.")

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
	dynamoService.DeleteTable()
}

func setup() {
	// create table
	if createTable {
		dynamoService.CreateTable()
	}

	// populate table
	for i := 1; i <= 10; i++ {
		dynamoService.InsertRow()
	}
}

func runTest() {
	// bootstrap the dynamodb table and upload test data
	setup()

	// run the test
	//runBenchmark()

	// cleanup
	cleanup()
}
