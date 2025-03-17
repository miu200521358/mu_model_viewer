package usecase

import "os"

func LoadModelPath() (modelPath string, err error) {
	if len(os.Args) <= 1 {
		return "", nil
	}

	modelPath = os.Args[1]
	return modelPath, nil
}
