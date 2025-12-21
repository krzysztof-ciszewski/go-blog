package query

import (
	"context"
	view "main/internal/Application/View"
	repository "main/internal/Domain/Repository"
)

type FindUserByQueryHandler struct {
	UserRepository repository.UserRepository
}

func (h FindUserByQueryHandler) Handle(ctx context.Context, query any) (any, error) {
	findUserQuery, ok := query.(FindUserByQuery)
	if !ok {
		return view.UserView{}, nil
	}

	userEntity, err := h.UserRepository.FindByProviderUserIdAndEmail(findUserQuery.Filters.ProviderUserId, findUserQuery.Filters.UserEmail)
	if err != nil {
		return view.UserView{}, err
	}

	return view.NewUserView(
		userEntity.Id(),
		userEntity.CreatedAt(),
		userEntity.UpdatedAt(),
		userEntity.Email,
		userEntity.Provider,
		userEntity.Name,
		userEntity.FirstName,
		userEntity.LastName,
		userEntity.ProviderUserId,
		userEntity.AvatarURL,
	), nil
}

func (h FindUserByQueryHandler) Supports(query any) bool {
	_, ok := query.(FindUserByQuery)
	return ok
}
