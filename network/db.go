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
	if row := database.Db("networks").Get(context.TODO(), locale); row == nil {
		return nil, fmt.Errorf("did not find a network for locale '%s'", locale)
	} else {
		err = row.ScanDoc(&network)
		if err != nil {
			fmt.Printf("Could not read document into network: %s\n", err.Error())
		}
		return
	}
}
