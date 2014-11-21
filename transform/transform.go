package transform

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type RequestJsonWrapper struct {
	Content string
}

type Transformer interface {
	Handler(writer http.ResponseWriter, request *http.Request)
	ApplyHeaders(header *http.Header)
}

type Base64JsonTransformer struct {
	Content   bytes.Buffer
	TargetUrl string
	AuthToken string
}

func (transformer *Base64JsonTransformer) Handler(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	log.Println("Received request on", request.RequestURI, "from", request.RemoteAddr)

	body, err := ioutil.ReadAll(request.Body)

	if err != nil {
		log.Println("Was unable to read the request body", err)
		return
	}

	writer.WriteHeader(http.StatusOK)

	var buf bytes.Buffer

	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	encoder.Write(body) // Write the original body of the request
	encoder.Close()

	wrapper := RequestJsonWrapper{Content: buf.String()}

	jsonEncoder := json.NewEncoder(&transformer.Content)
	jsonEncoder.Encode(wrapper)

	log.Println("Transformed request")
}

func (transformer *Base64JsonTransformer) ApplyHeaders(header *http.Header) {
	header.Set("Authorization", transformer.AuthToken)
}
