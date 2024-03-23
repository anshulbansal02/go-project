package roomaggregates

import "anshulbansal02/scribbly/pkg/repository"

func NewUserRoomRelation(repository repository.Repository) *UserRoomRelationRepository {
	return &UserRoomRelationRepository{
		Repository: repository,
	}
}
