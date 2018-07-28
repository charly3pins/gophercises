package urlshort

import (
	"log"
	"net/http"

	yaml "gopkg.in/yaml.v2"
)

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if v, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, v, http.StatusFound)
			return
		}

		fallback.ServeHTTP(w, r)
	}
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathToURLs, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	pathToURLMap := buildMap(*pathToURLs)
	return MapHandler(pathToURLMap, fallback), nil
}

type PathToURL struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

func parseYAML(yml []byte) (*[]PathToURL, error) {
	var pathToURLs []PathToURL
	err := yaml.Unmarshal(yml, &pathToURLs)
	if err != nil {
		log.Println("Error unmarshalling yaml file", err)
		return nil, err
	}

	return &pathToURLs, nil
}

func buildMap(yamlPaths []PathToURL) map[string]string {
	mapPath := map[string]string{}
	for _, v := range yamlPaths {
		mapPath[v.Path] = v.URL
	}
	return mapPath
}
