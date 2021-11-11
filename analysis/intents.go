package analysis

import (
	"sort"

	"github.com/jeffdoubleyou/olivia/modules"
	"github.com/jeffdoubleyou/olivia/util"
)

// Intent is a way to group sentences that mean the same thing and link them with a tag which
// represents what they mean, some responses that the bot can reply and a context
type Intent struct {
	Id        string                 `json:"_id"`
	Tag       string                 `json:"tag"`
	Patterns  []string               `json:"patterns"`
	Responses []string               `json:"responses"`
	Context   string                 `json:"context"`
	Data      map[string]interface{} `json:"data"`
	Locale    string                 `json:"locale"`
	Language  string                 `json:"language"`
}

// Document is any sentence from the intents' patterns linked with its tag
type Document struct {
	Sentence Sentence
	Tag      string
}

var intents = map[string][]Intent{}

// CacheIntents set the given intents to the global variable intents
func CacheIntents(locale string, _intents []Intent) {
	intents[locale] = _intents
}

// GetIntents returns the cached intents
func GetIntents(locale string) []Intent {
	return intents[locale]
}

// SerializeIntents returns a list of intents retrieved from the given intents file
func SerializeIntents(locale string) (_intents []Intent) {
	if _intents, err := LoadIntents(locale); err != nil {
		panic(err)
	} else {
		CacheIntents(locale, _intents)
		return _intents
	}
}

// SerializeModulesIntents retrieves all the registered modules and returns an array of Intents
func SerializeModulesIntents(locale string) []Intent {
	registeredModules := modules.GetModules(locale)
	intents := make([]Intent, len(registeredModules))

	for k, module := range registeredModules {
		intents[k] = Intent{
			Tag:       module.Tag,
			Patterns:  module.Patterns,
			Responses: module.Responses,
			Context:   "",
		}
	}

	return intents
}

// GetIntentByTag returns an intent found by given tag and locale
func GetIntentByTag(tag, locale string) Intent {
	for _, intent := range GetIntents(locale) {
		if tag != intent.Tag {
			continue
		}

		return intent
	}

	return Intent{}
}

// Organize intents with an array of all words, an array with a representative word of each tag
// and an array of Documents which contains a word list associated with a tag
func Organize(locale string) (words, classes []string, documents []Document) {
	// Append the modules intents to the intents from res/datasets/intents.json
	intents := append(
		SerializeIntents(locale),
		SerializeModulesIntents(locale)...,
	)

	for _, intent := range intents {
		for _, pattern := range intent.Patterns {
			// Tokenize the pattern's sentence
			patternSentence := Sentence{locale, pattern}
			patternSentence.arrange()

			// Add each word to response
			for _, word := range patternSentence.stem() {

				if !util.Contains(words, word) {
					words = append(words, word)
				}
			}

			// Add a new document
			documents = append(documents, Document{
				patternSentence,
				intent.Tag,
			})
		}

		// Add the intent tag to classes
		classes = append(classes, intent.Tag)
	}

	sort.Strings(words)
	sort.Strings(classes)

	return words, classes, documents
}
