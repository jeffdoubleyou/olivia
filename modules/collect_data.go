package modules

import (
	"fmt"
	"github.com/jeffdoubleyou/olivia/language"
	"github.com/jeffdoubleyou/olivia/util"
	"github.com/soudy/mathcat"
	"regexp"
)

var CollectDataTag = "collect data"

func CollectDataReplacer(locale, entry, response, _ string) (string, string) {
	// Language should be - type:<int|string> name:<string> callback:<url>
	operation := language.FindMathOperation(entry)

	// If there is no operation in the entry message reply with a "don't understand" message
	if operation == "" {
		responseTag := "don't understand"
		return responseTag, util.GetMessage(locale, responseTag)
	}

	res, err := mathcat.Eval(operation)
	// If the expression isn't valid reply with a message from res/datasets/messages.json
	if err != nil {
		responseTag := "math not valid"
		return responseTag, util.GetMessage(locale, responseTag)
	}
	// Use number of decimals from the query
	decimals := language.FindNumberOfDecimals(locale, entry)
	if decimals == 0 {
		decimals = 6
	}

	result := res.FloatString(decimals)

	// Remove trailing zeros of the result with a Regex
	trailingZerosRegex := regexp.MustCompile(`\.?0+$`)
	result = trailingZerosRegex.ReplaceAllString(result, "")

	return MathTag, fmt.Sprintf(response, result)
}
