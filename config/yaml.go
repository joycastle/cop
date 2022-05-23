package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/joycastle/cop/util"
	"gopkg.in/yaml.v2"
)

var (
	ErrCommonTemplete     = "[ReadYamlConfig] %s"
	ErrFileNotExists      = errors.New("[ReadYamlConfig] file not exists")
	ErrFileFormatNotMatch = errors.New("[ReadYamlConfig] file format not match need .yaml or .yml")
)

func ReadYmalFromFile(fileName string, out interface{}) error {
	if !util.IsVarPointor(out) {
		return fmt.Errorf(ErrCommonTemplete, "var out is not pointor or quote type")
	}

	if !util.FileExists(fileName) {
		return ErrFileNotExists
	}

	if !strings.HasSuffix(fileName, ".yaml") && !strings.HasSuffix(fileName, ".yml") {
		return ErrFileFormatNotMatch
	}

	fd, err := ioutil.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf(ErrCommonTemplete, err.Error)
	}

	if err := yaml.Unmarshal(fd, out); err != nil {
		return fmt.Errorf(ErrCommonTemplete, err.Error())
	}

	return nil
}
