package repository

import "github.com/technoZoomers/MasterHubBackend/models"

type CookiesRepoI interface {
	InsertCookie (cookie *models.CookieDB) error
	DeleteCookie (cookie string) error
	GetCookieByUser (userId string, cookieDB *models.CookieDB) error
	GetUserByCookie (cookie string, cookieDB *models.CookieDB) error
}
