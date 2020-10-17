package useCases

import "github.com/technoZoomers/MasterHubBackend/repository"

type UseCases struct {
	UsersUC       *UsersUC
	MastersUC     *MastersUC
	StudentsUC    *StudentsUC
	ThemesUC      *ThemesUC
	LanguagesUC   *LanguagesUC
	VideosUC      *VideosUC
	AvatarsUC     *AvatarsUC
	errorMessages ErrorMessagesUC
	errorId int64
}

func (useCases *UseCases) Init(usersRepo repository.UsersRepoI, mastersRepo repository.MastersRepoI, studentsRepo repository.StudentsRepoI,
	themesRepo repository.ThemesRepoI, languagesRepo repository.LanguagesRepoI,
	videosRepo repository.VideosRepoI, avatarsRepo repository.AvatarsRepoI) error {
	useCases.UsersUC = &UsersUC{useCases, usersRepo}
	useCases.MastersUC = &MastersUC{
		useCases:useCases,
		MastersRepo:mastersRepo,
		ThemesRepo:themesRepo,
		LanguagesRepo:languagesRepo,
	mastersConfig:MastersConfig{
			qualificationMap: map[int64]string{1: "self-educated", 2: "professional"},
			educationFormatMap: map[int64]string{1: "online", 2: "live"},
			qualificationMapBackwards: map[string]int64{"self-educated":1, "professional" :2},
			educationFormatMapBackwards: map[string]int64{"online":1, "live":2},
	}}
	useCases.StudentsUC = &StudentsUC{useCases, studentsRepo}
	useCases.ThemesUC = &ThemesUC{useCases, themesRepo}
	useCases.LanguagesUC = &LanguagesUC{useCases, languagesRepo}
	useCases.VideosUC = &VideosUC{
		useCases:    useCases,
		VideosRepo:  videosRepo,
		MastersRepo: mastersRepo,
		ThemesRepo:  themesRepo,
		videosConfig: VideoConfig{
			videosDefaultName: "noname",
			videosDir:           "./master_videos/",
			videoPrefixMaster: "master_",
			videoPrefixVideo: "video_",
			videoPrefixIntro: "intro",
		},
	}
	useCases.AvatarsUC = &AvatarsUC{useCases, avatarsRepo}
	useCases.errorMessages = ErrorMessagesUC{
		DbError: "database internal error",
		FileErrors: FileErrors{
			FileOpenError:          "error opening file",
			FileReadError:          "error reading file",
			FileReadExtensionError: "error reading file extension",
			FileCreateError:        "error creating file",
			FileRemoveError: "error removing file",
		},
	}
	useCases.errorId = 0
	return nil
}
