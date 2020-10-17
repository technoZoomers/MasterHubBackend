package handlers

import (
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/useCases"
)

type Handlers struct {
	UsersHandlers     *UsersHandlers
	MastersHandlers   *MastersHandlers
	StudentsHandlers  *StudentsHandlers
	LanguagesHandlers *LanguagesHandlers
	ThemesHandlers    *ThemesHandlers
	VideosHandlers    *VideosHandlers
	AvatarsHandlers   *AvatarsHandlers
	badRequestError *models.BadRequestError
}

func (handlers *Handlers) Init(usersUC useCases.UsersUCInterface, mastersUC useCases.MastersUCInterface, studentsUC useCases.StudentsUCInterface,
	themesUC useCases.ThemesUCInterface, languagesUC useCases.LanguagesUCInterface,
	videosUC useCases.VideosUCInterface, avatarsUC useCases.AvatarsUCInterface) error {
	handlers.UsersHandlers = &UsersHandlers{handlers, usersUC}
	handlers.MastersHandlers = &MastersHandlers{handlers, mastersUC}
	handlers.StudentsHandlers = &StudentsHandlers{handlers, studentsUC}
	handlers.LanguagesHandlers = &LanguagesHandlers{handlers, languagesUC}
	handlers.ThemesHandlers = &ThemesHandlers{handlers, themesUC}
	handlers.AvatarsHandlers = &AvatarsHandlers{handlers, avatarsUC}
	handlers.VideosHandlers = &VideosHandlers{
		handlers:handlers,
		VideosUC:videosUC,
	VideoParseConfig:VideoParseConfig{
		FormDataKey:  "video",
		VideoFormats: map[string]bool{
			"video/webm":               true,
			"audio/ogg":                true,
			"video/mp4":                true,
			"video/quicktime":          true,
			"video/x-msvideo":          true,
			"application/octet-stream": true,
		},
	}}
	return nil
}
