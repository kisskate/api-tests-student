package favorites

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	favoritesManager "api-tests-template/internal/managers/favorites"
	favoriteModels "api-tests-template/internal/managers/favorites/models"
	"api-tests-template/internal/utils"

	testdata "api-tests-template/tests"
)

const nonExistentAdvertisementID = "ffffffff-ffff-4fff-bfff-ffffffffffff"

type TestSuite struct {
	suite.Suite
	ownerLogin  authModels.LoginOkResponse
	secondLogin authModels.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	testdata.SetupSuite()
	testdata.Precondition("Авторизуем два аккаунта для сценариев избранного")
	require.NotEmpty(s.T(), os.Getenv("TEST_LOGIN"), "задайте TEST_LOGIN в .env.override")
	require.NotEmpty(s.T(), os.Getenv("TEST_PASSWORD"), "задайте TEST_PASSWORD в .env.override")
	require.NotEmpty(s.T(), os.Getenv("TEST_LOGIN2"), "задайте TEST_LOGIN2 в .env.override")
	require.NotEmpty(s.T(), os.Getenv("TEST_PASSWORD2"), "задайте TEST_PASSWORD2 в .env.override")

	s.ownerLogin = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
	s.secondLogin = auth.Login(s.T(), os.Getenv("TEST_LOGIN2"), os.Getenv("TEST_PASSWORD2"))
	require.NotEqual(s.T(), s.ownerLogin.User.Id, s.secondLogin.User.Id, "для сценария нужны два разных аккаунта")
}

func (s *TestSuite) TearDownSuite() {
	testdata.TearDownSuite()
}

func (s *TestSuite) TestAddAndRemoveFavorite() {
	t := s.T()
	advertisement := availableForeignAdvertisement(t, s.ownerLogin.Token, s.ownerLogin.User.Id)
	testdata.Precondition("Добавляем в избранное чужое объявление, отсутствующее в списке")

	added := favoritesManager.AddFavorite(t, s.ownerLogin.Token, advertisement.ID)
	require.Equal(t, http.StatusCreated, added.StatusCode, "не удалось добавить объявление в избранное: %s", added.Error.Message)
	markFavoriteCleaned := registerFavoriteCleanup(t, s.ownerLogin.Token, advertisement.ID)
	require.Equal(t, advertisement.ID, added.Favorite.Advertisement.ID)

	t.Run("объявление появляется в GET /favorites", func(t *testing.T) {
		utils.LogWithLabelAndTimeStamp("Step", "Проверяем появление объявления в избранном")
		list := favoritesManager.GetFavorites(t, s.ownerLogin.Token)
		require.Equal(t, http.StatusOK, list.StatusCode, "не удалось получить избранное: %s", list.Error.Message)
		assertions.ContainsFavoriteByAdvertisementID(t, list.List.Items, advertisement.ID)
	})

	t.Run("объявление удаляется из GET /favorites", func(t *testing.T) {
		utils.LogWithLabelAndTimeStamp("Step", "Удаляем объявление из избранного")
		removed := favoritesManager.DeleteFavorite(t, s.ownerLogin.Token, advertisement.ID)
		require.Equal(t, http.StatusNoContent, removed.StatusCode, "не удалось удалить из избранного: %s", removed.Error.Message)
		markFavoriteCleaned()

		list := favoritesManager.GetFavorites(t, s.ownerLogin.Token)
		require.Equal(t, http.StatusOK, list.StatusCode, "не удалось повторно получить избранное")
		assertions.DoesNotContainFavoriteByAdvertisementID(t, list.List.Items, advertisement.ID)
	})
}

func (s *TestSuite) TestAddDuplicateFavorite() {
	t := s.T()
	advertisement := availableForeignAdvertisement(t, s.ownerLogin.Token, s.ownerLogin.User.Id)
	testdata.Precondition("Проверяем конфликт при повторном добавлении одного объявления")

	firstAdd := favoritesManager.AddFavorite(t, s.ownerLogin.Token, advertisement.ID)
	require.Equal(t, http.StatusCreated, firstAdd.StatusCode, "не удалось подготовить избранное: %s", firstAdd.Error.Message)
	registerFavoriteCleanup(t, s.ownerLogin.Token, advertisement.ID)

	duplicate := favoritesManager.AddFavorite(t, s.ownerLogin.Token, advertisement.ID)
	assertions.APIError(t, http.StatusConflict, duplicate.StatusCode, "already_exists", duplicate.Error.Error, "Advertisement is already in favorites", duplicate.Error.Message)
}

func (s *TestSuite) TestAddOwnAdvertisementToFavorites() {
	t := s.T()
	advertisement, _ := testdata.CreateTestAdvertisement(t, s.ownerLogin.Token)
	testdata.Precondition("Проверяем запрет добавления собственного объявления в избранное")

	result := favoritesManager.AddFavorite(t, s.ownerLogin.Token, advertisement.ID)
	assertions.APIError(t, http.StatusForbidden, result.StatusCode, "forbidden", result.Error.Error, "Cannot add own advertisement to favorites", result.Error.Message)
}

func (s *TestSuite) TestFavoritesWithNonExistentAdvertisement() {
	t := s.T()
	testdata.Precondition("Проверяем ошибки для несуществующего объявления")

	t.Run("добавление в избранное возвращает 404", func(t *testing.T) {
		utils.LogWithLabelAndTimeStamp("Step", "Добавляем несуществующее объявление в избранное")
		result := favoritesManager.AddFavorite(t, s.ownerLogin.Token, nonExistentAdvertisementID)
		assertions.APIError(t, http.StatusNotFound, result.StatusCode, "not_found", result.Error.Error, "Advertisement not found", result.Error.Message)
	})

	t.Run("удаление из избранного возвращает 404", func(t *testing.T) {
		utils.LogWithLabelAndTimeStamp("Step", "Удаляем несуществующее объявление из избранного")
		result := favoritesManager.DeleteFavorite(t, s.ownerLogin.Token, nonExistentAdvertisementID)
		assertions.APIError(t, http.StatusNotFound, result.StatusCode, "not_found", result.Error.Error, "Favorite entry not found", result.Error.Message)
	})
}

func (s *TestSuite) TestAddFavoriteWithoutAdvertisementID() {
	testdata.Precondition("Проверяем валидацию обязательного advertisement_id")
	result := favoritesManager.AddFavorite(s.T(), s.ownerLogin.Token, "")
	assertions.APIError(
		s.T(),
		http.StatusBadRequest,
		result.StatusCode,
		"invalid_request",
		result.Error.Error,
		"invalid UUID length: 0",
		result.Error.Message,
	)
}

func (s *TestSuite) TestFavoritesWithoutAuthorization() {
	t := s.T()
	advertisement := availableForeignAdvertisement(t, s.ownerLogin.Token, s.ownerLogin.User.Id)
	testdata.Precondition("Проверяем, что все ручки избранного требуют авторизацию")

	t.Run("GET /favorites", func(t *testing.T) {
		utils.LogWithLabelAndTimeStamp("Step", "Запрашиваем избранное без токена")
		result := favoritesManager.GetFavorites(t, "")
		assertions.APIError(t, http.StatusUnauthorized, result.StatusCode, "unauthorized", result.Error.Error, "Authorization header required", result.Error.Message)
	})

	t.Run("POST /favorites", func(t *testing.T) {
		utils.LogWithLabelAndTimeStamp("Step", "Добавляем объявление в избранное без токена")
		result := favoritesManager.AddFavorite(t, "", advertisement.ID)
		assertions.APIError(t, http.StatusUnauthorized, result.StatusCode, "unauthorized", result.Error.Error, "Authorization header required", result.Error.Message)
	})

	t.Run("DELETE /favorites/{ad_id}", func(t *testing.T) {
		utils.LogWithLabelAndTimeStamp("Step", "Удаляем объявление из избранного без токена")
		result := favoritesManager.DeleteFavorite(t, "", advertisement.ID)
		assertions.APIError(t, http.StatusUnauthorized, result.StatusCode, "unauthorized", result.Error.Error, "Authorization header required", result.Error.Message)
	})
}

func (s *TestSuite) TestDeletingAdvertisementRemovesItFromSecondUsersFavorites() {
	t := s.T()
	testdata.Precondition("Создаем объявление владельца и добавляем его в избранное второго аккаунта")
	advertisement, markAdvertisementDeleted := testdata.CreateTestAdvertisement(t, s.ownerLogin.Token)

	added := favoritesManager.AddFavorite(t, s.secondLogin.Token, advertisement.ID)
	require.Equal(t, http.StatusCreated, added.StatusCode, "второй пользователь не смог добавить объявление в избранное: %s", added.Error.Message)
	markFavoriteCleaned := registerFavoriteCleanup(t, s.secondLogin.Token, advertisement.ID)

	utils.LogWithLabelAndTimeStamp("Step", "Удаляем объявление токеном владельца")
	deleted := advertisements.DeleteAdvertisement(t, s.ownerLogin.Token, advertisement.ID)
	require.Equal(t, http.StatusNoContent, deleted.StatusCode, "владелец не смог удалить объявление: %s", deleted.Error.Message)
	markAdvertisementDeleted()

	utils.LogWithLabelAndTimeStamp("Step", "Проверяем, что удаленное объявление исчезло из избранного второго пользователя")
	list := favoritesManager.GetFavorites(t, s.secondLogin.Token)
	require.Equal(t, http.StatusOK, list.StatusCode, "не удалось получить избранное второго пользователя: %s", list.Error.Message)
	assertions.DoesNotContainFavoriteByAdvertisementID(t, list.List.Items, advertisement.ID)
	markFavoriteCleaned()
}

func availableForeignAdvertisement(t *testing.T, token string, userID string) models.Advertisement {
	t.Helper()
	list := advertisements.GetAdvertisementsList(t, token, "")
	favorites := favoritesManager.GetFavorites(t, token)
	require.Equal(t, http.StatusOK, favorites.StatusCode, "не удалось получить избранное: %s", favorites.Error.Message)

	for _, item := range list.Items {
		if item.ID == "" || item.UserID == userID || favoriteModelsContain(favorites.List.Items, item.ID) {
			continue
		}
		return item
	}
	t.Fatalf("в каталоге нет чужого объявления, доступного для добавления в избранное")
	return models.Advertisement{}
}

func favoriteModelsContain(items []favoriteModels.Favorite, advertisementID string) bool {
	for _, item := range items {
		if item.Advertisement.ID == advertisementID {
			return true
		}
	}
	return false
}

func registerFavoriteCleanup(t *testing.T, token string, advertisementID string) func() {
	t.Helper()
	removed := false
	t.Cleanup(func() {
		if removed {
			return
		}
		testdata.Precondition("Очищаем временное избранное после теста")
		result := favoritesManager.DeleteFavorite(t, token, advertisementID)
		if result.StatusCode != http.StatusNoContent && result.StatusCode != http.StatusNotFound {
			t.Errorf("не удалось очистить избранное: status=%d, error=%s", result.StatusCode, result.Error.Message)
		}
	})
	return func() { removed = true }
}
