package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/jeffdoubleyou/olivia/locales"
	"github.com/jeffdoubleyou/olivia/training"

	"github.com/jeffdoubleyou/olivia/dashboard"

	"github.com/jeffdoubleyou/olivia/util"

	"github.com/gookit/color"

	"github.com/jeffdoubleyou/olivia/network"

	"github.com/jeffdoubleyou/olivia/server"
)

var neuralNetworks = map[string]network.Network{}

func main() {
	port := flag.String("port", "8080", "The port for the API and WebSocket.")
	localesFlag := flag.String("re-train", "", "The locale(s) to re-train.")
	flag.Parse()

	// If the locales flag isn't empty then retrain the given dao
	if *localesFlag != "" {
		reTrainModels(*localesFlag)
	}

	// Print the Olivia ascii text
	oliviaASCII := string(util.ReadFile("res/olivia-ascii.txt"))
	fmt.Println(color.FgLightGreen.Render(oliviaASCII))

	// Create the authentication token
	dashboard.Authenticate()

	for _, locale := range locales.Locales {
		util.SerializeMessages(locale.Tag)

		neuralNetworks[locale.Tag] = training.CreateNeuralNetwork(
			locale.Tag,
			false,
		)
	}

	// Get port from environment variables if there is
	if os.Getenv("PORT") != "" {
		*port = os.Getenv("PORT")
	}

	// Serves the server
	server.Serve(neuralNetworks, *port)
}

// reTrainModels retrain the given locales
func reTrainModels(localesFlag string) {
	fmt.Printf("Re-train models...")
	// Iterate locales by separating them by comma
	for _, localeFlag := range strings.Split(localesFlag, ",") {
		path := fmt.Sprintf("res/locales/%s/training.json", localeFlag)
		fmt.Printf("Remove '%s'\n", path)
		err := os.Remove(path)

		if err != nil {
			fmt.Printf("Cannot re-train %s model: %s", localeFlag, err.Error())
			return
		}
	}
}
