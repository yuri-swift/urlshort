package urlshort

import (
	"net/http"

	yaml "gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
// 返却値はhttp.HandlerFunc(http.Handlerを実装)
// http.HandlerFuncはパスを対応するURLにマップして返す
// マップ内にパスが指定されていない場合はfallbackとしてhttp.Headerが呼ばれる
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// リクエストからパスを取得
		path := r.URL.Path
		// リクエストされたパスからキーを見つけたら変数OKにはtrueが入る
		if dest, ok := pathsToUrls[path]; ok {
			//マッチするものが見つかればリダイレクトさせる
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
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yamlBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// yamlをパース
	pathUrls, err := parseYaml(yamlBytes)
	if err != nil {
		return nil, err
	}
	// yamlをmapに変換
	pathsToUrls := buildMap(pathUrls)
	// map handlerを返す
	return MapHandler(pathsToUrls, fallback), nil
}

func buildMap(pathUrls []pathUrl) map[string]string {
	fmt.Println("buildMap")
	pathsToUrls := make(map[string]string)
	for _, pu := range pathUrls {
		pathsToUrls[pu.Path] = pu.URL
	}
	return pathsToUrls
}

func parseYaml(data []byte) ([]pathUrl, error) {
	var pathUrls []pathUrl
	err := yaml.Unmarshal(data, &pathUrls)
	if err != nil {
		return nil, err
	}
	return pathUrls, nil
}

type pathUrl struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}