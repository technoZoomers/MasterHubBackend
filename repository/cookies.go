package repository

import (
	"context"
	"fmt"
	"github.com/google/logger"
	"github.com/technoZoomers/MasterHubBackend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (cookiesRepo *CookiesRepo) GetCookiesByUser(userId int64) ([]models.CookieDB, error) {
	var dbError error
	cookies := make([]models.CookieDB, 0)
	cookiesCollection := cookiesRepo.repository.mongoDB.Collection(cookiesRepo.collectionName)
	filter := bson.D{{cookiesRepo.userKey, userId}}
	findOptions := options.Find()
	cur, err := cookiesCollection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		dbError = fmt.Errorf("failed to retrieve cookies: %v", err.Error())
		logger.Errorf(dbError.Error())
		return cookies, dbError
	}

	for cur.Next(context.TODO()) {
		var cookie models.CookieDB
		err = cur.Decode(&cookie)
		if err != nil {
			dbError = fmt.Errorf("failed to retrieve cookie: %v", err.Error())
			logger.Errorf(dbError.Error())
			return cookies, dbError
		}
		cookies = append(cookies, cookie)

	}
	if err := cur.Err(); err != nil {
		dbError = fmt.Errorf("cursor error: %v", err.Error())
		logger.Errorf(dbError.Error())
		return cookies, dbError
	}
	err = cur.Close(context.TODO())
	if err != nil {
		dbError = fmt.Errorf("failed to close cursor: %v", err.Error())
		logger.Errorf(dbError.Error())
		return cookies, dbError
	}
	return cookies, nil
}

func (cookiesRepo *CookiesRepo) GetUserByCookie(cookie string, cookieDB *models.CookieDB) error {
	//log.Println(cookie)
	cookiesCollection := cookiesRepo.repository.mongoDB.Collection(cookiesRepo.collectionName)
	filter := bson.D{{cookiesRepo.cookieKey, cookie}}
	err := cookiesCollection.FindOne(context.TODO(), filter).Decode(&cookieDB)
	if err != nil {
		return err
	}
	//log.Println(cookieDB)
	return nil
}
