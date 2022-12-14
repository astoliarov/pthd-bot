package services

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"pthd-bot/tests/mocks"
	"testing"
)

type TeamKillServiceTestCase struct {
	suite.Suite

	controller *gomock.Controller
	daoMock    *mocks.MockITeamKillLogDAO

	service *TeamKillService
}

func (suite *TeamKillServiceTestCase) SetupTest() {
	suite.controller = gomock.NewController(suite.T())

	suite.daoMock = mocks.NewMockITeamKillLogDAO(suite.controller)

	responseService := &ResponseSelectorService{}

	suite.service = NewTeamKillService(
		suite.daoMock,
		responseService,
	)
}

func (suite *TeamKillServiceTestCase) Test_AddTeamKill_PassedRequest_SavedTeamKill() {
	request := &TeamKillRequest{
		Killer: "Roma",
		Victim: "Igor",
	}
	source := "test"

	suite.daoMock.EXPECT().Save(gomock.Any()).Return(nil)

	response, err := suite.service.AddTeamKill(request, source)

	assert.Nil(suite.T(), err)
	assert.NotEqual(suite.T(), "", response)
}

func (suite *TeamKillServiceTestCase) Test_AddTeamKill_SaveReturnedError_ProcessReturnedError() {
	request := &TeamKillRequest{
		Killer: "Roma",
		Victim: "Igor",
	}

	saveErr := errors.New("Failed to save request")

	suite.daoMock.EXPECT().Save(gomock.Any()).Return(saveErr)

	source := "test"

	response, err := suite.service.AddTeamKill(request, source)

	assert.Equal(suite.T(), err, saveErr)
	assert.Equal(suite.T(), response, "")
}

func (suite *TeamKillServiceTestCase) Test_validateName_NoMixOfAlphabets_NoErr() {

	err := suite.service.validateName("рома")

	assert.Nil(suite.T(), err)
}

func (suite *TeamKillServiceTestCase) Test_validateName_MixOfAlphabets_Err() {

	err := suite.service.validateName("рoма")

	assert.NotNil(suite.T(), err)
}

func TestTeamKillService(t *testing.T) {
	testSuite := TeamKillServiceTestCase{}
	suite.Run(t, &testSuite)
}
