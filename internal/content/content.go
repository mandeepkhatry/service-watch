package content

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"os"
	"service-watch/internal/schema"
	"service-watch/internal/utils"

	"github.com/getkin/kin-openapi/openapi3"
)

var ContentBasedData = map[string]func(configSchema *openapi3.SchemaRef, encoding map[string]*openapi3.Encoding, components openapi3.Components, endpoint string) (*bytes.Buffer, string, error){
	"application/json": func(configSchema *openapi3.SchemaRef, encoding map[string]*openapi3.Encoding, components openapi3.Components, endpoint string) (*bytes.Buffer, string, error) {
		dummyData := schema.GenerateSchemaData(configSchema, components)
		requestBytes, err := json.Marshal(dummyData)
		if err != nil {
			return bytes.NewBuffer([]byte{}), "", err
		}

		return bytes.NewBuffer(requestBytes), "application/json", nil
	},
	"multipart/form-data": func(configSchema *openapi3.SchemaRef, encoding map[string]*openapi3.Encoding, components openapi3.Components, endpoint string) (*bytes.Buffer, string, error) {
		dummyData := schema.GenerateSchemaData(configSchema, components)

		fileContent := utils.FindFileContent(configSchema, components, encoding)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		root, _ := os.Getwd()

		for fileField, fileType := range fileContent {

			path := root + "/static/" + endpoint + fileType

			file, err := os.Open(path)
			if err != nil {
				return nil, "", err
			}
			fileContents, err := ioutil.ReadAll(file)
			if err != nil {
				return nil, "", err
			}

			fi, err := file.Stat()
			if err != nil {
				return nil, "", err
			}

			file.Close()

			part, err := writer.CreateFormFile(fileField, fi.Name())
			if err != nil {
				return nil, "", err
			}
			part.Write(fileContents)

		}

		for k, v := range dummyData {
			_ = writer.WriteField(k, v.(string))
		}

		contentType := writer.FormDataContentType()

		err := writer.Close()
		if err != nil {
			return nil, "", err
		}

		return body, contentType, nil
	},
}
