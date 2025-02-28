package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type PointsRepositoryMockTestSuite struct {
	suite.Suite
	repository PointsRepositoryMock
}

func (suite *PointsRepositoryMockTestSuite) SetupTest() {
	suite.repository = PointsRepositoryMock{}
}


func (suite *PointsRepositoryMockTestSuite) TestCharge_Amount_Properly() {
	// Set the initial amount to be 100
	suite.repository.AmountCents = 100

	chargedAmount := 20
	expected := 80
	points, err := suite.repository.Charge(chargedAmount)
	require.NoError(suite.T(), err)
	require.NotEmpty(suite.T(), points)

	assert.Equal(suite.T(), expected, points)
}

func TestPointsRepositoryMockTestSuite( t *testing.T) {
	suite.Run(t, new(PointsRepositoryMockTestSuite))
}

