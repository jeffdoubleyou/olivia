package locales

import (
	"context"
	"fmt"
	"github.com/jeffdoubleyou/olivia/database"
	_ "github.com/jeffdoubleyou/olivia/res/locales/c1212"
	// Import these packages to trigger the init() function
	_ "github.com/jeffdoubleyou/olivia/res/locales/ca"
	_ "github.com/jeffdoubleyou/olivia/res/locales/de"
	_ "github.com/jeffdoubleyou/olivia/res/locales/el"
	_ "github.com/jeffdoubleyou/olivia/res/locales/en"
	_ "github.com/jeffdoubleyou/olivia/res/locales/es"
	_ "github.com/jeffdoubleyou/olivia/res/locales/fr"
	_ "github.com/jeffdoubleyou/olivia/res/locales/it"
	_ "github.com/jeffdoubleyou/olivia/res/locales/nl"
	_ "github.com/jeffdoubleyou/olivia/res/locales/tr"
)

var Locales []Locale

func init() {
	fmt.Printf("Initializing Locales..........\n")
	// http://localhost:5984/intents/_design/locales/_view/list?limit=20&reduce=true&group_level=1
	if rows, err := database.Db("intents").Query(context.TODO(), "_design/locales", "_view/list", map[string]interface{}{"group": true}); err != nil {
		panic(err)
	} else {
		for rows.Next() {
			var l string
			if err := rows.ScanKey(&l); err != nil {
				fmt.Printf("Failed to read row: %s\n", err.Error())
			} else {
				// TODO Register modules
				Locales = append(Locales, Locale{l, l})
			}
		}
	}
	for _, locale := range Locales {
		fmt.Printf("Found locale: %s\n", locale.Name)
	}
}

// Locales is the list of locales's tags and names
// Please check if the language is supported in https://github.com/tebeka/snowball,
// if it is please add the correct language name.

// A Locale is a registered locale in the file
type Locale struct {
	Tag  string
	Name string
}

// GetNameByTag returns the name of the given locale's tag
func GetNameByTag(tag string) string {
	for _, locale := range Locales {
		if locale.Tag != tag {
			continue
		}

		return locale.Name
	}

	return ""
}

// GetTagByName returns the tag of the given locale's name
func GetTagByName(name string) string {
	for _, locale := range Locales {
		if locale.Name != name {
			continue
		}

		return locale.Tag
	}

	return ""
}

// Exists checks if the given tag exists in the list of locales
func Exists(tag string) bool {
	for _, locale := range Locales {
		if locale.Tag == tag {
			return true
		}
	}

	return false
}
