package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/constants"
)

// reads json body into 'dest' make sure the `dest` is never nil otherwise it will panic
func ReadJson(w http.ResponseWriter, r *http.Request, dest interface{}) error {

	if dest == nil {
		panic(errors.New("must pass non nil pointer to decode json body into"))
	}

	//create reader with size limit
	r.Body = http.MaxBytesReader(w, r.Body, constants.MAX_BODY_SIZE)

	//decode json request body into dest struct
	err := json.NewDecoder(r.Body).Decode(dest)

	return err

}

// handles common json parsing errors and writes error response
func WriteRequestPayloadError(err error, w http.ResponseWriter) {

	var maxBytesError *http.MaxBytesError
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError

	switch {
	case errors.As(err, &maxBytesError):
		WriteErrorJsonResponse(w, constants.ErrorCodes.BadRequest, http.StatusRequestEntityTooLarge, []string{"request body is too large"})
		return
	case errors.As(err, &syntaxError):
		WriteErrorJsonResponse(w, constants.ErrorCodes.BadRequest, http.StatusBadRequest, []string{fmt.Sprintf("Body contains invalid/badly formed json at %d", syntaxError.Offset)})
		return
	case errors.Is(err, io.ErrUnexpectedEOF):
		WriteErrorJsonResponse(w, constants.ErrorCodes.BadRequest, http.StatusBadRequest, []string{"request body contains badly formed json"})
		return
	case errors.As(err, &unmarshalTypeError):
		WriteErrorJsonResponse(w, constants.ErrorCodes.BadRequest, http.StatusBadRequest, []string{fmt.Sprintf("invalid field type in json body for %s", unmarshalTypeError.Field)})
		return
	case errors.Is(err, io.EOF):
		WriteErrorJsonResponse(w, constants.ErrorCodes.BadRequest, http.StatusBadRequest, []string{"body cannot be empty"})
		return
	default:
		WriteInternalErrorJsonResponse(w)
		return
	}

}
