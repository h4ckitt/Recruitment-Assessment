package model

type (
	//Result : Used for storing results from operations
	Result struct {
		Data []Data `json:"data"`
		Meta Meta   `json:"meta"`
	}

	//Data : Stores phone number information
	Data struct {
		Country     string `json:"country"`
		State       string `json:"state"`
		CountryCode string `json:"countryCode"`
		PhoneNumber string `json:"phoneNumber"`
	}

	//Meta : contains pagination metadata
	Meta struct {
		CurrentPage string `json:"page"`
		Next        bool   `json:"next"`
		Prev        bool   `json:"prev"`
	}
)
