package parse_test

import (
	"fmt"
	"os"

	"github.com/tmc/parse"
)

func ExampleNewClient() {
	appID := os.Getenv("APPLICATION_ID")
	apiKey := os.Getenv("REST_API_KEY")
	_, err := parse.NewClient(appID, apiKey)
	fmt.Println(err)
	// output: <nil>
}
