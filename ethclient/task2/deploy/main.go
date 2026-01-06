package main

import (
	"context"
	"crypto/ecdsa"
	"ethclient/task2/counter"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

/**
 * 部署到sepolia
 */
func main() {
	// 连接到 Sepolia 测试网
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/pkkkk")
	if err != nil {
		log.Fatal(err)
	}

	//privateKey, err := crypto.GenerateKey()
	//privateKeyBytes := crypto.FromECDSA(privateKey)
	//privateKeyHex := hex.EncodeToString(privateKeyBytes)
	//fmt.Println("Private Key:", privateKeyHex)
	// 私钥
	privateKey, err := crypto.HexToECDSA("your_private_key") // 替换为钱包私钥
	if err != nil {
		log.Fatal(err)
	}

	// 获取公钥地址
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	chainId, err := client.NetworkID(context.Background()) // 11155111
	if err != nil {
		log.Fatal(err)
	}
	// 创建认证
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		log.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	// 部署合约
	address, tx, instance, err := counter.DeployCounter(auth, client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(address.Hex())   //
	fmt.Println(tx.Hash().Hex()) // 0x75ece13e173d43afc56b2b39f7c149cd71dacac45d653848a3e2651c52fa913c

	// 等待交易确认
	fmt.Println("等待交易确认...")
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.Fatal(err)
	}

	if receipt.Status == 1 {
		fmt.Println("✓ 合约部署确认！")
		fmt.Printf("区块号: %d\n", receipt.BlockNumber)
		fmt.Printf("Gas 使用量: %d\n", receipt.GasUsed)

		// 测试调用合约
		count, err := instance.GetCount(nil)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("初始计数: %d\n", count)
	} else {
		fmt.Println("✗ 合约部署失败")
	}

	// 保存合约地址到文件
	saveContractAddress(address.Hex())
	_ = instance
}

func saveContractAddress(address string) {
	file, err := os.Create("contract_address.txt")
	if err != nil {
		log.Printf("无法保存合约地址: %v", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(address)
	if err != nil {
		log.Printf("写入合约地址失败: %v", err)
	}
}

//结果：
//0x5E1bD45498136DD390DB9d5bb826d2b33BD5d8A3
//0x9ccbbf335f6ea31da59ec5e50bb2bb5d2509fcb63c2c3b6cd8bb043af592052a
//等待交易确认...
//✓ 合约部署确认！
//区块号: 9990582
//Gas 使用量: 298030
//初始计数: 0
