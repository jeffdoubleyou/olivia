package analysis

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jeffdoubleyou/olivia/database"
)

func LoadIntents(locale string) ([]Intent, error) {
	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"locale": locale,
		},
	}
	if rows, err := database.Db("intents").Find(context.TODO(), query); err != nil {
		return nil, err
	} else {
		var intents []Intent
		for rows.Next() {
			var intent Intent
			if err := rows.ScanDoc(&intent); err != nil {
				return nil, err
			} else {
				intents = append(intents, intent)
			}
		}
		return intents, nil
	}
}

func AddIntent(intent *Intent) error {
	if intent.Id == "" {
		intent.Id = uuid.New().String()
	}
	if rev, err := database.Db("intents").Put(context.TODO(), "intents", intent); err != nil {
		fmt.Printf("Failed to insert intent: %s\n", err.Error())
		return err
	} else {
		fmt.Printf("Inserted record with rev %s", rev)
		return nil
	}
}
