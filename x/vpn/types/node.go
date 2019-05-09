package types

import (
	"sort"

	csdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"

	sdkTypes "github.com/ironman0x7b2/sentinel-sdk/types"
)

type Node struct {
	ID            sdkTypes.ID
	Owner         csdkTypes.AccAddress
	OwnerPubKey   crypto.PubKey
	DepositAmount csdkTypes.Coin

	Moniker          string
	PricesPerGB      csdkTypes.Coins
	InternetSpeed    sdkTypes.Bandwidth
	EncryptionMethod string
	Type             string
	Version          string

	ModifiedAtHeight       int64
	Status                 string
	StatusModifiedAtHeight int64
}

func (n *Node) UpdateDetails(_node Node) {
	if len(_node.Moniker) != 0 {
		n.Moniker = _node.Moniker
	}
	if _node.PricesPerGB != nil && _node.PricesPerGB.Len() > 0 &&
		_node.PricesPerGB.IsValid() && _node.PricesPerGB.IsAllPositive() {
		n.PricesPerGB = _node.PricesPerGB
	}
	if !_node.InternetSpeed.IsNil() && _node.InternetSpeed.IsPositive() {
		n.InternetSpeed = _node.InternetSpeed
	}
	if len(_node.EncryptionMethod) != 0 {
		n.EncryptionMethod = _node.EncryptionMethod
	}
	if len(_node.Type) != 0 {
		n.Type = _node.Type
	}
	if len(_node.Version) != 0 {
		n.Version = _node.Version
	}
}

func (n Node) FindPricePerGB(denom string) csdkTypes.Coin {
	index := sort.Search(n.PricesPerGB.Len(), func(i int) bool {
		return n.PricesPerGB[i].Denom >= denom
	})

	if index == n.PricesPerGB.Len() ||
		(index < n.PricesPerGB.Len() && n.PricesPerGB[index].Denom != denom) {
		return csdkTypes.Coin{}
	}

	return n.PricesPerGB[index]
}

func (n Node) AmountToBandwidth(amount csdkTypes.Coin) (sdkTypes.Bandwidth, csdkTypes.Error) {
	pricePerGB := n.FindPricePerGB(amount.Denom)
	if len(pricePerGB.Denom) == 0 || pricePerGB.Amount.IsZero() {
		return sdkTypes.Bandwidth{}, ErrorInvalidPriceDenom()
	}

	upload := amount.Amount.Quo(pricePerGB.Amount).Mul(sdkTypes.GB)
	download := amount.Amount.Quo(pricePerGB.Amount).Mul(sdkTypes.GB)
	bandwidth := sdkTypes.NewBandwidth(upload, download)

	return bandwidth, nil
}
