package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/crypto"

	"scavenge/x/scavenge/types"
)

// 更新されたスカベンジをストアに書き込む(回答する)
func (k msgServer) RevealSolution(goCtx context.Context, msg *types.MsgRevealSolution) (*types.MsgRevealSolutionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	//ソリューションとスカベンジャーアドレスを連結してバイトに変換する
	var solutionScavengerBytes = []byte(msg.Solution + msg.Creator)
	var solutionScavengerHash = sha256.Sum256(solutionScavengerBytes)
	var solutionScavengerHashString = hex.EncodeToString(solutionScavengerHash[:])

	//特定のソリューションスカベンジャーハッシュを持つコミットがストアに存在することを確認
	commit, isFound := k.GetCommit(ctx, solutionScavengerHashString)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Commit with that hash already exists")
	}

	//指定されたソリューション ハッシュを持つスカベンジがストアに存在することを確認
	scavenge, isFound := k.GetScavenge(ctx, commit.SolutionHash)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Scavenge with that solution hash already exists")

	}
	//スカベンジがまだ解決されていないことを確認する
	_, err := sdk.AccAddressFromBech32(scavenge.Scavenger)

	//スカベンジがすでに解決されている場合はエラー
	if err == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Scavenge has already been solved")
	}

	//正解者,回答をスカベンジにセットする
	scavenge.Scavenger = msg.Creator
	scavenge.Solution = msg.Solution

	//モジュールアカウントのアドレスを取得する(正解者が出るまでの賞金を保持していたので)
	moduleAcct := sdk.AccAddress(crypto.AddressHash([]byte(types.ModuleName)))

	//正解者のアドレスをsdk.AccAddressに変換
	scavenger, err := sdk.AccAddressFromBech32(scavenge.Scavenger)
	if err != nil {
		panic(err)
	}

	//トークン(賞金)をsdk.Coinsに変換
	reward, err := sdk.ParseCoinsNormalized(scavenge.Reward)
	if err != nil {
		panic(err)
	}

	// モジュールアカウントから正解を出したアカウントにトークンを送る
	sdkError := k.bankKeeper.SendCoins(ctx, moduleAcct, scavenger, reward)
	if sdkError != nil {
		return nil, sdkError
	}

	//スカベンジをストアに書き込む(Scavenger,Solutioを更新)
	k.SetScavenge(ctx, scavenge)

	return &types.MsgRevealSolutionResponse{}, nil
}
