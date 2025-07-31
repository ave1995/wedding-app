package factory

import (
	"sync"
	"wedding-app/domain/service"
	"wedding-app/domain/store"
	"wedding-app/service/user"
	"wedding-app/store/mongodb"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	userStoreOnce sync.Once
	userStore     store.UserStore
)

func getUserStore(database *mongo.Database) store.UserStore {
	userStoreOnce.Do(func() {
		userStore = mongodb.NewUserStore(database)
	})

	return userStore
}

var (
	userServiceOnce sync.Once
	userService     service.UserService
)

func getUserService(store store.UserStore) service.UserService {
	userServiceOnce.Do(func() {
		userService = user.NewUserService(store)
	})

	return userService
}
