package handlers

import "github.com/technoZoomers/MasterHubBackend/useCases"


type Handlers struct {
	UsersHandlers *UsersHandlers
	MastersHandlers *MastersHandlers
	StudentsHandlers *StudentsHandlers
	LanguagesHandlers *LanguagesHandlers
	ThemesHandlers *ThemesHandlers
	VideosHandlers *VideosHandlers
	AvatarsHandlers *AvatarsHandlers
}

var h Handlers

func Init(usersUC useCases.UsersUCInterface, mastersUC useCases.MastersUCInterface, studentsUC useCases.StudentsUCInterface,
	themesUC useCases.ThemesUCInterface, languagesUC useCases.LanguagesUCInterface,
	videosUC useCases.VideosUCInterface, avatarsUC useCases.AvatarsUCInterface, ) error {
	h.UsersHandlers = &UsersHandlers{usersUC}
	h.MastersHandlers = &MastersHandlers{mastersUC}
	h.StudentsHandlers = &StudentsHandlers{studentsUC}
	h.LanguagesHandlers = &LanguagesHandlers{languagesUC}
	h.ThemesHandlers = &ThemesHandlers{themesUC}
	h.AvatarsHandlers = &AvatarsHandlers{avatarsUC}
	h.VideosHandlers = &VideosHandlers{videosUC}
	return nil
}

func GetUsersH() *UsersHandlers {
	return h.UsersHandlers
}

func GetMastersH() *MastersHandlers {
	return h.MastersHandlers
}

func GetStudentsH() *StudentsHandlers {
	return h.StudentsHandlers
}

func GetLanguagesH() *LanguagesHandlers {
	return h.LanguagesHandlers
}

func GetThemesH() *ThemesHandlers {
	return h.ThemesHandlers
}

func GetVideosH() *VideosHandlers {
	return h.VideosHandlers
}

func GetAvatarsH() *AvatarsHandlers {
	return h.AvatarsHandlers
}