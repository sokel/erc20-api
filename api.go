package api

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/common"
)

const defaultGasPrice = 20 * 1000000000

type Tokener interface {
	Approve(key *ecdsa.PrivateKey, to string, amount *big.Int) (*types.Transaction, error)
	Transfer(key *ecdsa.PrivateKey, to string, amount *big.Int) (*types.Transaction, error)
	TransferFrom(key *ecdsa.PrivateKey, from string, to string, amount *big.Int) (*types.Transaction, error)
	BalanceOf(address string) (*big.Int, error)
	AllowanceOf(from string, to string) (*big.Int, error)
	TotalSupply() (*big.Int, error)
}

type API struct {
	client   *ethclient.Client
	gasPrice int

	tokenContract *StandardToken
}

func NewAPI(ethEndpoint string, gasPrice *int, tokenContractAddress string) (*API, error) {
	client, err := ethclient.Dial(ethEndpoint)
	if err != nil {
		return nil, err
	}

	if gasPrice == nil {
		*gasPrice = defaultGasPrice
	}

	tokenContract, err := NewStandardToken(common.HexToAddress(tokenContractAddress), client)
	if err != nil {
		return nil, err
	}

	bch := &API{
		client:        client,
		gasPrice:      *gasPrice,
		tokenContract: tokenContract,
	}
	return bch, nil
}

func (bch *API) BalanceOf(address string) (*big.Int, error) {
	balance, err := bch.tokenContract.BalanceOf(&bind.CallOpts{Pending: true}, common.HexToAddress(address))
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (bch *API) AllowanceOf(from string, to string) (*big.Int, error) {
	allowance, err := bch.tokenContract.Allowance(&bind.CallOpts{Pending: true}, common.HexToAddress(from), common.HexToAddress(to))
	if err != nil {
		return nil, err
	}
	return allowance, nil
}

func (bch *API) Approve(key *ecdsa.PrivateKey, to string, amount *big.Int) (*types.Transaction, error) {
	opts := getTxOpts(key, 50000)

	tx, err := bch.tokenContract.Approve(opts, common.HexToAddress(to), amount)
	if err != nil {
		return nil, err
	}
	return tx, err
}

func (bch *API) Transfer(key *ecdsa.PrivateKey, to string, amount *big.Int) (*types.Transaction, error) {
	opts := getTxOpts(key, 50000)

	tx, err := bch.tokenContract.Transfer(opts, common.HexToAddress(to), amount)
	if err != nil {
		return nil, err
	}
	return tx, err
}

func (bch *API) TransferFrom(key *ecdsa.PrivateKey, from string, to string, amount *big.Int) (*types.Transaction, error) {
	opts := getTxOpts(key, 50000)

	tx, err := bch.tokenContract.TransferFrom(opts, common.HexToAddress(from), common.HexToAddress(to), amount)
	if err != nil {
		return nil, err
	}
	return tx, err
}

func (bch *API) TotalSupply() (*big.Int, error) {
	totalSupply, err := bch.tokenContract.TotalSupply(&bind.CallOpts{Pending: true})
	if err != nil {
		return nil, err
	}
	return totalSupply, nil
}

func getTxOpts(key *ecdsa.PrivateKey, gasLimit int64) (*bind.TransactOpts) {
	opts := bind.NewKeyedTransactor(key)
	opts.GasLimit = big.NewInt(gasLimit)
	return opts
}
