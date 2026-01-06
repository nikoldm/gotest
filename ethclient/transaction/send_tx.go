package transaction

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	_ "github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// SendTransaction 发送以太币转账交易
func SendTransaction(client *ethclient.Client, privateKeyStr string, fromAddress common.Address, toAddress common.Address, amount *big.Int) (string, error) {

	// 1. 解析私钥
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %v", err)
	}

	// 2. 获取公钥和地址
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("failed to cast public key to ECDSA")
	}

	fromAddr := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Printf("发送方地址: %s\n", fromAddr.Hex())

	// 3. 获取链ID
	ctx := context.Background()
	chainID, err := client.NetworkID(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get network ID: %v", err)
	}

	// 4. 获取nonce
	nonce, err := client.PendingNonceAt(ctx, fromAddr)
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %v", err)
	}
	fmt.Printf("Nonce: %d\n", nonce)

	// 5. 获取Gas价格
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to suggest gas price: %v", err)
	}
	fmt.Printf("Gas 价格: %d wei\n", gasPrice)

	// 6. 获取Gas限制
	gasLimit := uint64(21000) // 标准转账的Gas限制

	// 7. 构造交易
	tx := types.NewTransaction(
		nonce,
		toAddress,
		amount,
		gasLimit,
		gasPrice,
		nil,
	)

	// 8. 签名交易
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %v", err)
	}

	// 9. 发送交易
	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}

	// 10. 返回交易哈希
	txHash := signedTx.Hash().Hex()
	fmt.Printf("交易已发送!\n")
	fmt.Printf("交易哈希: %s\n", txHash)
	fmt.Printf("发送方: %s\n", fromAddr.Hex())
	fmt.Printf("接收方: %s\n", toAddress.Hex())
	fmt.Printf("金额: %s ETH\n", weiToEther(amount))
	fmt.Printf("Gas 价格: %s Gwei\n", weiToGwei(gasPrice))
	fmt.Printf("Gas 限制: %d\n", gasLimit)
	fmt.Printf("预计费用: %s ETH\n", weiToEther(calculateFee(gasPrice, gasLimit)))

	return txHash, nil
}

// CheckBalance 检查账户余额
func CheckBalance(client *ethclient.Client, address common.Address) error {
	ctx := context.Background()

	balance, err := client.BalanceAt(ctx, address, nil)
	if err != nil {
		return fmt.Errorf("failed to get balance: %v", err)
	}

	fmt.Printf("地址 %s 的余额: %s ETH\n", address.Hex(), weiToEther(balance))
	return nil
}

// WaitForTransaction 等待交易确认
func WaitForTransaction(client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	ctx := context.Background()

	fmt.Println("\n⏳ 等待交易确认...")

	// 方法1: 使用 bind.WaitMined（需要实际的交易对象）
	// 但由于我们没有原始交易对象，使用轮询方法

	// 方法2: 轮询检查交易收据
	for {
		receipt, err := client.TransactionReceipt(ctx, txHash)
		if err != nil {
			if err.Error() == "not found" {
				// 交易还未被打包，等待后重试
				fmt.Print(".")
				time.Sleep(5 * time.Second)
				continue
			}
			return nil, fmt.Errorf("failed to get transaction receipt: %v", err)
		}

		// 交易已确认
		fmt.Printf("\n交易已确认!\n")
		fmt.Printf("区块号: %d\n", receipt.BlockNumber.Uint64())
		fmt.Printf("状态: %s\n", getStatusString(receipt.Status))
		fmt.Printf("Gas 使用量: %d\n", receipt.GasUsed)
		fmt.Printf("实际费用: %s ETH\n",
			weiToEther(new(big.Int).Mul(receipt.EffectiveGasPrice, big.NewInt(int64(receipt.GasUsed)))))

		return receipt, nil
	}
}

// 辅助函数
func weiToEther(wei *big.Int) string {
	// 1 ETH = 10^18 wei
	ether := new(big.Float).SetInt(wei)
	ether.Quo(ether, big.NewFloat(1e18))
	return ether.Text('f', 8)
}

func weiToGwei(wei *big.Int) string {
	gwei := new(big.Float).SetInt(wei)
	gwei.Quo(gwei, big.NewFloat(1e9))
	return gwei.Text('f', 2)
}

func calculateFee(gasPrice *big.Int, gasLimit uint64) *big.Int {
	fee := new(big.Int).Mul(gasPrice, big.NewInt(int64(gasLimit)))
	return fee
}

func getStatusString(status uint64) string {
	if status == 1 {
		return "成功"
	}
	return "失败"
}
