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

func UpdateHandler(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	log.Println("Received request on", request.RequestURI, "from", request.RemoteAddr)

	body, err := ioutil.ReadAll(request.Body)

	if err != nil {
		log.Println("Was unable to read the request body")
		return
	}

	writer.WriteHeader(http.StatusOK)

	var buf bytes.Buffer

	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	encoder.Write(body) // Write the original body of the request
	encoder.Close()

	wrapper := RequestJsonWrapper{Content: buf.String()}

	jsonEncoder := json.NewEncoder(writer)
	jsonEncoder.Encode(wrapper)
	encoder.Close()

	log.Println("Send response")
}
