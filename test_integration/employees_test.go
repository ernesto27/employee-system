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

func TestGetEmployees(t *testing.T) {
	InitEnvTest()
	dbInstance := GetDBInstanceTest()

	tests := []Test{
		{
			name:           "Get Employees",
			statusExpected: http.StatusOK,
			countResponse:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/api/v1/employees", nil)
			if err != nil {
				t.Fatal(err)
			}

			newRecorder := httptest.NewRecorder()

			r := router.GetRouter(dbInstance.Db)
			r.ServeHTTP(newRecorder, req)

			if newRecorder.Result().StatusCode != tt.statusExpected {
				t.Fatalf("unexpected result: got %v want %v", newRecorder.Result().StatusCode, tt.statusExpected)
			}

			var resp response.Response
			err = json.NewDecoder(strings.NewReader(newRecorder.Body.String())).Decode(&resp)
			if err != nil {
				t.Fatal(err)
			}

			if len(resp.Data.([]interface{})) != tt.countResponse {
				t.Fatalf("unexpected result: got %v want %v", len(resp.Data.([]interface{})), tt.countResponse)
			}
		})
	}

}
