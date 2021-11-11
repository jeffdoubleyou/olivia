package network

import (
	"context"
	"fmt"
	"github.com/jeffdoubleyou/olivia/database"
)

func AddNetworkToDB(locale string, data Network) error {
	if exists, _ := GetNetworkFromDb(locale); exists != nil {
		data.Rev = exists.Rev
	}
	if rev, err := database.Db("networks").Put(context.TODO(), locale, &data); err != nil {
		return err
	} else {
		fmt.Printf("Inserted network record with rev %s\n", rev)
		return nil
	}
	return nil
}

func GetNetworkFromDb(locale string) (network *Network, err error) {
	fmt.Printf("Get network for locale: %s\n", locale)
	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"Locale": locale,
		},
	}
	if rows, err := database.Db("networks").Find(context.TODO(), query); err != nil {
		fmt.Printf("Did not find a network for locale '%s': %s", locale, err.Error())
		return nil, err
	} else {
		fmt.Printf("Got a network for '%s' - rows: %d\n", locale, rows.TotalRows())
		for rows.Next() {
			if err := rows.ScanDoc(&network); err != nil {
				fmt.Printf("Could not read document into network: %s\n", err.Error())
				return nil, err
			} else {
				network.Rev = rows.ID()
				fmt.Printf("Scanned network for %s with ID '%s'\n", locale, network.Rev)
			}
		}
	}
	return
}
