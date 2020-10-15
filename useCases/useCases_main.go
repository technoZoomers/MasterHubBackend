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
	mastersConfig:MastersConfig{qualificationMap: map[int64]string{1: "self-educated", 2: "professional"},
		educationFormatMap: map[int64]string{1: "online", 2: "live"}}}
	useCases.StudentsUC = &StudentsUC{useCases, studentsRepo}
	useCases.ThemesUC = &ThemesUC{useCases, themesRepo}
	useCases.LanguagesUC = &LanguagesUC{useCases, languagesRepo}
	useCases.VideosUC = &VideosUC{
		useCases:    useCases,
		VideosRepo:  videosRepo,
		MastersRepo: mastersRepo,
		ThemesRepo:  themesRepo,
		videosConfig: VideoConfig{videosDefaultName: "noname",
			videosDir:           "./master_videos/",
			videoFilenamePrefix: "master_video_"},
	}
	useCases.AvatarsUC = &AvatarsUC{useCases, avatarsRepo}
	useCases.errorMessages = ErrorMessagesUC{
		DbError: "database internal error",
		FileErrors: FileErrors{
			FileOpenError:          "error opening file",
			FileReadError:          "error reading file",
			FileReadExtensionError: "error reading file extension",
			FileCreateError:        "error creating file",
		},
	}
	return nil
}
