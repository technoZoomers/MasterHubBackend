package repository

import (
	"context"
	"github.com/technoZoomers/MasterHubBackend/models"
	"go.mongodb.org/mongo-driver/bson"
)

type CookiesRepo struct {
	repository     *Repository
	userKey        string
	cookieKey      string
	collectionName string
}

func (cookiesRepo *CookiesRepo) InsertCookie(cookie *models.CookieDB) error {
	//log.Println(cookie)
	cookiesCollection := cookiesRepo.repository.mongoDB.Collection(cookiesRepo.collectionName)
	_, err := cookiesCollection.InsertOne(context.TODO(), cookie)
	return err
}

func (cookiesRepo *CookiesRepo) DeleteCookie(cookie string) error {
	cookiesCollection := cookiesRepo.repository.mongoDB.Collection(cookiesRepo.collectionName)
	filter := bson.D{{cookiesRepo.cookieKey, cookie}}
	_, err := cookiesCollection.DeleteOne(context.TODO(), filter)
	return err
}

func (cookiesRepo *CookiesRepo) GetCookieByUser(userId int64, cookieDB *models.CookieDB) error {
	cookiesCollection := cookiesRepo.repository.mongoDB.Collection(cookiesRepo.collectionName)
	filter := bson.D{{cookiesRepo.userKey, userId}}
	err := cookiesCollection.FindOne(context.TODO(), filter).Decode(&cookieDB)
	if err != nil {
		return err
	}
	return nil
}

func (cookiesRepo *CookiesRepo) GetUserByCookie(cookie string, cookieDB *models.CookieDB) error {
	//log.Println(cookie)
	cookiesCollection := cookiesRepo.repository.mongoDB.Collection(cookiesRepo.collectionName)
	filter := bson.D{{cookiesRepo.cookieKey, cookie}}
	err := cookiesCollection.FindOne(context.TODO(), filter).Decode(&cookieDB)
	if err != nil {
		return err
	}
	return nil
}
