package lottery_test

import (
	"encoding/json"
	"lottery/app"
	"lottery/x/lottery/types"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctesting "github.com/cosmos/ibc-go/v6/testing"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

type LotteryTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	app app.App

	// testing chains used for convenience and readability
	lottery *ibctesting.TestChain
	ticket  *ibctesting.TestChain
}

func (suite *LotteryTestSuite) SetupTest() {
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)
	suite.lottery = suite.coordinator.GetChain(ibctesting.GetChainID(1))
	suite.ticket = suite.coordinator.GetChain(ibctesting.GetChainID(2))
}

func NewLotteryPath(lottery, ticket *ibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(lottery, ticket)
	path.EndpointA.ChannelConfig.PortID = types.PortID
	path.EndpointB.ChannelConfig.PortID = types.PortID
	path.EndpointA.ChannelConfig.Version = types.Version
	path.EndpointB.ChannelConfig.Version = types.Version

	return path
}

func (suite *LotteryTestSuite) TestLotteryLifecycleCancel() {
	path := NewLotteryPath(suite.ticket, suite.lottery)
	suite.coordinator.Setup(path)

	lotteryCreationMsg := types.MsgCreateLottery{
		Creator:  suite.lottery.SenderAccount.GetAddress().String(),
		Deadline: 50,
	}
	_, err := suite.lottery.SendMsgs(&lotteryCreationMsg)
	suite.Require().NoError(err)

	bettingAmt := sdk.NewCoin("stake", sdk.NewInt(50))
	buyer := suite.ticket.SenderAccount.GetAddress()

	ticketContext := suite.ticket.GetContext()
	packetTimeout := ticketContext.BlockTime().Add(50 * time.Second)

	// balanceBefore := suite.ticket.GetSimApp().BankKeeper.GetBalance(suite.ticket.GetContext(), buyer, bettingAmt.Denom)

	buyTicketMsg := types.MsgSendBuyTicket{
		LotteryId:        0,
		Price:            bettingAmt,
		Creator:          buyer.String(),
		Port:             path.EndpointB.ChannelConfig.PortID,
		ChannelID:        path.EndpointB.ChannelID,
		TimeoutTimestamp: uint64(packetTimeout.Unix()),
	}

	res, err := suite.ticket.SendMsgs(&buyTicketMsg)
	suite.Require().NoError(err)

	packet, err := ibctesting.ParsePacketFromEvents(res.GetEvents())
	suite.Require().NoError(err)

	// relay send
	err = path.RelayPacket(packet)
	suite.Require().NoError(err) // relay committed

	// check that on chain ticket coins were spent
	// balanceAfter := suite.ticket.GetSimApp().BankKeeper.GetBalance(suite.ticket.GetContext(), buyer, "stake")

	// suite.Require().Equal(balanceBefore.SubAmount(bettingAmt.Amount), balanceAfter)
}

func TestTransferTestSuite(t *testing.T) {
	suite.Run(t, new(LotteryTestSuite))
}

func SetupTestingApp() (ibctesting.TestingApp, map[string]json.RawMessage) {
	db := dbm.NewMemDB()

	encCdc := app.MakeEncodingConfig()
	application := app.New(log.NewNopLogger(), db, nil, true, map[int64]bool{}, simapp.DefaultNodeHome, 5, encCdc, simapp.EmptyAppOptions{})
	return application, app.NewDefaultGenesisState(encCdc.Marshaler)
}

func init() {
	ibctesting.DefaultTestingAppInit = SetupTestingApp
}
