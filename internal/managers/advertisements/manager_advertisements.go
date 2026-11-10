package advertisements

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
)

func GetAdvertisementsList(t *testing.T, token string, textForSearch string) models.AdvertisementsList {
	response := advertisements.HttpGetAdvertisements(t, token, textForSearch)
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	var advertisementsList models.AdvertisementsList
	err = json.Unmarshal(body, &advertisementsList)
	require.NoError(t, err)

	return advertisementsList
}
