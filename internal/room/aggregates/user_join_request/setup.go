package roomaggregates

import "anshulbansal02/scribbly/internal/repository"

func NewUserJoinRequestRepository(repository repository.Repository) *UserJoinRequestRepository {
	return &UserJoinRequestRepository{
		Repository: repository,
	}
}
