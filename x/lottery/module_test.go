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

	lotteryApp *app.App
	ticketApp  *app.App

	// testing chains used for convenience and readability
	lottery *ibctesting.TestChain
	ticket  *ibctesting.TestChain
}

var apps []*app.App

func (suite *LotteryTestSuite) SetupTest() {
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)
	suite.lottery = suite.coordinator.GetChain(ibctesting.GetChainID(1))
	suite.ticket = suite.coordinator.GetChain(ibctesting.GetChainID(2))
	suite.lotteryApp = apps[0]
	suite.ticketApp = apps[1]
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
 
	creator := suite.lottery.SenderAccount.GetAddress().String()
	lotteryCreationMsg := types.MsgCreateLottery{
		Creator:  creator,
		Deadline: 50,
	}
	_, err := suite.lottery.SendMsgs(&lotteryCreationMsg)
	suite.Require().NoError(err)

	lotteries := suite.lotteryApp.LotteryKeeper.GetAllLottery(suite.lottery.GetContext())
	suite.Require().Equal(1, len(lotteries))
	suite.Require().Equal(creator, lotteries[0].Creator)

	bettingAmt := sdk.NewCoin("stake", sdk.NewInt(50))
	buyer := suite.ticket.SenderAccount.GetAddress()

	ticketContext := suite.ticket.GetContext()
	packetTimeout := ticketContext.BlockTime().Add(50 * time.Second)

	balanceBefore := suite.ticketApp.BankKeeper.GetBalance(suite.ticket.GetContext(), buyer, bettingAmt.Denom)

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
	balanceAfter := suite.ticketApp.BankKeeper.GetBalance(suite.ticket.GetContext(), buyer, "stake")

	suite.Require().True(balanceBefore.SubAmount(bettingAmt.Amount).IsEqual(balanceAfter))

}

func TestTransferTestSuite(t *testing.T) {
	suite.Run(t, new(LotteryTestSuite))
}

func SetupTestingApp() (ibctesting.TestingApp, map[string]json.RawMessage) {
	db := dbm.NewMemDB()

	encCdc := app.MakeEncodingConfig()
	application := app.New(log.NewNopLogger(), db, nil, true, map[int64]bool{}, simapp.DefaultNodeHome, 5, encCdc, simapp.EmptyAppOptions{})
	apps = append(apps, application)

	return application, app.NewDefaultGenesisState(encCdc.Marshaler)
}

func init() {
	ibctesting.DefaultTestingAppInit = SetupTestingApp
}
