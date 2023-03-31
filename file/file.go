package file

import (
	"os"
)

func CrateFile(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err = os.Mkdir(filePath, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
