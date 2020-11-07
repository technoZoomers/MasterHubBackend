package useCases

import (
	"bufio"
	"fmt"
	"github.com/google/logger"
	"github.com/h2non/filetype"
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/repository"
	"io/ioutil"
	"mime/multipart"
	"os"
)

type AvatarsUC struct {
	useCases     *UseCases
	AvatarsRepo  repository.AvatarsRepoI
	avatarConfig AvatarConfig
}

type AvatarConfig struct {
	avatarsDir    string
	avatarPrefix  string
	avatarPostfix string
}

func (avatarsUC *AvatarsUC) createAvatarFilename(userId int64) (string, error) {
	var avatarExists models.AvatarDB
	var filename string
	_ = avatarsUC.AvatarsRepo.GetAvatarByUser(userId, &avatarExists)
	if avatarExists.User != avatarsUC.useCases.errorId {
		return filename, &models.ConflictError{Message: "avatar already exists"}
	}
	filename = fmt.Sprintf("%s%d%s", avatarsUC.avatarConfig.avatarPrefix, userId, avatarsUC.avatarConfig.avatarPostfix)
	return filename, nil
}

func (avatarsUC *AvatarsUC) createAvatarFile(file multipart.File, filename string) (string, string, error) {
	var newPath string
	var ext string
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fileError := fmt.Errorf("%s: %s", avatarsUC.useCases.errorMessages.FileErrors.FileReadError, err.Error())
		logger.Errorf(fileError.Error())
		return newPath, ext, fileError
	}
	defer file.Close()
	fileExtension, err := filetype.Match(fileBytes)
	ext = fileExtension.Extension
	if err != nil {
		fileError := fmt.Errorf("%s: %s", avatarsUC.useCases.errorMessages.FileErrors.FileReadExtensionError, err.Error())
		logger.Errorf(fileError.Error())
		return newPath, ext, fileError
	}

	newPath = fmt.Sprintf("%s%s%s.%s", avatarsUC.useCases.filesDir, avatarsUC.avatarConfig.avatarsDir, filename, fileExtension.Extension)
	newFile, err := os.Create(newPath)
	if err != nil {
		fileError := fmt.Errorf("%s: %s", avatarsUC.useCases.errorMessages.FileErrors.FileCreateError, err.Error())
		logger.Errorf(fileError.Error())
		return newPath, ext, fileError
	}
	defer newFile.Close()

	_, err = newFile.Write(fileBytes)
	if err != nil {
		os.Remove(newPath)
		fileError := fmt.Errorf("%s: %s", avatarsUC.useCases.errorMessages.FileErrors.FileCreateError, err.Error())
		logger.Errorf(fileError.Error())
		return newPath, ext, fileError
	}
	return newPath, ext, nil
}

func (avatarsUC *AvatarsUC) NewUserAvatar(file multipart.File, userId int64) error {
	err := avatarsUC.useCases.UsersUC.validateUserId(userId)
	if err != nil {
		return err
	}

	filename, err := avatarsUC.createAvatarFilename(userId)
	if err != nil {
		return err
	}
	newPath, ext, err := avatarsUC.createAvatarFile(file, filename)

	avatarDB := models.AvatarDB{
		Filename:  filename,
		Extension: ext,
		User:      userId,
	}
	err = avatarsUC.AvatarsRepo.InsertAvatar(&avatarDB)
	if err != nil {
		os.Remove(newPath)
		return fmt.Errorf(avatarsUC.useCases.errorMessages.DbError)
	}
	return nil
}

func (avatarsUC *AvatarsUC) deleteAvatarFile(filename string, ext string) error {
	oldPath := fmt.Sprintf("%s%s%s.%s", avatarsUC.useCases.filesDir, avatarsUC.avatarConfig.avatarsDir, filename, ext)
	err := os.Remove(oldPath)
	if err != nil {
		fileError := fmt.Errorf("%s: %s", avatarsUC.useCases.errorMessages.FileErrors.FileRemoveError, err.Error())
		logger.Errorf(fileError.Error())
		return fileError
	}
	return nil
}

func (avatarsUC *AvatarsUC) ChangeUserAvatar(file multipart.File, userId int64) error {
	err := avatarsUC.useCases.UsersUC.validateUserId(userId)
	if err != nil {
		return err
	}
	var avatar models.AvatarDB
	err = avatarsUC.AvatarsRepo.GetAvatarByUser(userId, &avatar)
	if err != nil {
		return fmt.Errorf(avatarsUC.useCases.errorMessages.DbError)
	}
	if avatar.User == avatarsUC.useCases.errorId {
		absenceError := &models.NoContentError{Message: "avatar doesn't exist", RequestId: userId}
		logger.Errorf(absenceError.Error())
		return absenceError
	}
	err = avatarsUC.deleteAvatarFile(avatar.Filename, avatar.Extension)
	if err != nil {
		return err
	}
	_, ext, err := avatarsUC.createAvatarFile(file, avatar.Filename)
	if err != nil {
		return err
	}
	avatar.Extension = ext
	err = avatarsUC.AvatarsRepo.UpdateAvatarByUserId(userId, &avatar)
	if err != nil {
		return fmt.Errorf(avatarsUC.useCases.errorMessages.DbError)
	}
	return nil
}

func (avatarsUC *AvatarsUC) GetUserAvatar(userId int64) ([]byte, error) {
	var videoBytes []byte
	var filename string
	err := avatarsUC.useCases.UsersUC.validateUserId(userId)
	if err != nil {
		return videoBytes, err
	}
	var avatar models.AvatarDB
	err = avatarsUC.AvatarsRepo.GetAvatarByUser(userId, &avatar)
	if err != nil {
		return videoBytes, fmt.Errorf(avatarsUC.useCases.errorMessages.DbError)
	}
	if avatar.User == avatarsUC.useCases.errorId {
		filename = fmt.Sprintf("%s%s%s.%s", avatarsUC.useCases.filesDir, avatarsUC.avatarConfig.avatarsDir, "default", "png")
	} else {
		filename = fmt.Sprintf("%s%s%s.%s", avatarsUC.useCases.filesDir, avatarsUC.avatarConfig.avatarsDir, avatar.Filename, avatar.Extension)
	}
	imageFile, err := os.Open(filename)
	if err != nil {
		fileError := fmt.Errorf("%s: %s", avatarsUC.useCases.errorMessages.FileErrors.FileOpenError, err.Error())
		logger.Errorf(fileError.Error())
		return videoBytes, fileError
	}
	defer imageFile.Close()

	reader := bufio.NewReader(imageFile)
	imgFileInfo, err := imageFile.Stat()
	if err != nil {
		fileError := fmt.Errorf("%s: %s", avatarsUC.useCases.errorMessages.FileErrors.FileOpenError, err.Error())
		logger.Errorf(fileError.Error())
		return videoBytes, fileError
	}
	imgFileSize := imgFileInfo.Size()

	videoBytes = make([]byte, imgFileSize)
	_, err = reader.Read(videoBytes)
	if err != nil {
		fileError := fmt.Errorf("%s: %s", avatarsUC.useCases.errorMessages.FileErrors.FileReadError, err.Error())
		logger.Errorf(fileError.Error())
		return videoBytes, fileError
	}
	return videoBytes, nil
}
