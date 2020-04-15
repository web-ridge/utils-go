package api

import (
	"net/http"
)

const PARSE_ERROR_MESSAGE = "could not parse request"
const PARSE_ERROR_CODE = http.StatusBadRequest
const PARSE_ERROR_INTERNAL_CODE = "PARSE_ERROR"

const RATE_LIMIT_MESSAGE = "Rate limit exceeded"
const RATE_LIMIT_ERROR_CODE = http.StatusTooManyRequests
const RATE_LIMIT_INTERNAL_CODE = "LIMIT_EXCEEDED"
