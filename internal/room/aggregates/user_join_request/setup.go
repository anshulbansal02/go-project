package roomaggregates

import "anshulbansal02/scribbly/pkg/repository"

func NewUserJoinRequestRepository(repository repository.Repository) *UserJoinRequestRepository {
	return &UserJoinRequestRepository{
		Repository: repository,
	}
}
