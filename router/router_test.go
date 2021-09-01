package router

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func getHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "getHome")
}

func postProfilePictureUpload(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "postProfilePictureUpload")
}

func getProfile(w http.ResponseWriter, r *http.Request) {
	userId := GetParam(r, 0)
	fmt.Fprintf(w, "userid=%s", userId)
}

func TestRouter(t *testing.T) {
	routes := []Route{
		NewRoute("test1", http.MethodGet, "/", getHome),
		NewRoute("test2", http.MethodPost, "/profilePicture", postProfilePictureUpload),
		NewRoute("test3", http.MethodGet, "/profile/([0-9]+)", getProfile),
	}

	tests := map[string]struct {
		method           string
		url              string
		expectedResponse []byte
	}{
		"get request without any parameters": {
			method:           http.MethodGet,
			url:              "/",
			expectedResponse: []byte("getHome"),
		},
		"post request without any parameters": {
			method:           http.MethodPost,
			url:              "/profilePicture",
			expectedResponse: []byte("postProfilePictureUpload"),
		},
		"get request with a numeric parameter": {
			method:           http.MethodGet,
			url:              "/profile/12345",
			expectedResponse: []byte("userid=12345"),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			u, err := url.Parse(test.url)
			if err != nil {
				t.Fatalf("Error parsing URL: %s", err)
			}
			r := &http.Request{
				Method: test.method,
				URL:    u,
			}
			w := httptest.NewRecorder()

			Serve(routes, w, r)
			w.Flush()
			responseBytes := w.Body.Bytes()
			if !bytes.Equal(responseBytes, test.expectedResponse) {
				t.Fatalf("Did not get the expected response (expected: \"%s\", received: \"%s\")", test.expectedResponse, responseBytes)
			}
		})
	}
}
