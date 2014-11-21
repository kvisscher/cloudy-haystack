package transform_test

import (
	"encoding/json"
	"github.com/kvisscher/cloudy-haystack/transform"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestHandler(t *testing.T) {
	transformer := transform.Base64JsonTransformer{}

	server := httptest.NewServer(http.HandlerFunc(transformer.Handler))
	defer server.Close()

	serverUrl, _ := url.Parse(server.URL)

	request := &http.Request{
		Method:        "POST",
		URL:           serverUrl,
		Body:          ioutil.NopCloser(strings.NewReader("<root><element>123</element></root>")),
		ContentLength: 35,
	}

	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		t.Errorf("unexpected error %s", err)
	}

	expectedResponse := transform.RequestJsonWrapper{Content: "PHJvb3Q+PGVsZW1lbnQ+MTIzPC9lbGVtZW50Pjwvcm9vdD4="}

	var got transform.RequestJsonWrapper

	decoder := json.NewDecoder(response.Body)
	decoder.Decode(&got)

	if !reflect.DeepEqual(expectedResponse, got) {
		t.Errorf("expected %+v got %+v", expectedResponse, got)
	}
}
