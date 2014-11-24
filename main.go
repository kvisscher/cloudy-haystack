package main

import (
	"flag"
	"fmt"
	"github.com/kvisscher/cloudy-haystack/config"
	"github.com/kvisscher/cloudy-haystack/transform"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	var bindPort int

	flag.IntVar(&bindPort, "port", 8080, "The port to bind on to receive requests")
	flag.Parse()

	logFile, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatal("Was unable to open log file", err)
	}

	defer logFile.Close()

	log.SetOutput(io.MultiWriter(os.Stdout, logFile))

	configFile, err := os.Open("config.json")

	if err != nil {
		log.Fatal("Was unable to open config.json", err)
	}

	mappingConfig := config.Parse(configFile)

	defer configFile.Close()

	log.Println("Binding to port", fmt.Sprintf(":%d", bindPort))

	ApplyMappings(&mappingConfig)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", bindPort), nil))
}

// Applies all of the mappings that have been configured
func ApplyMappings(mappingConfig *config.MappingConfig) {
	for _, mapping := range mappingConfig.Mappings {
		transformer := transform.Base64JsonTransformer{
			TargetUrl: fmt.Sprintf("%s%s", mappingConfig.TargetBaseUrl, mapping.To),
			AuthToken: mappingConfig.AuthToken,
		}

		http.HandleFunc(mapping.From, func(writer http.ResponseWriter, request *http.Request) {
			// Let the transformer do its magic
			transformer.Handler(writer, request)

			// Create a new request with the transformed content
			transformedRequest, err := http.NewRequest("POST", transformer.TargetUrl, &transformer.Content)

			if err != nil {
				log.Println("Something went wrong while making a request to", transformer.TargetUrl, err)
				return
			}

			// Let the request transformer apply additional headers
			transformer.ApplyHeaders(&transformedRequest.Header)

			// Send the request to the target
			client := http.Client{}
			response, err := client.Do(transformedRequest)

			if err != nil {
				log.Println("Error while posting to", transformer.TargetUrl, err)
				return
			}

			// Check if something went wrong on the server that should process the request
			if response.StatusCode != http.StatusOK {
				log.Println("Expected status code 200 while posting to", transformer.TargetUrl, "got", response.StatusCode)
			}

			// Discard all of the bytes of the response
			io.Copy(ioutil.Discard, response.Body)
		})

		log.Printf("Mapped %s -> %s%s\n", mapping.From, mappingConfig.TargetBaseUrl, mapping.To)
	}
}
