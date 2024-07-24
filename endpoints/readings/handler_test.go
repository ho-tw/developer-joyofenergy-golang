package readings

import (
	"encoding/json"
	"io"
	"joi-energy-golang/domain"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"

	"github.com/stretchr/testify/assert"
)

func callEndpoint(handler http.HandlerFunc, url string, body io.Reader, t *testing.T) *httptest.ResponseRecorder {
	return callEndpointMethod(http.MethodPost, handler, url, body, t)
}

func callEndpointMethod(method string, handler http.HandlerFunc, url string, body io.Reader, t *testing.T) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatalf("request creation failed: %s", err)
	}
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}

func TestStoreReadingsReturnResultFromService(t *testing.T) {
	s := &MockService{}
	h := NewHandler(s)
	params := httprouter.Params{}
	storeReadings := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		h.StoreReadings(writer, request, params)
	})

	body := `{"smartMeterId": "smartMeterId", "electricityReadings": []}`
	rr := callEndpoint(storeReadings, "/readings/store", strings.NewReader(body), t)
	assert.Equal(t, http.StatusOK, rr.Code, "handler returned status code %v on valid request", rr.Code)
}

func TestStoreReadingsWithInvalidInput(t *testing.T) {
	s := &MockService{}
	h := NewHandler(s)
	params := httprouter.Params{}
	storeReadings := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		h.StoreReadings(writer, request, params)
	})

	body := ""
	rr := callEndpoint(storeReadings, "/readings/store", strings.NewReader(body), t)
	assert.Equal(t, http.StatusBadRequest, rr.Code, "handler returned status code %v on invalid request", rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	expectedMessage := domain.ErrorResponse{
		Message: "unmarshal request body failed: unexpected end of JSON input",
	}
	var actualMessage domain.ErrorResponse

	err := json.Unmarshal(rr.Body.Bytes(), &actualMessage)
	assert.NoError(t, err)
	assert.Equal(t, expectedMessage, actualMessage)
}

func TestGetLastWeekReadings(t *testing.T) {
	s := &MockService{}
	h := NewHandler(s)
	params := httprouter.Params{{Key: "smartMeterId", Value: "123"}}
	getLastWeekReadings := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		h.GetLastWeekUsage(writer, request, params)
	})

	rr := callEndpointMethod(http.MethodGet, getLastWeekReadings, "/readings/last-week/123", nil, t)
	assert.Equal(t, http.StatusOK, rr.Code, "handler returned status code %v on valid request", rr.Code)

	expectedResult := domain.LastWeekUsage{
		SmartMeterId: "123",
		Cost:         0,
	}
	var actualResult domain.LastWeekUsage

	err := json.Unmarshal(rr.Body.Bytes(), &actualResult)
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

type MockService struct {
	Service
}

func (s *MockService) StoreReadings(smartMeterId string, reading []domain.ElectricityReading) {}
func (s *MockService) CalculateLastWeekUsageCost(smartMeterId string, tariff float64) (float64, error) {
	return 0, nil
}
