package service

import (
	"errors"
	"strconv"
)

/*validateParams : helps validate the input parameters sent by the client
Returns:
	- page <string>
	- limit <string>
	- error <error>
*/
func validateParams(state, page, limit string) (int, int, error) {
	// give page a default value of 1 if it is empty
	if page == "" {
		page = "1"
	}

	// give limit a default value of 5 if it is empty or 0
	if limit == "" || limit == "0" {
		limit = "5"
	}

	// ensure that the page and limit are digits
	// ensure that the state values are one of OK, NOK, or empty
	if !numberRegex.MatchString(limit) || !numberRegex.MatchString(page) || (state != "OK" && state != "NOK" && state != "") {
		return -1, -1, errors.New("invalid data provided")
	}

	// convert the page and limit to integers
	pg, _ := strconv.Atoi(page)
	lim, _ := strconv.Atoi(limit)

	return pg, lim, nil
}
