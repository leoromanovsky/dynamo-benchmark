package benchmark

import "fmt"

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
