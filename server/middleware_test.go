package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func MockMiddleware(tokenValid bool, path string) AuthorizationMiddleware {
	var code int
	if tokenValid {
		switch path {
		case "/chat":
			code = 101
		case "/sign_out":
			code = 200
		}
	} else if path == "/sign_up" || path == "/sign_in" {
		code = 200
	} else {
		code = 303
	}
	return AuthorizationMiddleware{
		handler: MockedHandler{
			status: code,
		},
		verifier: MockVerifier(tokenValid),
	}
}

type MockedHandler struct {
	status int
}

func (mh MockedHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(mh.status)
}

func MockVerifier(tokenValid bool) func(c *http.Cookie) bool {
	if tokenValid {
		return func(c *http.Cookie) bool {
			return true
		}
	} else {
		return func(c *http.Cookie) bool {
			return false
		}
	}
}

func stubRequest(cookie bool, url string, method string) *http.Request {
	request := httptest.NewRequest(method, url, nil)
	if cookie {
		cookie := &http.Cookie{
			Name:  "token",
			Value: "true",
		}
		request.AddCookie(cookie)
	}
	return request
}

func TestMiddleware(t *testing.T) {
	url := "http://localhost:8090"
	t.Run("When target path is restricted and token is not valid returns 303 ", func(t *testing.T) {
		middleware := MockMiddleware(false, "/chat")
		request := stubRequest(false, url+"/chat", "GET")
		recorder := httptest.NewRecorder()
		middleware.ServeHTTP(recorder, request)
		response := recorder.Result()
		if response.StatusCode != 401 {
			t.Errorf("Response code should be 401, instead got %d", response.StatusCode)
		}
		if redirectURL, _ := response.Location(); redirectURL.Path != "/sign_up" {
			t.Errorf("Response location should be /sign_up, instead got %s", redirectURL.Path)
		}
	})
	t.Run("When target path is /sign_up or /sign_in and token is valid returns 303", func(t *testing.T) {
		middleware := MockMiddleware(true, "/sign_up")
		request := stubRequest(true, url+"/sign_up", "GET")
		recorder := httptest.NewRecorder()
		middleware.ServeHTTP(recorder, request)
		response := recorder.Result()
		if response.StatusCode != 303 {
			t.Errorf("Response code should be 303, instead got %d", response.StatusCode)
		}
	})

	t.Run("When target path is /chat and token is valid returns 101", func(t *testing.T) {
		middleware := MockMiddleware(true, "/chat")
		request := stubRequest(true, url+"/chat", "GET")
		recorder := httptest.NewRecorder()
		middleware.ServeHTTP(recorder, request)
		response := recorder.Result()
		if response.StatusCode != 101 {
			t.Errorf("Response code should be 101, instead got %d", response.StatusCode)
		}
	})

	t.Run("When target path is /sign_out and token is valid returns 200", func(t *testing.T) {
		middleware := MockMiddleware(true, "/sign_out")
		request := stubRequest(true, url+"/sign_out", "GET")
		recorder := httptest.NewRecorder()
		middleware.ServeHTTP(recorder, request)
		response := recorder.Result()
		if response.StatusCode != 200 {
			t.Errorf("Response code should be 200, instead got %d", response.StatusCode)
		}
	})

}
