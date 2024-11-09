package files

import (
	"errors"
	"fmt"
	"github.com/go-clarum/agent/validators/strings"
	"gopkg.in/yaml.v3"
	"os"
)

func ReadYamlFileToStruct[S any](filePath string) (*S, error) {
	if strings.IsBlank(filePath) {
		return nil, errors.New("unable to read file - file path is empty")
	}

	buf, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to read file - %s", err))
	}

	out := new(S)

	if err := yaml.Unmarshal(buf, out); err != nil {
		return nil, errors.New(fmt.Sprintf("failed to unmarshal yaml file %s: %s", filePath, err))
	}

	return out, err
}
