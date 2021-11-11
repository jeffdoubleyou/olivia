package training

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/jeffdoubleyou/olivia/analysis"
	"github.com/jeffdoubleyou/olivia/network"
	"github.com/jeffdoubleyou/olivia/util"
)

// TrainData returns the inputs and outputs for the neural network
func TrainData(locale string) (inputs, outputs [][]float64) {
	words, classes, documents := analysis.Organize(locale)

	for _, document := range documents {
		outputRow := make([]float64, len(classes))
		bag := document.Sentence.WordsBag(words)

		// Change value to 1 where there is the document Tag
		outputRow[util.Index(classes, document.Tag)] = 1

		// Append data to inputs and outputs
		inputs = append(inputs, bag)
		outputs = append(outputs, outputRow)
	}

	return inputs, outputs
}

// CreateNeuralNetwork returns a new neural network which is loaded from res/training.json or
// trained from TrainData() inputs and targets.
func CreateNeuralNetwork(locale string, ignoreTrainingFile bool) (neuralNetwork network.Network) {
	// Decide if the network is created by the save or is a new one
	/*saveFile := "res/locales/" + locale + "/training.json"

	_, err := os.Open(saveFile)*/
	net := network.LoadNetwork(locale)

	// Train the model if there is no training file
	if net == nil || ignoreTrainingFile {
		inputs, outputs := TrainData(locale)

		neuralNetwork = network.CreateNetwork(locale, 0.1, inputs, outputs, 50)
		neuralNetwork.Train(200)
		neuralNetwork.Save()
	} else {
		fmt.Printf(
			"%s %s\n",
			color.FgBlue.Render("Loading the neural network from"),
			color.FgRed.Render(locale),
		)
		// Initialize the intents
		analysis.SerializeIntents(locale)
		neuralNetwork = *net
	}

	return
}
