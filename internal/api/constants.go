package api

import "time"

const CHIRPY_LIMIT int = 140
const TOKEN_EXPIRATION_LIMIT_SECONDS_INT int = 3600

// constants in time.Duration
const TOKEN_EXPIRATION_LIMIT_SECONDS time.Duration = time.Second * 3600
const REFRESH_TOKEN_EXPIRATION_LIMIT_SECONDS time.Duration = time.Second * 5184000
