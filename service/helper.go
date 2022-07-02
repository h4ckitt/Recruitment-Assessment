package service

import (
	"errors"
	"strconv"
)

func validateParams(state, page, limit string) (int, int, error) {
	if page == "" {
		page = "1"
	}

	if limit == "" || limit == "0" {
		limit = "5"
	}

	if !numberRegex.MatchString(limit) || !numberRegex.MatchString(page) || (state != "OK" && state != "NOK" && state != "") {
		return -1, -1, errors.New("invalid data provided")
	}

	off, _ := strconv.Atoi(page)
	lim, _ := strconv.Atoi(limit)

	return off, lim, nil
}
