package urlshort

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
// func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
func MakeHandler(dataFilePath string, fallback http.Handler) (http.HandlerFunc, error) {
	pathsToUrls := make(map[string]string)

	// 1. parse yaml or json
	// parsing data file extension
	reYaml, _ := regexp.Compile(`\.yaml$`)
	reJson, _ := regexp.Compile(`\.json$`)

	// set data file path to lower case
	normalizedDataPath := strings.ToLower(dataFilePath)

	switch {

	case reYaml.MatchString(normalizedDataPath):
		yamlParsed, err := parseYaml(dataFilePath)
		if err != nil {
			return nil, err
		}

		// 2. build map for yaml
		pathsToUrls = buildMapYaml(yamlParsed)

	case reJson.MatchString(normalizedDataPath):
		jsonParsed, err := parseJson(dataFilePath)
		if err != nil {
			return nil, err
		}

		// 2. build map for json
		pathsToUrls = buildMapJson(jsonParsed)

	default:
		log.Fatal("wrong data file extension; must be 'yaml' or 'json'")
	}

	// 3. return map handler
	return MapHandler(pathsToUrls, fallback), nil
}

type pathUrlYaml struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

type pathUrlJson struct {
	Path string `json:"path"`
	Url  string `js:"url"`
}

func parseYaml(filepath string) ([]pathUrlYaml, error) {
	var result []pathUrlYaml

	file, err := os.ReadFile(filepath)
	if err != nil {
		return result, fmt.Errorf("failed to open yaml file: \n\t%v", err)
	}

	errY := yaml.Unmarshal(file, &result)
	if errY != nil {
		return result, fmt.Errorf("failed to unmarshal yaml data: \n\t%v", errY)
	}

	return result, nil
}

func parseJson(filepath string) ([]pathUrlJson, error) {
	var result []pathUrlJson

	file, err := os.ReadFile(filepath)
	if err != nil {
		return result, fmt.Errorf("failed to open json file: \n\t%v", err)
	}

	errJ := json.Unmarshal(file, &result)
	if errJ != nil {
		return result, fmt.Errorf("failed to unmarshal json data: \n\t%v", errJ)
	}

	return result, nil
}

func buildMapYaml(yamlParsed []pathUrlYaml) map[string]string {
	result := make(map[string]string)

	for _, v := range yamlParsed {
		result[v.Path] = v.Url
	}

	return result
}

func buildMapJson(jsonParsed []pathUrlJson) map[string]string {
	result := make(map[string]string)

	for _, v := range jsonParsed {
		result[v.Path] = v.Url
	}

	return result
}
