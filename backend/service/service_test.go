package service

import (
	repoMock "assessment/repository/mock"
	serviceMock "assessment/service/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type testSuite struct {
	suite.Suite
	svc *NumberService
}

func (t *testSuite) SetupSuite() {
	mockRepo := new(repoMock.PhoneNumberRepository)
	mockValidator := new(serviceMock.NumberValidator)

	mockValidator.On("Validate", "(237) 697151594").Return("Cameroon", "+237", "697151594", true)
	mockValidator.On("Validate", "(237) 699209115").Return("Cameroon", "+237", "699209115", false)
	mockValidator.On("GetCodeFromCountry", "cameroon").Return("237", nil)

	// ============== Test Data For All Phone Numbers  ===================== \\
	mockRepo.On("FetchPaginatedPhoneNumbers", 0, 6).Return([]string{
		"(237) 697151594",
		"(237) 697151594",
		"(237) 697151594",
		"(237) 697151594",
		"(237) 697151594",
		"(237) 697151594",
	}, nil)

	mockRepo.On("FetchPaginatedPhoneNumbers", 5, 6).Return([]string{
		"(237) 697151594",
		"(237) 697151594",
		"(237) 697151594",
		"(237) 697151594",
		"(237) 697151594",
		"(237) 697151594",
	}, nil)

	mockRepo.On("FetchPaginatedPhoneNumbers", 10, 6).Return([]string{
		"(237) 697151594",
		"(237) 697151594",
		"(237) 697151594",
	}, nil)

	mockRepo.On("FetchPaginatedPhoneNumbers", 15, 6).Return([]string{}, nil)

	mockRepo.On("FetchPaginatedPhoneNumbers", 0, 3).Return([]string{
		"(237) 699209115",
		"(237) 699209115",
	}, nil)
	// ============================================================================== \\

	// =========================== Test Data For Filter By State And Filter By Country ==================== \\
	mockRepo.On("FetchPaginatedPhoneNumbers", 0, 11).Return([]string{
		"(237) 699209115",
		"(237) 699209115",
		"(237) 699209115",
		"(237) 699209115",
		"(237) 699209115",
	}, nil)
	mockRepo.On("FetchPaginatedPhoneNumbers", 10, 11).Return([]string{}, nil)
	t.svc = NewNumberService(mockValidator, mockRepo)

	// ============================================================================== \\

	// ============================ Test Data For Filter By Country And State ====================== \\
	mockRepo.On("FetchPaginatedPhoneNumbersByCode", "237", 0, 5).Return([]string{
		"(237) 697151594",
		"(237) 697151594",
		"(237) 697151594",
		"(237) 697151594",
		"(237) 697151594",
	}, nil)

	mockRepo.On("FetchPaginatedPhoneNumbersByCode", "237", 0, 4).Return([]string{
		"(237) 697151594",
		"(237) 699209115",
		"(237) 697151594",
		"(237) 699209115",
	}, nil)
	mockRepo.On("FetchPaginatedPhoneNumbersByCode", "237", 3, 4).Return([]string{
		"(237) 699209115",
		"(237) 697151594",
		"(237) 697151594",
		"(237) 699209115",
	}, nil)

}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(testSuite))
}

func (t *testSuite) Test_FetchPhoneNumbers() {
	result, err := t.svc.FetchPhoneNumbers("1", "5")

	require.NoError(t.T(), err, "Expected: nil\nGot: %v\n", err)

	for _, d := range result.Data {
		require.Equal(t.T(), "OK", d.State)
	}

	result, err = t.svc.FetchPhoneNumbers("1", "2")

	require.NoError(t.T(), err, "Expected: nil\nGot: %v\n", err)

	for _, d := range result.Data {
		require.Equal(t.T(), "NOK", d.State)
	}

	require.Equal(t.T(), 2, len(result.Data))
	require.Equal(t.T(), false, result.Meta.Next)
	require.Equal(t.T(), false, result.Meta.Prev)

	result, err = t.svc.FetchPhoneNumbers("2", "5")

	require.NoError(t.T(), err, "Expected: nil\nGot: %v\n", err)

	require.Equal(t.T(), 5, len(result.Data))
	require.Equal(t.T(), true, result.Meta.Next)
	require.Equal(t.T(), true, result.Meta.Prev)

	result, err = t.svc.FetchPhoneNumbers("3", "5")

	require.NoError(t.T(), err, "Expected: nil\nGot: %v\n", err)

	require.Equal(t.T(), 3, len(result.Data))
	require.Equal(t.T(), false, result.Meta.Next)
	require.Equal(t.T(), true, result.Meta.Prev)

	result, err = t.svc.FetchPhoneNumbers("", "")

	require.NoError(t.T(), err, "Expected: nil\nGot: %v\n", err)

	require.Equal(t.T(), 5, len(result.Data))
	require.Equal(t.T(), true, result.Meta.Next)
	require.Equal(t.T(), false, result.Meta.Prev)

	result, err = t.svc.FetchPhoneNumbers("4", "5")

	require.NoError(t.T(), err, "Expected: nil\nGot: %v\n", err)

	require.Equal(t.T(), 0, len(result.Data))
	require.Equal(t.T(), false, result.Meta.Next)
	require.Equal(t.T(), false, result.Meta.Prev)

	result, err = t.svc.FetchPhoneNumbers("-1", "4")

	require.Error(t.T(), err, "Expected An Error\nGot: %v\n", err)

	result, err = t.svc.FetchPhoneNumbers("1", "-4")

	require.Error(t.T(), err, "Expected An Error\nGot: %v\n", err)

}

func (t *testSuite) Test_FilterByState() {

	result, err := t.svc.FilterByState("OK", "1", "5")
	require.NoError(t.T(), err, "Expected: nil\nGot: %v\n", err)

	for _, d := range result.Data {
		require.Equal(t.T(), "OK", d.State)
	}

	result, err = t.svc.FilterByState("NOK", "1", "10")
	require.NoError(t.T(), err, "Expected: nil\nGot: %v\n", err)

	for _, d := range result.Data {
		require.Equal(t.T(), "NOK", d.State)
	}

	require.Equal(t.T(), false, result.Meta.Next)
	require.Equal(t.T(), false, result.Meta.Prev)

	_, err = t.svc.FilterByState("INVALID", "1", "10")
	require.Error(t.T(), err, "Expected An Error\nGot: %v\n", err)

	_, err = t.svc.FilterByState("VALID", "1", "10")
	require.Error(t.T(), err, "Expected An Error\nGot: %v\n", err)

	_, err = t.svc.FilterByState("INVALID", "-1", "10")
	require.Error(t.T(), err, "Expected An Error\nGot: %v\n", err)

	_, err = t.svc.FilterByState("INVALID", "1", "10a")
	require.Error(t.T(), err, "Expected An Error\nGot: %v\n", err)

}

func (t *testSuite) Test_FilterByCountry() {
	result, err := t.svc.FilterByCountry("cameroon", "1", "4")
	require.NoError(t.T(), err, "Expected: nil\nGot: %v\n", err)

	for _, d := range result.Data {
		require.Equal(t.T(), "Cameroon", d.Country)
		require.Equal(t.T(), "+237", d.CountryCode)
	}

	_, err = t.svc.FilterByCountry("cameroon", "-1", "4")
	require.Error(t.T(), err, "Expected An Error\nGot: %v\n", err)

	_, err = t.svc.FilterByCountry("cameroon", "1", "4+")
	require.Error(t.T(), err, "Expected An Error\nGot: %v\n", err)

}

func (t *testSuite) Test_FilterByCountryAndState() {
	result, err := t.svc.FilterByCountryAndState("cameroon", "OK", "1", "3")
	require.NoError(t.T(), err, "Expected: nil\nGot: %v\n", err)

	for _, d := range result.Data {
		require.Equal(t.T(), 3, len(result.Data))
		require.Equal(t.T(), "Cameroon", d.Country)
		require.Equal(t.T(), "+237", d.CountryCode)
		require.Equal(t.T(), "OK", d.State)
		require.Equal(t.T(), true, result.Meta.Next)
	}

	_, err = t.svc.FilterByCountryAndState("cameroon", "MOK", "1", "3")
	require.Error(t.T(), err, "Expected An Error\nGot: %v\n", err)

	_, err = t.svc.FilterByCountryAndState("cameroon", "NOK", "-1", "3")
	require.Error(t.T(), err, "Expected An Error\nGot: %v\n", err)

	_, err = t.svc.FilterByCountryAndState("cameroon", "NOK", "1", "jumia")
	require.Error(t.T(), err, "Expected An Error\nGot: %v\n", err)

}
