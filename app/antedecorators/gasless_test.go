package antedecorators_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sei-protocol/sei-chain/app/antedecorators"
	oraclekeeper "github.com/sei-protocol/sei-chain/x/oracle/keeper"
	"github.com/stretchr/testify/require"
)

var output = ""

type FakeAnteDecoratorOne struct{}

func (ad FakeAnteDecoratorOne) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	output = fmt.Sprintf("%sone", output)
	return next(ctx, tx, simulate)
}

type FakeAnteDecoratorTwo struct{}

func (ad FakeAnteDecoratorTwo) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	output = fmt.Sprintf("%stwo", output)
	return next(ctx, tx, simulate)
}

type FakeAnteDecoratorThree struct{}

func (ad FakeAnteDecoratorThree) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	output = fmt.Sprintf("%sthree", output)
	return next(ctx, tx, simulate)
}

type FakeTx struct{}

func (tx FakeTx) GetMsgs() []sdk.Msg {
	return []sdk.Msg{}
}

func (tx FakeTx) ValidateBasic() error {
	return nil
}

func TestGaslessDecorator(t *testing.T) {
	output = ""
	anteDecorators := []sdk.AnteDecorator{
		FakeAnteDecoratorOne{},
		antedecorators.NewGaslessDecorator([]sdk.AnteDecorator{FakeAnteDecoratorTwo{}}, oraclekeeper.Keeper{}),
		FakeAnteDecoratorThree{},
	}
	chainedHandler := sdk.ChainAnteDecorators(anteDecorators...)
	chainedHandler(sdk.Context{}, FakeTx{}, false)
	require.Equal(t, "onetwothree", output)
}
