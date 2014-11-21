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

  // Fake a request that comes from somebody
	request := &http.Request{
		Method:        "POST",
		URL:           serverUrl,
		Body:          ioutil.NopCloser(strings.NewReader("<root><element>123</element></root>")),
		ContentLength: 35,
	}

	client := &http.Client{}

  // Execute the request that we made earlier
	response, err := client.Do(request)

	if err != nil {
		t.Errorf("unexpected error %s", err)
	}

  if response.StatusCode != http.StatusOK {
    t.Errorf("Expected status code 200 got %d", response.StatusCode)
  }

  // Expected result is a JSON object with the content
  // encoded of the original request in base64
	want := transform.RequestJsonWrapper{Content: "PHJvb3Q+PGVsZW1lbnQ+MTIzPC9lbGVtZW50Pjwvcm9vdD4="}

	var got transform.RequestJsonWrapper

	json.Unmarshal(transformer.Content.Bytes(), &got)

	if !reflect.DeepEqual(want, got) {
		t.Errorf("expected %+v got %+v", want, got)
	}
}
