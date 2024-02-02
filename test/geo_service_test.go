package test

import (
	"testing"

	"github.com/rum-people-preseed/telegram-weather-svc/internal/models"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/services"
	"github.com/rum-people-preseed/telegram-weather-svc/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestValidateCountryPass(t *testing.T) {
	mockSender, service := setUpGeoServiceTest()

	mockSender.On("SendGetRequest", mock.Anything).Return([]byte(getMockedPassResponse()), nil)

	countryName, err := service.ValidateCountry("TestCountry")

	assert.NoError(t, err)
	assert.Equal(t, countryName, "TestName")
	mockSender.AssertExpectations(t)
}

func TestValidateCountryFail(t *testing.T) {
	mockSender, service := setUpGeoServiceTest()

	mockSender.On("SendGetRequest", mock.Anything).Return([]byte("[]"), nil)

	countryName, err := service.ValidateCountry("TestCountry")

	assert.Error(t, err)
	assert.Equal(t, countryName, "")
	mockSender.AssertExpectations(t)
}

func TestValidateCountryTotalResultsPass(t *testing.T) {
	mockSender, service := setUpGeoServiceTest()

	mockSender.On("SendGetRequest", mock.Anything).Return([]byte(getMockedPassResponse()), nil)

	countryName, err := service.ValidateCountry("TestCountry")

	assert.NoError(t, err)
	assert.Equal(t, countryName, "TestName")
	mockSender.AssertExpectations(t)
}

func TestValidateCountryTotalResultsFail(t *testing.T) {
	mockSender, service := setUpGeoServiceTest()

	mockSender.On("SendGetRequest", mock.Anything).Return([]byte(getMockedFailResponse()), nil)

	countryName, err := service.ValidateCountry("TestCountry")

	assert.Error(t, err)
	assert.Equal(t, countryName, "")
	mockSender.AssertExpectations(t)
}

func getMockedFailResponse() string {
	return `{
    "totalResultsCount": 0,
    "geonames": [{
            "adminCode1": "08",
            "lng": "32.61458",
            "geonameId": 706448,
            "toponymName": "TestName",
            "countryId": "690791",
            "fcl": "P",
            "population": 283649,
            "countryCode": "UA",
            "name": "TestName",
            "fclName": "city, village,...",
            "adminCodes1": {
                "ISO3166_2": "65"
            },
            "countryName": "TestName",
            "fcodeName": "seat of a first-order administrative division",
            "adminName1": "TestName",
            "lat": "46.63695",
            "fcode": "PPLA"
        }
    ]
}`
}
func getMockedPassResponse() string {
	return `{
    "totalResultsCount": 1,
    "geonames": [{
            "adminCode1": "08",
            "lng": "32.61458",
            "geonameId": 706448,
            "toponymName": "TestName",
            "countryId": "690791",
            "fcl": "P",
            "population": 283649,
            "countryCode": "UA",
            "name": "TestName",
            "fclName": "city, village,...",
            "adminCodes1": {
                "ISO3166_2": "65"
            },
            "countryName": "TestName",
            "fcodeName": "seat of a first-order administrative division",
            "adminName1": "TestName",
            "lat": "46.63695",
            "fcode": "PPLA"
        }
    ]
}`
}

func setUpGeoServiceTest() (*mocks.MockSender, services.GeoNameService) {
	mockSender := new(mocks.MockSender)
	log := models.GetNewLogger()
	service := services.NewGeoNameService(mockSender, log, 0.6)
	return mockSender, service
}
