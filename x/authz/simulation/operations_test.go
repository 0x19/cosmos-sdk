package simulation_test

import (
	"math/rand"
	"testing"
	"time"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/simapp"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/authz/simulation"
	"github.com/cosmos/cosmos-sdk/x/authz/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

type SimTestSuite struct {
	suite.Suite

	ctx sdk.Context
	app *simapp.SimApp
}

func (suite *SimTestSuite) SetupTest() {
	checkTx := false
	app := simapp.Setup(checkTx)
	suite.app = app
	suite.ctx = app.BaseApp.NewContext(checkTx, tmproto.Header{})
}

func (suite *SimTestSuite) TestWeightedOperations() {
	cdc := suite.app.AppCodec()
	appParams := make(simtypes.AppParams)

	weightesOps := simulation.WeightedOperations(appParams, cdc, suite.app.AccountKeeper,
		suite.app.BankKeeper, suite.app.MsgAuthKeeper, cdc,
	)

	// setup 3 accounts
	s := rand.NewSource(1)
	r := rand.New(s)
	accs := suite.getTestingAccounts(r, 3)

	expected := []struct {
		weight     int
		opMsgRoute string
		opMsgName  string
	}{
		{simappparams.DefaultWeightMsgDelegate, types.ModuleName, types.TypeMsgGrantAuthorization},
		{simappparams.DefaultWeightMsgUndelegate, types.ModuleName, types.TypeMsgRevokeAuthorization},
		{simappparams.DefaultWeightMsgSend, types.ModuleName, types.TypeMsgExecDelegated},
	}

	for i, w := range weightesOps {
		operationMsg, _, _ := w.Op()(r, suite.app.BaseApp, suite.ctx, accs, "")
		// the following checks are very much dependent from the ordering of the output given
		// by WeightedOperations. if the ordering in WeightedOperations changes some tests
		// will fail
		suite.Require().Equal(expected[i].weight, w.Weight(), "weight should be the same")
		suite.Require().Equal(expected[i].opMsgRoute, operationMsg.Route, "route should be the same")
		suite.Require().Equal(expected[i].opMsgName, operationMsg.Name, "operation Msg name should be the same")
	}
}

func (suite *SimTestSuite) getTestingAccounts(r *rand.Rand, n int) []simtypes.Account {
	accounts := simtypes.RandomAccounts(r, n)

	initAmt := sdk.TokensFromConsensusPower(200000)
	initCoins := sdk.NewCoins(sdk.NewCoin("foo", initAmt))

	// add coins to the accounts
	for _, account := range accounts {
		acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, account.Address)
		suite.app.AccountKeeper.SetAccount(suite.ctx, acc)
		err := suite.app.BankKeeper.SetBalances(suite.ctx, account.Address, initCoins)
		suite.Require().NoError(err)
	}

	return accounts
}

func (suite *SimTestSuite) TestSimulateRevokeAuthorization() {
	// setup 3 accounts
	s := rand.NewSource(1)
	r := rand.New(s)
	accounts := suite.getTestingAccounts(r, 3)

	// begin a new block
	suite.app.BeginBlock(abci.RequestBeginBlock{
		Header: tmproto.Header{
			Height:  suite.app.LastBlockHeight() + 1,
			AppHash: suite.app.LastCommitID().Hash,
		}})

	initAmt := sdk.TokensFromConsensusPower(200000)
	initCoins := sdk.NewCoins(sdk.NewCoin("foo", initAmt))

	granter := accounts[0]
	grantee := accounts[1]
	authorization := types.NewSendAuthorization(initCoins)

	suite.app.MsgAuthKeeper.Grant(suite.ctx, grantee.Address, granter.Address, authorization, time.Now().Add(30*time.Hour))

	// execute operation
	op := simulation.SimulateMsgRevokeAuthorization(suite.app.AccountKeeper, suite.app.BankKeeper, suite.app.MsgAuthKeeper)
	operationMsg, futureOperations, err := op(r, suite.app.BaseApp, suite.ctx, accounts, "")
	suite.Require().NoError(err)

	var msg types.MsgRevokeAuthorization
	types.ModuleCdc.UnmarshalJSON(operationMsg.Msg, &msg)

	suite.Require().True(operationMsg.OK)
	suite.Require().Equal(granter.Address.String(), msg.Granter)
	suite.Require().Equal(grantee.Address.String(), msg.Grantee)
	suite.Require().Equal(types.SendAuthorization{}.MethodName(), msg.AuthorizationMsgType)
	suite.Require().Len(futureOperations, 0)

}

func (suite *SimTestSuite) TestSimulateExecAuthorization() {
	// setup 3 accounts
	s := rand.NewSource(1)
	r := rand.New(s)
	accounts := suite.getTestingAccounts(r, 3)

	// begin a new block
	suite.app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: suite.app.LastBlockHeight() + 1, AppHash: suite.app.LastCommitID().Hash}})

	initAmt := sdk.TokensFromConsensusPower(200000)
	initCoins := sdk.NewCoins(sdk.NewCoin("foo", initAmt))

	granter := accounts[0]
	grantee := accounts[1]
	authorization := types.NewSendAuthorization(initCoins)

	suite.app.MsgAuthKeeper.Grant(suite.ctx, grantee.Address, granter.Address, authorization, time.Now().Add(30*time.Hour))

	// execute operation
	op := simulation.SimulateMsgExecuteAuthorized(suite.app.AccountKeeper, suite.app.BankKeeper, suite.app.MsgAuthKeeper, suite.app.AppCodec())
	operationMsg, futureOperations, err := op(r, suite.app.BaseApp, suite.ctx, accounts, "")
	suite.Require().NoError(err)

	var msg types.MsgExecAuthorized
	types.ModuleCdc.UnmarshalJSON(operationMsg.Msg, &msg)

	suite.Require().True(operationMsg.OK)
	suite.Require().Equal(grantee.Address.String(), msg.Grantee)
	suite.Require().Len(futureOperations, 0)

}

func TestSimTestSuite(t *testing.T) {
	suite.Run(t, new(SimTestSuite))
}
