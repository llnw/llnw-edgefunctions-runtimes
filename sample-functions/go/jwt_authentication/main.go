package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/llnw/llnw-edgefunctions-runtimes/go/edgefunction"
	"github.com/llnw/llnw-edgefunctions-runtimes/go/edgefunction/events"
)

const JWTENV = "JWT_SECRET"

// JWT Secret
var secret []byte

func handler(request *events.EPInvokeRequest) (*events.EPInvokeResponse, error) {
	// Load our JWT authorization secret
	if err := loadSecret(); err != nil {
		return nil, err
	}

	resp := &events.EPInvokeResponse{
		StatusCode: http.StatusUnauthorized,
	}
	tokenString, err := getJWTToken(request)
	if err != nil {
		fmt.Println(err.Error())
		return resp, nil
	}

	valid, err := isValid(*tokenString)
	if err != nil {
		fmt.Println(err.Error())
		return resp, nil
	}

	if !valid {
		return resp, nil
	}

	// Proxy request
	return proxy(request)
}

func loadSecret() error {
	if len(secret) == 0 {
		value, ok := os.LookupEnv(JWTENV)
		if !ok {
			return fmt.Errorf("failed to lookup %s ENV", JWTENV)
		}

		secret = []byte(value)
	}

	return nil
}

func getJWTToken(request *events.EPInvokeRequest) (*string, error) {
	value, ok := request.Headers["Authorization"]
	if !ok {
		return nil, errors.New("no Authorization header found")
	}

	if !strings.Contains(strings.ToUpper(value), "BEARER") {
		return nil, errors.New("Authorization header missing BEARER prefix")
	}

	tokenString := value[7:]
	return &tokenString, nil
}

func isValid(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return secret, nil
	})

	if err != nil {
		return false, err
	}

	return token.Valid, nil
}

func proxy(epRequest *events.EPInvokeRequest) (*events.EPInvokeResponse, error) {
	// Parse the ep request into an http request
	req, err := parseRequest(epRequest)
	if err != nil {
		return nil, err
	}

	// Do http call
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse http response into ep response
	return parseResponse(resp)
}

func parseRequest(epRequest *events.EPInvokeRequest) (*http.Request, error) {
	var body io.Reader
	if epRequest.Body != "" {
		body = bytes.NewBuffer([]byte(epRequest.Body))
	} else {
		body = http.NoBody
	}

	url := fmt.Sprintf("http://%s%s", epRequest.Host, epRequest.Path)
	// Format request
	req, err := http.NewRequest(epRequest.Method, url, body)
	if err != nil {
		return nil, err
	}

	// Copy Headers
	for k, v := range epRequest.Headers {
		req.Header.Add(k, v)
	}

	for k, values := range epRequest.MultiValueHeaders {
		for _, v := range values {
			req.Header.Add(k, v)
		}
	}

	// Add Queries
	params := req.URL.Query()
	for k, v := range epRequest.Queries {
		params.Add(k, v)
	}

	for k, values := range epRequest.MultiValueQueries {
		for _, v := range values {
			params.Add(k, v)
		}
	}

	req.URL.RawQuery = params.Encode()

	return req, nil
}

func parseResponse(resp *http.Response) (*events.EPInvokeResponse, error) {
	epResp := &events.EPInvokeResponse{
		StatusCode:        resp.StatusCode,
		Headers:           make(map[string]string),
		MultiValueHeaders: make(map[string][]string),
	}

	// Copy Body
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	epResp.Body = string(data)

	// Copy Headers
	copyRespHeaders(resp.Header, epResp.Headers, epResp.MultiValueHeaders)

	return epResp, nil
}

func copyRespHeaders(inMap map[string][]string, singleValueLookup map[string]string, multiValueLookup map[string][]string) {
	for key, values := range inMap {
		for _, value := range values {
			if len(values) == 1 {
				singleValueLookup[key] = value
			} else {
				copyIntoMultiValues(value, multiValueLookup, key)
			}
		}
	}
}

func copyIntoMultiValues(value string, multiValueLookup map[string][]string, key string) {
	multivalues, ok := multiValueLookup[key]
	if !ok {
		multivalues = make([]string, 0)
	}

	multivalues = append(multivalues, value)
	multiValueLookup[key] = multivalues
}

func main() {
	edgefunction.Start(handler)
}
