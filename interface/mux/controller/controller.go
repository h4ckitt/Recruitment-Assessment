package controller

import (
	"assessment/interface/mux/helper"
	"assessment/model"
	"assessment/service"
	"net/http"
)

type Controller struct {
	numberService *service.NumberService
}

func NewNumberController(numberService *service.NumberService) *Controller {
	return &Controller{numberService: numberService}
}

func (controller *Controller) FetchAllPhoneNumbers(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()

	limit := queries.Get("limit")
	page := queries.Get("page")
	country := queries.Get("country")
	state := queries.Get("state")

	var (
		result model.Result
		err    error
	)

	switch {
	case country == "" && state == "":
		result, err = controller.numberService.FetchPhoneNumbers(page, limit)

	case country != "" && state == "":
		result, err = controller.numberService.FilterByCountry(country, page, limit)

	case country == "" && state != "":
		result, err = controller.numberService.FilterByState(state, page, limit)

	case country != "" && page != "":
		result, err = controller.numberService.FilterByCountryAndState(country, state, page, limit)
	}

	if err != nil {
		helper.ReturnFailure(w, err)
		return
	}

	helper.ReturnSuccess(w, result)
}
