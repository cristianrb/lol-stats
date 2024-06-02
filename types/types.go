package types

type LolAccount struct {
	Puuid    string `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
}

type ChampionMastery struct {
	ChampionId     int `json:"championId"`
	ChampionPoints int `json:"championPoints"`
	ChampionLevel  int `json:"championLevel"`
}
