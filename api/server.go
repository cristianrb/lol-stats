package api

import (
	"fmt"
	"log"
	"lol-stats/cristianrb/internal"
	"lol-stats/cristianrb/types"
	"net/http"
)

type Server struct {
	addr   string
	client HTTPClient
	cache  internal.Cache
}

var platformToRegion = map[string]string{
	"euw1": "europe",
	"eun1": "europe",
	"tr1":  "europe",
	"ru":   "europe",
	"na1":  "americas",
	"br1":  "americas",
	"la1":  "americas",
	"la2":  "americas",
	"jp1":  "asia",
	"kr":   "asia",
	"oc1":  "sea",
	"ph2":  "sea",
	"sg2":  "sea",
	"th2":  "sea",
	"tw2":  "sea",
	"vn2":  "sea",
}

func NewServer(addr string, client HTTPClient, cache internal.Cache) *Server {
	return &Server{
		addr:   addr,
		client: client,
		cache:  cache,
	}
}

func (s *Server) Run() error {
	srv := &http.Server{
		Addr:    s.addr,
		Handler: s.routes(),
	}

	log.Printf("Started lol-stats at %s\n", srv.Addr)
	return srv.ListenAndServe()
}

func (s *Server) GetChampionMastery(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	platform := s.readString(qs, "region", "")
	gameName := s.readString(qs, "gameName", "")
	tagLine := s.readString(qs, "tag", "")
	region := platformToRegion[platform]

	cacheKey := fmt.Sprintf("%s#%s#%s", region, gameName, tagLine)
	puuid, err := s.cache.Get(cacheKey)
	if err != nil {
		puuid, err = s.getAccount(region, gameName, tagLine)
		if err != nil {
			return
		}

		err = s.cache.Save(cacheKey, puuid)
		if err != nil {
			return
		}
	}

	baseURL := fmt.Sprintf("https://%s.api.riotgames.com", platform)
	endpoint := fmt.Sprintf("%s/lol/champion-mastery/v4/champion-masteries/by-puuid/%s", baseURL, puuid)
	championMasteries := &[]types.ChampionMastery{}
	err = s.client.Do(endpoint, championMasteries)
	if err != nil {
		return
	}

	s.writeJSON(w, http.StatusOK, championMasteries, nil)
}

func (s *Server) getAccount(region, gameName, tagLine string) (string, error) {
	baseURL := fmt.Sprintf("https://%s.api.riotgames.com", region)
	endpoint := fmt.Sprintf("%s/riot/account/v1/accounts/by-riot-id/%s/%s", baseURL, gameName, tagLine)
	lolAccount := &types.LolAccount{}
	err := s.client.Do(endpoint, lolAccount)
	if err != nil {
		return "", err
	}

	return lolAccount.Puuid, nil
}
