package useCases

import (
	"github.com/technoZoomers/MasterHubBackend/repository"
	"time"
)

type UseCases struct {
	UsersUC       *UsersUC
	MastersUC     *MastersUC
	StudentsUC    *StudentsUC
	ThemesUC      *ThemesUC
	LanguagesUC   *LanguagesUC
	VideosUC      *VideosUC
	AvatarsUC     *AvatarsUC
	ChatsUC       *ChatsUC
	WebsocketsUC  *WebsocketsUC
	LessonsUC     *LessonsUC
	errorMessages ErrorMessagesUC
	errorId       int64
	filesDir      string
}

func (useCases *UseCases) Init(usersRepo repository.UsersRepoI, mastersRepo repository.MastersRepoI, studentsRepo repository.StudentsRepoI,
	themesRepo repository.ThemesRepoI, languagesRepo repository.LanguagesRepoI,
	videosRepo repository.VideosRepoI, avatarsRepo repository.AvatarsRepoI, chatsRepo repository.ChatsRepoI,
	wsRepo repository.WebsocketsRepo, cookiesRepo repository.CookiesRepoI, lessonsRepo repository.LessonsRepoI) error {
	useCases.UsersUC = &UsersUC{useCases, usersRepo, cookiesRepo}
	useCases.MastersUC = &MastersUC{
		useCases:      useCases,
		MastersRepo:   mastersRepo,
		UsersRepo:     usersRepo,
		StudentsRepo:  studentsRepo,
		ThemesRepo:    themesRepo,
		LanguagesRepo: languagesRepo,
		mastersConfig: MastersConfig{
			qualificationMap:            map[int64]string{1: "самоучка", 2: "профессионал"},
			educationFormatMap:          map[int64]string{1: "онлайн", 2: "вживую"},
			qualificationMapBackwards:   map[string]int64{"самоучка": 1, "профессионал": 2},
			educationFormatMapBackwards: map[string]int64{"онлайн": 1, "вживую": 2},
		}}
	useCases.StudentsUC = &StudentsUC{useCases, usersRepo, mastersRepo, studentsRepo, languagesRepo}
	useCases.ThemesUC = &ThemesUC{useCases, themesRepo}
	useCases.LanguagesUC = &LanguagesUC{useCases, languagesRepo}
	useCases.VideosUC = &VideosUC{
		useCases:    useCases,
		VideosRepo:  videosRepo,
		MastersRepo: mastersRepo,
		ThemesRepo:  themesRepo,
		videosConfig: VideoConfig{
			videosDefaultName: "noname",
			videosDir:         "/master_videos/",
			videoPrefixMaster: "master_",
			videoPrefixVideo:  "_video_",
			videoPrefixIntro:  "_intro",
		},
	}
	useCases.AvatarsUC = &AvatarsUC{
		useCases:    useCases,
		AvatarsRepo: avatarsRepo,
		avatarConfig: AvatarConfig{
			avatarsDir:    "/users_avatars/",
			avatarPrefix:  "user_",
			avatarPostfix: "_avatar",
		},
	}
	useCases.ChatsUC = &ChatsUC{
		useCases:     useCases,
		MastersRepo:  mastersRepo,
		StudentsRepo: studentsRepo,
		ChatsRepo:    chatsRepo,
		chatsConfig: ChatsConfig{
			userMap: map[string]int64{
				"master":  1,
				"student": 2,
			},
			userMapBackwards: map[int64]string{
				1: "master",
				2: "student",
			},
			chatTypes: map[string]int64{
				"unseen":             1,
				"approved":           2,
				"disapproved":        3,
				"deleted by master":  4,
				"deleted by student": 5,
			},
			chatTypesBackwards: map[int64]string{
				1: "unseen",
				2: "approved",
				3: "disapproved",
				4: "deleted by master",
				5: "deleted by student",
			},
			messagesTypesMap: map[int64]bool{
				1: false,
				2: true,
			},
			messagesTypesMapBackwards: map[bool]int64{
				false: 1,
				true:  2,
			},
		},
	}
	useCases.WebsocketsUC = &WebsocketsUC{
		useCases:       useCases,
		WebsocketsRepo: wsRepo,
		ChatsRepo:      chatsRepo,
		messagesTypesMap: map[int64]bool{
			1: false,
			2: true,
		},
	}
	useCases.LessonsUC = &LessonsUC{
		useCases:     useCases,
		LessonsRepo:  lessonsRepo,
		MastersRepo:  mastersRepo,
		StudentsRepo: studentsRepo,
		lessonsConfig: LessonsConfig{
			educationFormatMap:          map[int64]string{1: "онлайн", 2: "вживую"},
			educationFormatMapBackwards: map[string]int64{"онлайн": 1, "вживую": 2},
			layoutISODate:               "2006-01-02",
			layoutISOTime:               "15:04:05",
			zeroTime:                    "00:00:00",
			lessonRequestStatus:         map[int64]string{1: "approved", 2: "disapproved"},
		},
	}
	useCases.LessonsUC.lessonsConfig.zeroTimeParsed, _ = time.Parse(useCases.LessonsUC.lessonsConfig.layoutISOTime, useCases.LessonsUC.lessonsConfig.zeroTime)

	useCases.errorMessages = ErrorMessagesUC{
		DbError:             "database internal error",
		InternalServerError: "internal server error",
		FileErrors: FileErrors{
			FileOpenError:          "error opening file",
			FileReadError:          "error reading file",
			FileReadExtensionError: "error reading file extension",
			FileCreateError:        "error creating file",
			FileRemoveError:        "error removing file",
		},
	}
	useCases.errorId = 0
	useCases.filesDir = "/masterhub_files"
	//useCases.filesDir = ""

	go useCases.WebsocketsUC.Start() // GOROUTINE FOR WEBSOCKETS
	return nil
}
