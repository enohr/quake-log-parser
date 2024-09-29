package parser

import "regexp"

const (
	INIT_GAME_REGEX         = `InitGame:\s.*`
	KILLED_REGEX            = `Kill: (?P<killer_id>\d+) (?P<victim_id>\d+) (?P<mean_id>\d+)`
	JOIN_GAME_REGEX         = `ClientConnect: (?P<player_id>\d+)`
	USER_INFO_CHANGED_REGEX = `ClientUserinfoChanged: (?P<player_id>\d+) n\\(?P<player>.*?)\\`
	DISCONNECT_GAME_REGEX   = `ClientDisconnect: (?P<player_id>\d+)`
)

var (
	initGameRegex        = regexp.MustCompile(INIT_GAME_REGEX)
	joinGameRegex        = regexp.MustCompile(JOIN_GAME_REGEX)
	userInfoChangedRegex = regexp.MustCompile(USER_INFO_CHANGED_REGEX)
	killedRegex          = regexp.MustCompile(KILLED_REGEX)
	disconnectGameRegex  = regexp.MustCompile(DISCONNECT_GAME_REGEX)
)
