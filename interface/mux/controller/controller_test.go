package controller_test

import (
	"assessment/interface/mux/controller"
	"assessment/interface/mux/router"
	repoMock "assessment/repository/mock"
	"assessment/service"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

var rt *mux.Router

func TestController_FetchAllPhoneNumbers(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mockRepository := repoMock.NewMockPhoneNumberRepository(mockController)

	mockRepository.EXPECT().FetchPaginatedPhoneNumbers(0, 11).
		Return([]string{
			"(237) 697151594",
			"(212) 654642448",
			"(258) 042423566",
			"(256) 7734127498",
		}, nil).AnyTimes()
	validator := service.NewValidator()

	svc := service.NewNumberService(validator, mockRepository)

	ctrl := controller.NewNumberController(svc)

	rt = router.InitRouter(ctrl)

	req := httptest.NewRequest(http.MethodGet, "/numbersvc?limit=10&page=1", nil)

	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req = httptest.NewRequest(http.MethodGet, "/numbersvc", nil)

	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req = httptest.NewRequest(http.MethodGet, "/numbersvc?limit=-1&page=1", nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestController_FilterPhoneNumbersByCountryAndState(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mockRepository := repoMock.NewMockPhoneNumberRepository(mockController)

	mockRepository.EXPECT().FetchPaginatedPhoneNumbersByCode("237", 0, 11).
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
		}, nil).AnyTimes()

	validator := service.NewValidator()

	svc := service.NewNumberService(validator, mockRepository)

	ctrl := controller.NewNumberController(svc)

	rt = router.InitRouter(ctrl)

	req := httptest.NewRequest(http.MethodGet, "/numbersvc?limit=10&page=1&country=cameroon", nil)

	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req = httptest.NewRequest(http.MethodGet, "/numbersvc?limit=10&page=1&country=nigeria&state=OK", nil)

	response = executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	req = httptest.NewRequest(http.MethodGet, "/numbersvc?limit=10&page=1&country=nigeria&state=INVALID", nil)

	response = executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)

	req = httptest.NewRequest(http.MethodGet, "/numbersvc?limit=10&page=1&state=NOK", nil)

	response = executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()

	rt.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	require.Equal(t, expected, actual)
}
