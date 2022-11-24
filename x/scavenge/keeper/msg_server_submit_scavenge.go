package keeper

import (
	"context"

	"scavenge/x/scavenge/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/crypto"
)

// スカベンジを作成する
func (k msgServer) SubmitScavenge(goCtx context.Context, msg *types.MsgSubmitScavenge) (*types.MsgSubmitScavengeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var scavenge = types.Scavenge{
		Index:        msg.SolutionHash,
		Description:  msg.Description,
		SolutionHash: msg.SolutionHash,
		Reward:       msg.Reward,
	}

	//特定のソリューションハッシュを持つスカベンジが存在しないことを確認する
	_, isFound := k.GetScavenge(ctx, scavenge.SolutionHash)

	if isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Scavenge with that solution hash already exists")
	}

	//Scavengeモジュールアカウントのアドレスを取得する(正解者が出るまでの賞金を保持する)
	moduleAcct := sdk.AccAddress(crypto.AddressHash([]byte(types.ModuleName)))

	//スカベンジ作成者のアドレスをsdk.AccAddressに変換
	scavenger, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}

	//トークン(賞金)をsdk.Coinsに変換
	reward, err := sdk.ParseCoinsNormalized(scavenge.Reward)
	if err != nil {
		panic(err)
	}

	//スカベンジ作成者からモジュールアカウントにトークンを送信する
	sdkError := k.bankKeeper.SendCoins(ctx, scavenger, moduleAcct, reward)
	if sdkError != nil {
		return nil, sdkError
	}

	//スカベンジをストアに書き込む
	k.SetScavenge(ctx, scavenge)

	return &types.MsgSubmitScavengeResponse{}, nil
}
