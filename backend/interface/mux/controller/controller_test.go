package controller_test

import (
	"assessment/interface/mux/controller"
	"assessment/interface/mux/router"
	repoMock "assessment/repository/mock"
	"assessment/service"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

var rt *mux.Router

type testSuite struct {
	suite.Suite
	ctrl *controller.Controller
}

func (t *testSuite) SetupSuite() {
	mockRepo := new(repoMock.PhoneNumberRepository)
	mockRepo.On("FetchPaginatedPhoneNumbers", 0, 11).
		Return([]string{
			"(237) 697151594",
			"(212) 654642448",
			"(258) 042423566",
			"(256) 7734127498",
		}, nil)
	mockRepo.On("FetchPaginatedPhoneNumbers", 0, 6).
		Return([]string{}, nil)
	mockRepo.On("FetchPaginatedPhoneNumbersByCode", "237", 0, 11).
		Return([]string{
			"(237) 23456789",
			"(237) 23456789",
			"(237) 23456789",
			"(237) 23456789",
			"(237) 23456789",
			"(237) 23456789",
			"(237) 23456789",
			"(237) 23456789",
			"(237) 23456789",
			"(237) 23456789",
			"(237) 23456789",
		}, nil)

	mockRepo.On("FetchPaginatedPhoneNumbers", mock.Anything, 11).
		Return([]string{
			"(237) 697151594",
			"(212) 654642448",
			"(258) 042423566",
			"(256) 7734127498",
		}, nil)

	validator := service.NewValidator()

	svc := service.NewNumberService(validator, mockRepo)

	t.ctrl = controller.NewNumberController(svc)

	rt = router.InitRouter(t.ctrl)
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(testSuite))
}
func (t *testSuite) TestController_FetchAllPhoneNumbers() {

	req := httptest.NewRequest(http.MethodGet, "/phone-numbers?limit=10&page=1", nil)

	response := executeRequest(req)

	checkResponseCode(t.T(), http.StatusOK, response.Code)

	req = httptest.NewRequest(http.MethodGet, "/phone-numbers", nil)

	response = executeRequest(req)

	checkResponseCode(t.T(), http.StatusOK, response.Code)

	req = httptest.NewRequest(http.MethodGet, "/phone-numbers?limit=-1&page=1", nil)
	response = executeRequest(req)

	checkResponseCode(t.T(), http.StatusBadRequest, response.Code)
}

func (t *testSuite) TestController_FetchPhoneNumberByCountry() {
	req := httptest.NewRequest(http.MethodGet, "/phone-numbers?limit=10&page=1&country=cameroon", nil)

	response := executeRequest(req)

	checkResponseCode(t.T(), http.StatusOK, response.Code)
}

func (t *testSuite) TestController_FetchPhoneNumbersByState() {
	req := httptest.NewRequest(http.MethodGet, "/phone-numbers?limit=10&page=1&state=NOK", nil)

	response := executeRequest(req)

	checkResponseCode(t.T(), http.StatusOK, response.Code)

	req = httptest.NewRequest(http.MethodGet, "/phone-numbers?limit=10&page=1&state=OK", nil)

	response = executeRequest(req)

	checkResponseCode(t.T(), http.StatusOK, response.Code)

	req = httptest.NewRequest(http.MethodGet, "/phone-numbers?limit=10&page=1&state=VALID", nil)

	response = executeRequest(req)

	checkResponseCode(t.T(), http.StatusBadRequest, response.Code)

}

func (t *testSuite) TestController_FetchPhoneNumbersByCountryAndState() {
	req := httptest.NewRequest(http.MethodGet, "/phone-numbers?limit=10&page=1&country=nigeria&state=OK", nil)

	response := executeRequest(req)

	checkResponseCode(t.T(), http.StatusNotFound, response.Code)

	req = httptest.NewRequest(http.MethodGet, "/phone-numbers?limit=10&page=1&country=nigeria&state=NOK", nil)

	response = executeRequest(req)

	checkResponseCode(t.T(), http.StatusNotFound, response.Code)

	req = httptest.NewRequest(http.MethodGet, "/phone-numbers?limit=10&page=1&country=nigeria&state=INVALID", nil)

	response = executeRequest(req)

	checkResponseCode(t.T(), http.StatusNotFound, response.Code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()

	rt.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	require.Equal(t, expected, actual)
}
