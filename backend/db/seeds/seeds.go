package seeds

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Apply seed with version
func Apply(instance *sql.DB, version uint) error {
	name := fmt.Sprintf("%d_seed.json", version)
	relativePath := filepath.Join("db", "seeds", name)
	path, _ := filepath.Abs(relativePath)
	dataFile, err := os.Open(path)

	if err != nil {
		return err
	}

	defer dataFile.Close()
	var data map[string][]interface{}

	byteValue, _ := ioutil.ReadAll(dataFile)
	err = json.Unmarshal(byteValue, &data)

	if err != nil {
		return err
	}

	steps, ok := data["db"]

	if !ok {
		return nil
	}

	for _, step := range steps {
		step := step.(map[string]interface{})
		sql := step["sql"].(string)
		items := step["items"].([]interface{})

		for _, item := range items {
			nItem := item.([]interface{})
			_, err := instance.Exec(sql, nItem...)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ApplyAll applies all seeds up to version
func ApplyAll(instance *sql.DB, version uint) error {
	var i uint
	for i = 1; i <= version; i++ {
		err := Apply(instance, i)

		if err != nil {
			return err
		}
	}

	return nil
}
