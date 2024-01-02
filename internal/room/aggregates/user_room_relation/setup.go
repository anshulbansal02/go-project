package roomaggregates

import "anshulbansal02/scribbly/internal/repository"

func NewUserRoomRelation(repository repository.Repository) *UserRoomRelationRepository {
	return &UserRoomRelationRepository{
		Repository: repository,
	}
}
