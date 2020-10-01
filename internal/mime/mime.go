package mime

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"service-watch/internal/models"
	"strings"
)

var (
	mimeBuffer   = make(map[string]models.MimeFormat)
	extensions   = make(map[string]string)
	regexExtract = regexp.MustCompile(`^\s*([^;\s]*)(?:;|\s|$)`)
)

func init() {
	root, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	path := root + "/internal/mime/mime-collection.json"

	mimeFile, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	mimeBytes, _ := ioutil.ReadAll(mimeFile)

	err = json.Unmarshal(mimeBytes, &mimeBuffer)
	if err != nil {
		panic(err)
	}

	for k, v := range mimeBuffer {
		for _, extension := range v.Extensions {
			extensions[extension] = k
		}
	}

}

func GetMimeType(name string) string {
	ext := strings.Replace(path.Ext(name), ".", "", 1)
	return extensions[strings.ToLower(ext)]
}

func GetExtension(mimetype string) string {
	match := regexExtract.FindAllStringSubmatch(mimetype, -1)
	m := mimeBuffer[strings.TrimSpace(strings.ToLower(match[0][1]))]
	if len(m.Extensions) == 0 {
		return ""
	}
	return m.Extensions[0]
}
