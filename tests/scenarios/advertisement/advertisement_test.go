package advertisement

import (
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	advertisementModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"

	testdata "api-tests-template/tests"
)

type TestSuite struct {
	suite.Suite
	loginData authModels.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	testdata.SetupSuite()
	testdata.Precondition("Авторизация пользователя для сценариев объявления")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

func (s *TestSuite) TearDownSuite() {
	testdata.TearDownSuite()
}

func (s *TestSuite) TestCreateAdvertisementPositive() {
	t := s.T()
	testdata.Precondition("Создаем объявление со всеми полями и фотографией")
	created, _ := testdata.CreateTestAdvertisement(t, s.loginData.Token)

	t.Run("GET /advertisement/{id} возвращает параметры создания", func(t *testing.T) {
		utils.LogWithLabelAndTimeStamp("Step", "Получаем созданное объявление по ID")
		result := advertisements.GetAdvertisement(t, s.loginData.Token, created.ID)
		require.Equal(t, http.StatusOK, result.StatusCode, "не удалось получить созданное объявление: %s", result.Error.Message)
		require.Equal(t, created.ID, result.Advertisement.ID)
		require.Equal(t, created.Title, result.Advertisement.Title)
		require.Equal(t, created.Description, result.Advertisement.Description)
		require.Equal(t, created.Price, result.Advertisement.Price)
		require.Equal(t, created.Quantity, result.Advertisement.Quantity)
	})

	t.Run("GET /advertisements?search={title} находит объявление по полному названию", func(t *testing.T) {
		utils.LogWithLabelAndTimeStamp("Step", "Ищем созданное объявление по полному названию")
		list := advertisements.GetAdvertisementsList(t, s.loginData.Token, created.Title)
		assertions.ContainsAdvertisementByID(t, list.Items, created.ID)
	})

	t.Run("GET /advertisements/{id}/photos возвращает загруженную фотографию", func(t *testing.T) {
		utils.LogWithLabelAndTimeStamp("Step", "Получаем загруженные фотографии объявления")
		result := advertisements.GetAdvertisementPhotos(t, s.loginData.Token, created.ID)
		require.Equal(t, http.StatusOK, result.StatusCode, "не удалось получить фотографии: %s", result.Error.Message)
		require.NotEmpty(t, result.Photos, "API не вернул загруженные фотографии")
		require.NotEmpty(t, result.Photos[0].URL, "у фотографии отсутствует URL")
		if len(created.Photos) > 0 {
			require.Equal(t, created.Photos[0].ID, result.Photos[0].ID, "вернулась другая фотография")
		}
	})
}

func (s *TestSuite) TestCreateAdvertisementInvalidFields() {
	cases := []struct {
		name            string
		expectedMessage string
		change          func(*advertisementModels.CreateAdvertisementRequest)
	}{
		{name: "пустой title", expectedMessage: "title is required", change: func(request *advertisementModels.CreateAdvertisementRequest) { request.Title = "" }},
		{name: "отрицательная цена", expectedMessage: "price must be between 0 and 1000000 inclusive", change: func(request *advertisementModels.CreateAdvertisementRequest) { request.Price = -1 }},
		{name: "слишком длинный title", expectedMessage: "title must not exceed 50 characters provided 2048", change: func(request *advertisementModels.CreateAdvertisementRequest) {
			request.Title = strings.Repeat("x", 2048)
		}},
		{name: "пустое описание", expectedMessage: "description is required", change: func(request *advertisementModels.CreateAdvertisementRequest) { request.Description = "" }},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			t := s.T()
			utils.LogWithLabelAndTimeStamp("Step", "Проверяем отклонение объявления: "+tc.name)
			request := testdata.NewTestAdvertisementRequest()
			tc.change(&request)
			result := advertisements.CreateAdvertisement(t, s.loginData.Token, request, []string{testdata.CreateTestPhoto(t)})
			assertions.APIError(t, http.StatusBadRequest, result.StatusCode, "invalid_request", result.Error.Error, tc.expectedMessage, result.Error.Message)
		})
	}
}

func (s *TestSuite) TestCreateAdvertisementWithoutAuthorization() {
	utils.LogWithLabelAndTimeStamp("Step", "Пытаемся создать объявление без токена")
	result := advertisements.CreateAdvertisement(
		s.T(),
		"",
		testdata.NewTestAdvertisementRequest(),
		[]string{testdata.CreateTestPhoto(s.T())},
	)
	assertions.APIError(s.T(), http.StatusUnauthorized, result.StatusCode, "unauthorized", result.Error.Error, "Authorization header required", result.Error.Message)
}
