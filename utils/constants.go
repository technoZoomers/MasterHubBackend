package utils

const LogFile = "log.log"
const DBName = "masterhubdb"
const PortNum = ":5000"


const ERROR_ID = 0
const DEFAULT_VIDEO_NAME = "noname"

const (
	NO_ERROR = iota
	USER_ERROR
	SERVER_ERROR
)

const VIDEO_FORMAT  = "application/octet-stream"
const FORM_DATA_VIDEO_KEY = "video"