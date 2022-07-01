package model

type (
	Result struct {
		Data []Data `json:"data"`
		Meta Meta   `json:"meta"`
	}

	Data struct {
		Country     string `json:"country"`
		State       string `json:"state"`
		CountryCode string `json:"countryCode"`
		PhoneNumber string `json:"phoneNumber"`
	}

	Meta struct {
		CurrentPage string `json:"page"`
		Next        bool   `json:"next"`
		Prev        bool   `json:"prev"`
	}
)
