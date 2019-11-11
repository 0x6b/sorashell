package shell

import (
	"bytes"
	"fmt"
	"github.com/rakyll/statik/fs"
	sl "github.com/soracom/soracom-cli/generators/lib"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

// loadAPIDef loads API definitions from the specified file
// based on https://github.com/soracom/soracom-cli/blob/master/generators/lib/apidef_loader.go
func loadAPIDef(apiDefYAMLFile string) (*sl.APIDefinitions, error) {
	apiDefYAML, err := loadAPIDefYAML(apiDefYAMLFile)
	if err != nil {
		return nil, err
	}

	apiDefMap := make(map[interface{}]interface{})
	err = yaml.Unmarshal(bytes.NewBufferString(apiDefYAML).Bytes(), &apiDefMap)
	if err != nil {
		return nil, err
	}

	methods, err := loadMethods(apiDefMap)
	if err != nil {
		return nil, err
	}

	return &sl.APIDefinitions{
		Host:     apiDefMap["host"].(string),
		BasePath: apiDefMap["basePath"].(string),
		Methods:  methods,
	}, nil
}

func loadAPIDefYAML(inputFile string) (string, error) {
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}
	f, err := statikFS.Open(inputFile)
	if err != nil {
		return "", err
	}

	defer func() {
		err := f.Close()
		if err != nil {
			fmt.Printf("warning: unable to close file %s", inputFile)
		}
	}()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}
	return bytes.NewBuffer(data).String(), nil
}

func loadMethods(apiDefMap map[interface{}]interface{}) ([]sl.APIMethod, error) {
	result := make([]sl.APIMethod, 0, len(apiDefMap))
	paths := apiDefMap["paths"].(map[interface{}]interface{})
	for path, p := range paths {
		methods := p.(map[interface{}]interface{})
		for method, m := range methods {
			mi, err := decodeAPIMethod(m.(map[interface{}]interface{}))
			if err != nil {
				return nil, err
			}
			mi.Path = path.(string)
			mi.Method = method.(string)
			result = append(result, *mi)
		}
	}
	return result, nil
}

func decodeAPIMethod(data map[interface{}]interface{}) (*sl.APIMethod, error) {
	y, err := yaml.Marshal(&data)
	if err != nil {
		return nil, err
	}
	var result sl.APIMethod
	err = yaml.Unmarshal(y, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
