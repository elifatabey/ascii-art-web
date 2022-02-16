package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

)

func TestAsciiArt(t *testing.T) {
	testCases := []struct {
		banner           string
		input            string
		expectedResponse int
	}{
		{
			// when valid banner and input is passed
			banner:           "shadow",
			input:            "this is history",
			expectedResponse: http.StatusOK,
		},
		{
			// when an invalid input file is passed
			banner:           "shadow",
			input:            "",
			expectedResponse: http.StatusBadRequest,
		},
		{
			// when an invalid banner file
			banner:           "shadwo",
			input:            "this is history",
			expectedResponse: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		form := url.Values{}
		form.Add("bannertype", tc.banner)
		form.Add("inputtext", tc.input)

		request := httptest.NewRequest("POST", "/ascii-art", strings.NewReader(form.Encode()))
		request.PostForm = form
		responseRecorder := httptest.NewRecorder()

		Display(responseRecorder, request)
		if responseRecorder.Code != tc.expectedResponse {
			t.Errorf("Want status '%d', got '%d'", tc.expectedResponse, responseRecorder.Code)
		}
		// assert.Equal(t, responseRecorder.Code, tc.expectedResponse)
	}
}

func TestPostEndpoint(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(Display))
	defer svr.Close()
	form := url.Values{}
	form.Add("bannertype", "thinkertoy")
	form.Add("inputtext", "hello sad world")

	request, err := http.NewRequest("POST", svr.URL+"/ascii-art", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Want status '%d', got '%d'", 200, resp.StatusCode)
	}
}

func TestGetEndpoint(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(Home))
	defer svr.Close()
	resp, err := http.Get(svr.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("Want status '%d', got '%d'", 200, resp.StatusCode)
	}
}
