package compiler

import (
	"encoding/json"
	"io/ioutil"
)

func writeToFile(stuff interface{}, pathOfFile string) error {
	var stuffJSON, err = json.Marshal(stuff)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(pathOfFile, stuffJSON, 0666)
}
