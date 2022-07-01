package service

import (
	"assessment/apperror"
	repoMock "assessment/repository/mock"
	serviceMock "assessment/service/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNumberService_FetchPhoneNumbers(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mockRepository := repoMock.NewMockPhoneNumberRepository(mockController)
	mockValidator := serviceMock.NewMockNumberValidator(mockController)

	mockRepository.EXPECT().FetchPaginatedPhoneNumbers(0, 5).
		Return([]string{
			"(237) 697151594",
			"(212) 654642448",
			"(258) 042423566",
			"(256) 7734127498",
			"(237) 696443597",
		}, nil)

	mockValidator.EXPECT().Validate(gomock.Any()).Return("COUNTRY", "CODE", "NUMBER", true).AnyTimes()

	numSvc := NewNumberService(mockValidator, mockRepository)

	result, err := numSvc.FetchPhoneNumbers("1", "4")

	require.NoError(t, err, "Expected No Errors\nGot: %v\n", err)

	for _, d := range result.Data {
		require.Equal(t, "OK", d.State)
	}

	require.Equal(t, 4, len(result.Data))
	require.Equal(t, true, result.Meta.Next)
	require.Equal(t, false, result.Meta.Prev)

	result, err = numSvc.FetchPhoneNumbers("-1", "4")

	require.Error(t, err, "Expected Error\nGot No Error")

	result, err = numSvc.FetchPhoneNumbers("1", "-4")

	require.Error(t, err, "Expected Error\nGot No Error")

	result, err = numSvc.FetchPhoneNumbers("1a", "4")

	require.Error(t, err, "Expected Error\nGot No Error")

	mockRepository.EXPECT().FetchPaginatedPhoneNumbers(0, 5).
		Return([]string{
			"(237) 697151594",
			"(212) 654642448",
			"(258) 042423566",
			"(256) 7734127498",
		}, nil)

	result, err = numSvc.FetchPhoneNumbers("1", "4")

	require.Equal(t, false, result.Meta.Next)
	require.Equal(t, false, result.Meta.Prev)

	mockRepository.EXPECT().FetchPaginatedPhoneNumbers(4, 5).
		Return([]string{
			"(237) 697151594",
			"(212) 654642448",
			"(258) 042423566",
			"(256) 7734127498",
		}, nil)

	result, err = numSvc.FetchPhoneNumbers("2", "4")

	require.Equal(t, false, result.Meta.Next)
	require.Equal(t, true, result.Meta.Prev)

	mockRepository.EXPECT().FetchPaginatedPhoneNumbers(4, 5).
		Return([]string{
			"(237) 697151594",
			"(212) 654642448",
			"(258) 042423566",
			"(256) 7734127498",
			"(258) 042423566",
		}, nil)

	result, err = numSvc.FetchPhoneNumbers("2", "4")

	require.Equal(t, true, result.Meta.Next)
	require.Equal(t, true, result.Meta.Prev)

}

func TestNumberService_FetchPhoneNumbersFilterByCountryAndState_Single_Page(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mockRepository := repoMock.NewMockPhoneNumberRepository(mockController)
	mockValidator := serviceMock.NewMockNumberValidator(mockController)

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
		}, nil).AnyTimes()

	mockRepository.EXPECT().FetchPaginatedPhoneNumbersByCode("237", 10, 11).Return(nil, apperror.NotFound).AnyTimes()

	mockValidator.EXPECT().Validate("(237) 23456789").Return("Cameroon", "+237", "23456789", false).AnyTimes()
	mockValidator.EXPECT().GetCodeFromCountry("Cameroon").Return("237", nil)

	numSvc := NewNumberService(mockValidator, mockRepository)

	results, err := numSvc.FilterByCountryAndState("Cameroon", "NOK", "1", "10")

	require.NoError(t, err, "Expected No Error\nGot: %v\n", err)

	for _, d := range results.Data {
		require.Equal(t, "NOK", d.State)
	}

	require.Equal(t, 10, len(results.Data))
	require.Equal(t, false, results.Meta.Next)
	require.Equal(t, false, results.Meta.Prev)
}

func TestNumberService_FilterByCountryAndState_Multi_Page(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mockRepository := repoMock.NewMockPhoneNumberRepository(mockController)
	mockValidator := serviceMock.NewMockNumberValidator(mockController)

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

	mockRepository.EXPECT().FetchPaginatedPhoneNumbersByCode("237", 10, 11).Return(nil, apperror.NotFound).AnyTimes()

	mockValidator.EXPECT().Validate("(237) 23456789").Return("Cameroon", "+237", "23456789", false).AnyTimes()
	mockValidator.EXPECT().GetCodeFromCountry("Cameroon").Return("237", nil)

	numSvc := NewNumberService(mockValidator, mockRepository)

	results, err := numSvc.FilterByCountryAndState("Cameroon", "NOK", "1", "10")

	require.NoError(t, err, "Expected No Error\nGot: %v\n", err)

	for _, d := range results.Data {
		require.Equal(t, "NOK", d.State)
	}

	require.Equal(t, len(results.Data), 10)
	require.Equal(t, true, results.Meta.Next)
	require.Equal(t, false, results.Meta.Prev)
}
