package handlers

import "github.com/technoZoomers/MasterHubBackend/useCases"

type Handlers struct {
	UsersHandlers     *UsersHandlers
	MastersHandlers   *MastersHandlers
	StudentsHandlers  *StudentsHandlers
	LanguagesHandlers *LanguagesHandlers
	ThemesHandlers    *ThemesHandlers
	VideosHandlers    *VideosHandlers
	AvatarsHandlers   *AvatarsHandlers
}

func (handlers *Handlers) Init(usersUC useCases.UsersUCInterface, mastersUC useCases.MastersUCInterface, studentsUC useCases.StudentsUCInterface,
	themesUC useCases.ThemesUCInterface, languagesUC useCases.LanguagesUCInterface,
	videosUC useCases.VideosUCInterface, avatarsUC useCases.AvatarsUCInterface) error {
	handlers.UsersHandlers = &UsersHandlers{usersUC}
	handlers.MastersHandlers = &MastersHandlers{mastersUC}
	handlers.StudentsHandlers = &StudentsHandlers{studentsUC}
	handlers.LanguagesHandlers = &LanguagesHandlers{languagesUC}
	handlers.ThemesHandlers = &ThemesHandlers{themesUC}
	handlers.AvatarsHandlers = &AvatarsHandlers{avatarsUC}
	handlers.VideosHandlers = &VideosHandlers{
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
