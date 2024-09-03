package helpers

import "os"

var IsProd = os.Getenv("API_ENV") == "prod"
var FrontendURL = os.Getenv("FRONTEND_URL")
