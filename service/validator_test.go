package service

import (
	"assessment/model"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidator_GetCodeFromCountry(t *testing.T) {
	v := NewValidator()

	var expected = []string{
		"237",
		"251",
		"212",
		"258",
		"256",
	}

	var countries = []string{
		"Cameroon",
		"Ethiopia",
		"Morocco",
		"Mozambique",
		"Uganda",
	}

	for index, country := range countries {
		code, err := v.GetCodeFromCountry(country)
		require.NoError(t, err, "Expected No Error\nGot: %v\n", err)
		require.Equal(t, expected[index], code)
	}

	_, err := v.GetCodeFromCountry("Nigeria")
	require.Error(t, err, "Expected An Error\nGet: %v\n", err)
}

func TestValidator_Validate(t *testing.T) {
	v := NewValidator()

	var testCases = []string{
		"(212) 698054317",
		"(212) 6546545369",
		"(212) 6617344445",
		"(258) 847651504",
		"(258) 846565883",
		"(258) 849181828",
		"(256) 7503O6263",
		"(256) 704244430",
		"(251) 9773199405",
		"(237) 697151594",
	}

	var expected = []model.Data{
		{
			Country:     "Morocco",
			State:       "OK",
			CountryCode: "+212",
			PhoneNumber: "698054317",
		},
		{
			Country:     "Morocco",
			State:       "NOK",
			CountryCode: "+212",
			PhoneNumber: "6546545369",
		},
		{
			Country:     "Morocco",
			State:       "NOK",
			CountryCode: "+212",
			PhoneNumber: "6617344445",
		},
		{
			Country:     "Mozambique",
			State:       "OK",
			CountryCode: "+258",
			PhoneNumber: "847651504",
		},
		{
			Country:     "Mozambique",
			State:       "OK",
			CountryCode: "+258",
			PhoneNumber: "846565883",
		},
		{
			Country:     "Mozambique",
			State:       "OK",
			CountryCode: "+258",
			PhoneNumber: "849181828",
		},
		{
			Country:     "Uganda",
			State:       "NOK",
			CountryCode: "+256",
			PhoneNumber: "7503O6263",
		},
		{
			Country:     "Uganda",
			State:       "OK",
			CountryCode: "+256",
			PhoneNumber: "704244430",
		},
		{
			Country:     "Ethiopia",
			State:       "NOK",
			CountryCode: "+251",
			PhoneNumber: "9773199405",
		},
		{
			Country:     "Cameroon",
			State:       "OK",
			CountryCode: "+237",
			PhoneNumber: "697151594",
		},
	}

	for index, tCase := range testCases {
		country, code, number, valid := v.Validate(tCase)
		require.Equal(t, expected[index].Country, country)
		if expected[index].State == "OK" {
			require.Equal(t, true, valid)
		} else {
			require.Equal(t, false, valid)
		}
		require.Equal(t, expected[index].CountryCode, code)
		require.Equal(t, expected[index].PhoneNumber, number)
	}
}
