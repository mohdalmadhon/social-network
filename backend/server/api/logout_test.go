package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogOutUserClearsTokenCookie(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/api/logout", nil)
	response := httptest.NewRecorder()

	App{}.LogOutUser(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}

	cookies := response.Result().Cookies()
	if len(cookies) != 1 || cookies[0].Name != "token" || cookies[0].Value != "" || cookies[0].MaxAge != -1 {
		t.Fatal("logout did not clear the token cookie")
	}
}
