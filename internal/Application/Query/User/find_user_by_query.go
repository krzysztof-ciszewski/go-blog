package user_query

type Filters struct {
	ProviderUserId string
	UserEmail      string
}

type FindUserByQuery struct {
	Filters Filters
}

func NewFindUserByQuery(providerUserId string, userEmail string) FindUserByQuery {
	return FindUserByQuery{Filters: Filters{ProviderUserId: providerUserId, UserEmail: userEmail}}
}
