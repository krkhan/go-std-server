package main

import (
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"testing"
	"time"
)

const (
	contentType = "application/x-www-form-urlencoded"
)

func TestHTTPServer(t *testing.T) {
	serverAddr := fmt.Sprintf("127.0.0.1:%d", (rand.Uint32()%1000)+40000)
	httpServerExitDone := &sync.WaitGroup{}
	httpServerExitDone.Add(1)
	server, _ := startHttpServer(httpServerExitDone, serverAddr)

	tests := map[string]struct {
		method             string
		url                string
		body               []byte
		expectedStatusCode int
		expectedResponse   []byte
	}{
		"valid request": {
			method:             http.MethodPost,
			url:                "hash",
			body:               []byte(fmt.Sprintf("%s=angryMonkey", PasswordKey)),
			expectedStatusCode: http.StatusOK,
			expectedResponse:   []byte("1"),
		},
		"multiple passwords": {
			method:             http.MethodPost,
			url:                "hash",
			body:               []byte(fmt.Sprintf("%s=angryMonkey&%s=angrierPrimate", PasswordKey, PasswordKey)),
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   nil,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			parsedUrl, err := url.Parse(fmt.Sprintf("http://%s/%s", serverAddr, test.url))
			if err != nil {
				t.Fatalf("Error parsing URL: %s", err)
			}
			var resp *http.Response
			switch test.method {
			case http.MethodGet:
				t.Logf("Making GET call to %s", parsedUrl)
				resp, err = http.Get(parsedUrl.String())
				if err != nil {
					t.Fatalf("Error making GET request: %s", err)
				}
			case http.MethodPost:
				t.Logf("Making POST call to %s with %d bytes in body", parsedUrl, len(test.body))
				resp, err = http.Post(parsedUrl.String(), contentType, bytes.NewBuffer(test.body))
				if err != nil {
					t.Fatalf("Error making POST request: %s", err)
				}
			default:
				t.Fatalf("Cannot test with method: %s", test.method)
			}
			if resp.StatusCode != test.expectedStatusCode {
				t.Fatalf("Did not get the expected status code (expected: %d, received: %d)", test.expectedStatusCode, resp.StatusCode)
			}
			responseBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Error reading response body: %s", err)
			}
			if test.expectedResponse != nil && !bytes.Equal(responseBytes, test.expectedResponse) {
				t.Fatalf("Did not get the expected response (expected: \"%s\", received: \"%s\")", test.expectedResponse, responseBytes)
			}
		})
	}

	if err := server.Shutdown(context.TODO()); err != nil {
		panic(err)
	}
}

func TestDelay(t *testing.T) {
	serverAddr := fmt.Sprintf("127.0.0.1:%d", (rand.Uint32()%1000)+40000)
	httpServerExitDone := &sync.WaitGroup{}
	httpServerExitDone.Add(1)
	server, _ := startHttpServer(httpServerExitDone, serverAddr)

	postHashUrl, err := url.Parse(fmt.Sprintf("http://%s/hash", serverAddr))
	if err != nil {
		t.Fatalf("Error parsing URL: %s", err)
	}

	testContent := "blahblahblahblah"
	testDigest := sha512.Sum512([]byte(testContent))
	testDigestB64 := []byte(base64.StdEncoding.EncodeToString(testDigest[:]))

	resp, err := http.Post(postHashUrl.String(), contentType, bytes.NewBuffer([]byte(fmt.Sprintf("password=%s", testContent))))
	if err != nil {
		t.Fatalf("Error sending hash request: %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Got invalid status code: %d", resp.StatusCode)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %s", err)
	}

	digestKey, err := strconv.ParseInt(string(respBody), 10, 64)
	if err != nil {
		t.Fatalf("Error parsing digest key: %s", err)
	}

	getHashUrl, err := url.Parse(fmt.Sprintf("http://%s/hash/%d", serverAddr, digestKey))
	if err != nil {
		t.Fatalf("Error parsing URL: %s", err)
	}

	// The test obviously assumes that Sha512DelaySeconds is at least longer than the time
	// it took between the two requests. In other words, we should NOT be getting the digest
	// right now, until we have waits the required number of seconds
	//
	// If, e.g., Sha512DelaySeconds is 1 and this test is slower than a second
	// for some reason it might generate a false negative

	resp, err = http.Get(getHashUrl.String())
	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %s", err)
	}
	expectedRespBody := []byte(fmt.Sprintf("Key not found: %d", digestKey))

	if !bytes.Equal(respBody, expectedRespBody) {
		t.Fatalf(`Invalid response. Got: "%s", expected: "%s"`, respBody, expectedRespBody)
	}

	time.Sleep(Sha512DelaySeconds*time.Second + 1)

	resp, err = http.Get(getHashUrl.String())
	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %s", err)
	}

	if !bytes.Equal(respBody, testDigestB64) {
		t.Fatalf(`Invalid response. Got: "%s", expected: "%s"`, respBody, testDigestB64)
	}

	if err := server.Shutdown(context.TODO()); err != nil {
		panic(err)
	}
}
