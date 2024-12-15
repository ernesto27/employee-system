package test_integration

import (
	"employees-system/response"
	"employees-system/router"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestVersionAPI(t *testing.T) {
	InitEnvTest()
	dbInstance := GetDBInstanceTest()

	req, err := http.NewRequest("GET", "/api/v1/", nil)
	if err != nil {
		t.Fatal(err)
	}

	newRecorder := httptest.NewRecorder()

	r := router.GetRouter(dbInstance.Db)
	r.ServeHTTP(newRecorder, req)

	if newRecorder.Result().StatusCode != http.StatusOK {
		t.Fatalf("unexpected result: got %v want %v", newRecorder.Result().StatusCode, http.StatusOK)
	}

	var resp response.Response
	err = json.NewDecoder(strings.NewReader(newRecorder.Body.String())).Decode(&resp)
	if err != nil {
		t.Fatal(err)
	}

	want := "employees system api v1"
	if resp.Message != want {
		t.Fatalf("unexpected result: got %v want %v", resp.Message, want)
	}

}
