package config

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "postgres"
	DB_HOST     = ""
	DB_SCHEME   = "postgres"
	DB_PORT     = 5432

	REDDIT_REDIRECT_URL           = ""
	REDDIT_TOKEN_EXCHANGE_URL     = ""
	REDDIT_SEARCH_URL             = ""
	REDDIT_AUTHORIZE_URL_TEMPLATE = "" +
		"?client_id=%s" +
		"&response_type=code" +
		"&state=%s" +
		"&redirect_uri=%s" +
		"&duration=permanent" +
		"&scope=%s"

	REDDIT_SCOPE      = "read"
	REDDIT_APP_ID     = ""
	REDDIT_APP_SECRET = ""
	REDDIT_USER_AGENT = ""
	ENCRYPTION_SALT   = ""
)
