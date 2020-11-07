package repository

import (
	"context"
	"fmt"
	"github.com/google/logger"
	"github.com/technoZoomers/MasterHubBackend/models"
	"go.mongodb.org/mongo-driver/bson"
)

type AvatarsRepo struct {
	repository     *Repository
	userKey        string
	avatarKey      string
	extKey         string
	collectionName string
}

func (avatarsRepo *AvatarsRepo) InsertAvatar(avatar *models.AvatarDB) error {
	avatarCollection := avatarsRepo.repository.mongoDB.Collection(avatarsRepo.collectionName)
	_, err := avatarCollection.InsertOne(context.TODO(), avatar)
	if err != nil {
		dbError := fmt.Errorf("failed to insert avatar: %v", err.Error())
		logger.Errorf(dbError.Error())
		return dbError
	}
	return nil
}

func (avatarsRepo *AvatarsRepo) UpdateAvatarByUserId(userId int64, avatar *models.AvatarDB) error {
	avatarCollection := avatarsRepo.repository.mongoDB.Collection(avatarsRepo.collectionName)
	filter := bson.D{{avatarsRepo.userKey, userId}}
	update := bson.M{
		"$set": bson.M{
			avatarsRepo.extKey:    avatar.Extension,
			avatarsRepo.avatarKey: avatar.Filename,
		},
	}
	_, err := avatarCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		dbError := fmt.Errorf("failed to update avatar: %v", err.Error())
		logger.Errorf(dbError.Error())
		return dbError
	}
	return nil
}

func (avatarsRepo *AvatarsRepo) GetAvatarByUser(userId int64, avatar *models.AvatarDB) error {
	cookiesCollection := avatarsRepo.repository.mongoDB.Collection(avatarsRepo.collectionName)
	filter := bson.D{{avatarsRepo.userKey, userId}}
	err := cookiesCollection.FindOne(context.TODO(), filter).Decode(&avatar)
	if err != nil {
		dbError := fmt.Errorf("failed to retrieve avatar: %v", err.Error())
		logger.Errorf(dbError.Error())
		return dbError
	}
	return nil
}
