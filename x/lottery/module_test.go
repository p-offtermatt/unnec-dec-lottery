package lottery_test

import (
	"encoding/json"
	"lottery/app"
	"lottery/x/lottery/types"
	"testing"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/simapp"
	ibctesting "github.com/cosmos/ibc-go/v6/testing"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/libs/log"
)

type LotteryTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	app ibctesting.TestingApp

	// testing chains used for convenience and readability
	lottery *ibctesting.TestChain
	ticket  *ibctesting.TestChain
}

func (suite *LotteryTestSuite) SetupTest() {
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)
	suite.lottery = suite.coordinator.GetChain(ibctesting.GetChainID(1))
	suite.ticket = suite.coordinator.GetChain(ibctesting.GetChainID(2))
}

func NewTransferPath(lottery, ticket *ibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(lottery, ticket)
	path.EndpointA.ChannelConfig.PortID = types.PortID
	path.EndpointB.ChannelConfig.PortID = types.PortID
	path.EndpointA.ChannelConfig.Version = types.Version
	path.EndpointB.ChannelConfig.Version = types.Version

	return path
}

func (suite *LotteryTestSuite) TestLotteryLifecycleCancel() {
	path := NewTransferPath(suite.lottery, suite.ticket)
	suite.coordinator.Setup(path)

	lotteryCreationMsg := types.MsgCreateLottery{
		Creator:  suite.lottery.SenderAccount.GetAddress().String(),
		Deadline: 50,
	}
	_, err := suite.lottery.SendMsgs(&lotteryCreationMsg)
	suite.Require().NoError(err)
}

func TestTransferTestSuite(t *testing.T) {
	suite.Run(t, new(LotteryTestSuite))
}

var DefaultTestingAppInit = SetupTestingApp

func SetupTestingApp() (ibctesting.TestingApp, map[string]json.RawMessage) {
	db := dbm.NewMemDB()

	encCdc := simapp.MakeTestEncodingConfig()
	app := app.New(log.NewNopLogger(), db, nil, true, map[int64]bool{}, simapp.DefaultNodeHome, 5, encCdc, simapp.EmptyAppOptions{})
	return app, simapp.NewDefaultGenesisState(encCdc.Codec)
}
