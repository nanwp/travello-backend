package service

type FavoriteService interface {
	Create(userId string, destinationId string) (bool, error)
	Delete()
}
