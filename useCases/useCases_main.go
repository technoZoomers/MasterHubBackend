package useCases

import "github.com/technoZoomers/MasterHubBackend/repository"

type UseCases struct {
	UsersUC     *UsersUC
	MastersUC   *MastersUC
	StudentsUC  *StudentsUC
	ThemesUC    *ThemesUC
	LanguagesUC *LanguagesUC
	VideosUC    *VideosUC
	AvatarsUC   *AvatarsUC
}

var uc UseCases

func Init(usersRepo repository.UsersRepoI, mastersRepo repository.MastersRepoI, studentsRepo repository.StudentsRepoI,
	themesRepo repository.ThemesRepoI, languagesRepo repository.LanguagesRepoI,
	videosRepo repository.VideosRepoI, avatarsRepo repository.AvatarsRepoI) error {
	uc.UsersUC = &UsersUC{usersRepo}
	uc.MastersUC = &MastersUC{mastersRepo, themesRepo, languagesRepo}
	uc.StudentsUC = &StudentsUC{studentsRepo}
	uc.ThemesUC = &ThemesUC{themesRepo}
	uc.LanguagesUC = &LanguagesUC{languagesRepo}
	uc.VideosUC = &VideosUC{videosRepo, mastersRepo}
	uc.AvatarsUC = &AvatarsUC{avatarsRepo}
	return nil
}

func GetUsersUC() UsersUCInterface {
	return uc.UsersUC
}

func GetMastersUC() MastersUCInterface {
	return uc.MastersUC
}

func GetStudentsUC() StudentsUCInterface {
	return uc.StudentsUC
}

func GetThemesUC() ThemesUCInterface {
	return uc.ThemesUC
}

func GetLanguagesUC() LanguagesUCInterface {
	return uc.LanguagesUC
}

func GetVideosUC() VideosUCInterface {
	return uc.VideosUC
}

func GetAvatarsUC() AvatarsUCInterface {
	return uc.AvatarsUC
}
