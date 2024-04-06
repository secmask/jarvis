package common

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tranvictor/jarvis/config"
	"github.com/tranvictor/jarvis/networks"
	"math/big"
)

// RawTxToHash returns valid hex data of a transaction to
// transaction hash
func RawTxToHash(data string) string {
	return crypto.Keccak256Hash(hexutil.MustDecode(data)).Hex()
}

func BuildExactTx(nonce uint64, to string, ethAmount *big.Int, gasLimit uint64, priceGwei float64, data []byte) (tx *types.Transaction) {
	toAddress := common.HexToAddress(to)
	gasPrice := GweiToWei(priceGwei)
	tipInt := GweiToWei(config.TipGas)
	chainIDInt := big.NewInt(networks.CurrentNetwork().GetChainID())
	if config.TxType == config.TxTypeDynamicFee {
		return types.NewTx(&types.DynamicFeeTx{
			ChainID:   chainIDInt,
			Nonce:     nonce,
			GasTipCap: tipInt,
			GasFeeCap: gasPrice,
			Gas:       gasLimit,
			To:        &toAddress,
			Value:     ethAmount,
			Data:      data,
		})
	} else {
		return types.NewTx(&types.LegacyTx{
			Nonce:    nonce,
			GasPrice: gasPrice,
			Gas:      gasLimit,
			To:       &toAddress,
			Value:    ethAmount,
			Data:     data,
		})
	}
}

func BuildTx(nonce uint64, to string, ethAmount float64, gasLimit uint64, priceGwei float64, data []byte) (tx *types.Transaction) {
	amount := FloatToBigInt(ethAmount, 18)
	return BuildExactTx(nonce, to, amount, gasLimit, priceGwei, data)
}

func BuildExactSendETHTx(nonce uint64, to string, ethAmount *big.Int, gasLimit uint64, priceGwei float64,
	tipGas float64, txType string, chainID int64,
) (tx *types.Transaction) {
	return BuildExactTx(nonce, to, ethAmount, gasLimit, priceGwei, []byte{})
}

func BuildContractCreationTx(nonce uint64, ethAmount *big.Int, gasLimit uint64, priceGwei float64, data []byte) (tx *types.Transaction) {
	gasPrice := GweiToWei(priceGwei)
	tipInt := GweiToWei(config.TipGas)
	chainIDInt := big.NewInt(networks.CurrentNetwork().GetChainID())
	if config.TxType == config.TxTypeDynamicFee {
		return types.NewTx(&types.DynamicFeeTx{
			ChainID:   chainIDInt,
			Nonce:     nonce,
			GasTipCap: tipInt,
			GasFeeCap: gasPrice,
			Gas:       gasLimit,
			Value:     ethAmount,
			Data:      data,
		})
	} else {
		return types.NewTx(&types.LegacyTx{
			Nonce:    nonce,
			GasPrice: gasPrice,
			Gas:      gasLimit,
			Value:    ethAmount,
			Data:     data,
		})
	}
}

// CheckDynamicFeeTxAvailable use to detect if current network that connect via node url is support dynamic fee tx,
// this is done by a trick where we check if block info contain baseFee > 0, that may not always work but should enough
// for now.
func CheckDynamicFeeTxAvailable(ethclient *ethclient.Client) bool {
	block, err := ethclient.BlockByNumber(context.Background(), nil)
	if err != nil {
		fmt.Printf("couldn't get block by number: %s\n", err)
		return false
	}
	return block.BaseFee() != nil && block.BaseFee().Cmp(common.Big0) > 0
}
