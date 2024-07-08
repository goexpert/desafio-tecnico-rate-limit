package database

type TokenLimit struct {
	Name  string `json:"name"`
	Limit int    `json:"limit"`
}

func HidrateListaTokens() []TokenLimit {
	lista := []TokenLimit{
		{Name: "Token010", Limit: 10},
		{Name: "Token020", Limit: 20},
		{Name: "Token050", Limit: 50},
		{Name: "Token100", Limit: 100},
	}
	return lista
}
