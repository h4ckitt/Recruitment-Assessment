package service

import (
	"assessment/apperror"
	"assessment/model"
	"assessment/repository"
	"log"
	"regexp"
	"strings"
)

var numberRegex = regexp.MustCompile(`^\d+$`)

type NumberService struct {
	validator  NumberValidator
	repository repository.PhoneNumberRepository
}

/*NewNumberService : This starts a new service which handles the business logic of returning
phone numbers with the specified criteria
*/
func NewNumberService(validator NumberValidator, repository repository.PhoneNumberRepository) *NumberService {
	return &NumberService{
		validator:  validator,
		repository: repository,
	}
}

/*FetchPhoneNumbers : Fetches all the phone numbers in the database
Returns a paginated list of all phone numbers in the database.
*/
func (s *NumberService) FetchPhoneNumbers(page, limit string) (model.Result, error) {
	off, lim, err := validateParams("", page, limit) // validates the input parameters

	// return an error if an unsupported parameter is received
	if err != nil {
		return model.Result{}, apperror.BadRequest
	}

	// calculate the offset to be used for fetching subsequent
	// e.g page 2 with a limit of 5 per page will begin search from position 5 in the database
	off = lim*off - lim

	// fetch the requested phone numbers from the database using value of specified limit + 1.
	// the reason for this is to simulate a lookahead for ensuring that there's still more data even after the requested limit is satisfied
	result, err := s.repository.FetchPaginatedPhoneNumbers(off, lim+1)

	// ensure that no error was returned
	// this would typically be a serious error such as db outage or unavailability
	if err != nil {
		log.Println(err)                            // log the error
		return model.Result{}, apperror.ServerError // return an internal server error
	}

	// declare variable for holding result metadata
	var meta model.Meta

	// check whether the returned data has an extra data that serves as lookahead.
	// existence of this extra data informs that there is still more data to be read from the db
	if len(result) == lim+1 {
		result = result[:lim]
		meta.Next = true
	}

	// if the offset is greater than the limit then we definitely aren't on the first page
	if off >= lim {
		meta.Prev = true

	}

	meta.CurrentPage = page

	var data []model.Data

	// load the data from the db into the result object
	for _, number := range result {
		country, code, number, valid := s.validator.Validate(number)

		state := "OK"

		// if the state of the phone number is not valid then set it as Not Okay (NOK)
		if !valid {
			state = "NOK"
		}

		data = append(data, model.Data{
			Country:     country,
			CountryCode: code,
			PhoneNumber: number,
			State:       state,
		})
	}

	// if an empty result set was returned them there's no next or previous.
	// this would usually be populated if a client tries to fetch past the available pages for a resource
	// i.e. total page is 5 and client tries to fetch 7.
	if len(data) == 0 {
		meta.Next = false
		meta.Prev = false
	}

	// prepare the final result object
	finalResult := model.Result{
		Data: data,
		Meta: meta,
	}

	return finalResult, nil
}

//FilterByState : Filter phone numbers by the validity specified by the client
func (s *NumberService) FilterByState(state, page, limit string) (model.Result, error) {
	pg, lim, err := validateParams(state, page, limit) // validate the input

	// return an error if unacceptable input is returned
	if err != nil {
		return model.Result{}, apperror.BadRequest
	}

	// The offset here is 0 instead of being calculated (further details below)
	offset := 0

	// declare variables required for storing crucial data
	var (
		data        = make([]model.Data, lim)
		meta        model.Meta
		numbers     []string
		next        bool
		currentPage = 1
	)

	/*Fetching paginated data from the database based on state is a bit more complex because that data is non-contiguous. e.g:
		1. OK
		2. NOK
		3. NOK
		4. OK
		5. NOK
		6. OK
		7. OK
		8. NOK

	Fetching page 1 of phone numbers which are nok with a limit of 3 puts the final offset for page 1 on number 5 because that's where the third NOK number is.

	In order to fetch page 2, the offset will need to start from number 6 to prevent number 5 from getting included the result set.

	What this function does is start from page 1 up to page n and try to calculate the start offset to use for page n since it's not really deterministic.
	*/

	// loop for as long as the current page is lesser than the page for which data is needed
	// current page always starts as 1.
	for currentPage <= pg {
		// Check whether the length of the data is lesser than the limit
		// this usually executes when the length of data is less than the limit required meaning there's no longer data to be read.
		if len(data) < lim {
			data = nil
			break
		}

		// recover excess data i.e. data after the required limit.
		// e.g. page 4 returns 1, 2, 3, 4, 5 and our limit is 3
		// 4, 5 need to be saved because the database query will start from a further offset and returned data will not include them.
		data = data[lim:]

		// loop for as long as the length of the data we have is lesser than the limit
		for len(data) < lim {

			numbers, err = s.repository.FetchPaginatedPhoneNumbers(offset, lim+1)

			if err != nil {
				log.Println(err)
				return model.Result{}, apperror.ServerError
			}

			// if returned data from the db is empty
			if len(numbers) == 0 {
				next = false // there will be no next page
				break
			}

			// if we get extra (lookahead) data then there's extra data to be read
			if len(numbers) == lim+1 {
				next = true
				numbers = numbers[:lim] // truncate the returned data
			} else {
				next = false // if we don't get extra (lookahead) data then there's no further data to be read
			}

			// fill in the data slice with results from the database.
			for _, number := range numbers {
				country, code, number, valid := s.validator.Validate(number)

				// ensure that phone number status matches the requested status
				// i.e. ensure that current phone number is OK and requested status is OK.
				if (state == "OK" && valid) || (state == "NOK" && !valid) {
					data = append(data, model.Data{
						Country:     country,
						CountryCode: code,
						PhoneNumber: number,
						State:       state,
					})
				}
			}

			// increase the offset
			offset += lim
		}
		// if the length of data returned is greater than the limit requested, then there's more data to be read from the db.
		// due to the non-contiguous nature of the data in the db, it's possible for the length of data to be greater than limit.
		if len(data) > lim {
			next = true
		}
		// increment the current page
		currentPage++
	}

	// make the pagination indicators false if the length of data returned is 0
	if len(data) == 0 {
		meta.Next = false
		meta.Prev = false
	} else {
		// if the length of returned data is greater than the limit then truncate it to the requested limit
		if len(data) > lim {
			data = data[:lim]
		}
		meta.Next = next
		// set the previous flag if the offset is greater than the limit.
		if (pg*lim - lim) >= lim {
			meta.Prev = true
		}
	}

	meta.CurrentPage = page

	return model.Result{
		Data: data,
		Meta: meta,
	}, nil

}

//FilterByCountry : Filter Numbers From The Database By The Country They Belong To.
func (s *NumberService) FilterByCountry(country, page, limit string) (model.Result, error) {

	p, lim, err := validateParams("", page, limit)

	if err != nil {
		return model.Result{}, apperror.BadRequest
	}

	// get the code for the specified country since that's what will be used for the database query
	code, err := s.validator.GetCodeFromCountry(country)

	if err != nil {
		return model.Result{}, apperror.NotFound
	}

	// calculate the offset to be used for fetching data for the requested page
	offset := lim*p - lim

	var (
		state   string
		data    []model.Data
		meta    model.Meta
		numbers []string
		next    bool
	)

	// loop until the length of returned data is same as the length of the requested limit
	for len(data) < lim {
		numbers, err = s.repository.FetchPaginatedPhoneNumbersByCode(code, offset, lim+1)

		if err != nil {
			log.Println(err)
			return model.Result{}, apperror.ServerError
		}

		// if results from the database are empty break out of the loop
		if len(numbers) == 0 {
			break
		}

		// if lookahead data is present then set the next flag to true
		if len(numbers) == lim+1 {
			next = true
			numbers = numbers[:lim] // truncate returned data
		} else {
			next = false
		}

		// fill in results from the db into the data object
		for _, number := range numbers {
			country, code, number, valid := s.validator.Validate(number)
			if valid {
				state = "OK"
			} else {
				state = "NOK"
			}
			data = append(data, model.Data{
				Country:     country,
				CountryCode: code,
				PhoneNumber: number,
				State:       state,
			})
		}
		// increase the offset
		offset += lim
	}

	meta.Next = next

	// use the offset to determine whether there should be a page behind
	if (p*lim - lim) >= lim {
		meta.Prev = true
	}

	// strip the pagination flags if returned data is empty
	if len(data) == 0 {
		meta.Next = false
		meta.Prev = false
	}

	meta.CurrentPage = page

	return model.Result{
		Data: data,
		Meta: meta,
	}, nil
}

//FilterByCountryAndState : Filter phone numbers based on the specified country and data
func (s *NumberService) FilterByCountryAndState(country, state, page, limit string) (model.Result, error) {

	p, lim, err := validateParams(state, page, limit)

	if err != nil {
		return model.Result{}, apperror.NotFound
	}

	code, err := s.validator.GetCodeFromCountry(country)

	if err != nil {
		return model.Result{}, apperror.NotFound
	}

	// switch the state variable to uppercase
	state = strings.ToUpper(state)

	// set the offset to 0 because that's what page 1 resolves to (detailed explanation below)
	offset := 0

	var (
		data        = make([]model.Data, lim)
		meta        model.Meta
		numbers     []string
		next        bool
		currentPage = 1
	)

	/*Fetching paginated data from the database based on status is a bit more complex because that data is non-contiguous. e.g:
		1. OK
		2. NOK
		3. NOK
		4. OK
		5. NOK
		6. OK
		7. OK
		8. NOK

	Fetching page 1 of phone numbers which are nok with a limit of 3 puts the final offset for page 1 on number 5 because that's where the third NOK number is.

	In order to fetch page 2, the offset will need to start from number 6 in order to prevent number 5 from entering the result set.

	What this function does is to start from page 1 up to page n and try to calculate the start offset to use for page n since it's not really deterministic.

	This block runs in 0(N^2) Time
	*/

	// loop until the current page is equal to the required page
	for currentPage <= p {
		// Results From The Previous Page Are Lesser Than The Limit, Meaning There's No More Data To Be Read From The DB.
		if len(data) < lim {
			data = nil
			break
		}

		// recover excess data i.e. data after the required limit.
		// e.g. page 4 returns 1, 2, 3, 4, 5 and our limit is 3
		// 4, 5 need to be saved because the database query will start from a further offset and returned data will not include them.
		data = data[lim:]

		// loop until the data is equal to the limit
		for len(data) < lim {
			numbers, err = s.repository.FetchPaginatedPhoneNumbersByCode(code, offset, lim+1)

			if err != nil {
				log.Println(err)
				return model.Result{}, apperror.ServerError
			}

			// there will be no next page if returned result from the db is empty, so exit the loop
			if len(numbers) == 0 {
				next = false
				break
			}

			// if we get extra (lookahead) data then there's extra data to be read
			if len(numbers) == lim+1 {
				next = true
				numbers = numbers[:lim] // truncate the returned data
			} else {
				next = false
			}

			// fill in the data slice with results from the database.
			for _, number := range numbers {
				country, code, number, valid := s.validator.Validate(number)

				// ensure that phone number status matches the requested status
				// i.e. ensure that current phone number is OK and requested status is OK.
				if (state == "OK" && valid) || (state == "NOK" && !valid) {
					data = append(data, model.Data{
						Country:     country,
						CountryCode: code,
						PhoneNumber: number,
						State:       state,
					})
				}
			}
			// calculate the offset for the next page (offset + limit)
			offset += lim
		}
		if len(data) > lim {
			next = true
		}

		// go on to the next page
		currentPage++
	}

	// make the pagination indicators false if the length of data returned is 0
	if len(data) == 0 {
		meta.Next = false
		meta.Prev = false
	} else {
		// if the length of returned data is greater than the limit then truncate it to the requested limit
		if len(data) > lim {
			data = data[:lim]
		}
		meta.Next = next

		// set the previous flag if the offset is greater than the limit.
		if (p*lim - lim) >= lim {
			meta.Prev = true
		}
	}

	meta.CurrentPage = page

	return model.Result{
		Data: data,
		Meta: meta,
	}, nil
}
