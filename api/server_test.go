package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"io"
	"lol-stats/cristianrb/internal"
	mockhttpclient "lol-stats/cristianrb/mocks"
	"lol-stats/cristianrb/types"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_GetChampionMastery(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHTTPClient := mockhttpclient.NewMockHTTPClient(ctrl)
	cache := internal.NewKeyValueCache()

	lolAccount := &types.LolAccount{}
	lolAccountFilled := types.LolAccount{
		Puuid: "123",
	}

	endpoint := "https://europe.api.riotgames.com/riot/account/v1/accounts/by-riot-id/user/tag"
	mockHTTPClient.
		EXPECT().
		Do(endpoint, lolAccount).
		SetArg(1, lolAccountFilled).
		Return(nil)

	championMasteries := &[]types.ChampionMastery{}
	championMasteriesFilled := []types.ChampionMastery{
		{
			ChampionId:     1,
			ChampionLevel:  7,
			ChampionPoints: 1000,
		},
		{
			ChampionId:     2,
			ChampionLevel:  4,
			ChampionPoints: 300,
		},
	}
	masteryEndpoint := fmt.Sprintf("https://euw1.api.riotgames.com/lol/champion-mastery/v4/champion-masteries/by-puuid/%s", lolAccountFilled.Puuid)

	mockHTTPClient.
		EXPECT().
		Do(masteryEndpoint, championMasteries).
		SetArg(1, championMasteriesFilled).
		Return(nil)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/champion-mastery?region=euw1&gameName=user&tag=tag", nil)
	s := NewServer(":8080", mockHTTPClient, cache)
	s.GetChampionMastery(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	requireBodyMatchChampionMasteries(t, rr.Body, championMasteriesFilled)
}

func requireBodyMatchChampionMasteries(t *testing.T, body *bytes.Buffer, championMasteries []types.ChampionMastery) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotChampionMasteries []types.ChampionMastery
	err = json.Unmarshal(data, &gotChampionMasteries)
	require.NoError(t, err)
	require.Equal(t, championMasteries, gotChampionMasteries)
}
