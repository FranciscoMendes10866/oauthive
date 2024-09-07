package helpers

import "os"

var IsProd = os.Getenv("API_ENV") == "prod"
var FrontendURL = os.Getenv("FRONTEND_URL")

const AuthSessionCookie = "auth_session"
const CtxSessionID = "sessionID"
const AuthSessionMaxAge = 7 * 24 * 60 * 60  // 7 days
const AuthRenewThreshold = 2 * 24 * 60 * 60 // 2 days
