package main

import (
	"context"
	_ "crypto/ecdsa"
	"ethclient/task2/counter"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	//// 加载环境变量
	//if err := godotenv.Load(); err != nil {
	//	log.Fatal("加载 .env 文件失败")
	//}
	//
	//// 配置
	//rpcURL := os.Getenv("SEPOLIA_RPC_URL")
	//if rpcURL == "" {
	//	rpcURL = "https://rpc.sepolia.org"
	//}

	// 连接到 Sepolia
	privateKeyHex := "your_private_key" // 替换为钱包私钥
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/pkkk")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// 读取合约地址
	contractAddressHex, err := os.ReadFile("contract_address.txt")
	if err != nil {
		log.Fatal("请先部署合约并生成 contract_address.txt 文件")
	}

	contractAddress := common.HexToAddress(string(contractAddressHex))

	// 创建合约实例
	instance, err := counter.NewCounter(contractAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("连接到合约: %s\n", contractAddress.Hex())

	// 1. 读取当前计数
	count, err := instance.GetCount(nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("当前计数: %d\n", count)

	// 2. 获取合约所有者
	owner, err := instance.Owner(nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("合约所有者: %s\n", owner.Hex())

	// 3. 递增计数（需要交易）
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}

	//publicKey := privateKey.Public()
	//publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	//if !ok {
	//	log.Fatal("无法获取公钥")
	//}
	//fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 获取链 ID
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}

	// 设置 gas 限制和价格
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice

	fmt.Println("\n正在递增计数...")

	// 发送递增交易
	tx, err := instance.Increment(auth)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("交易已发送: %s\n", tx.Hash().Hex())
	fmt.Println("等待交易确认...")

	// 等待交易确认
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		log.Fatal(err)
	}

	if receipt.Status == 1 {
		fmt.Println("✓ 交易成功！")

		// 再次读取计数
		newCount, err := instance.GetCount(nil)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("新计数: %d\n", newCount)

		// 检查事件日志
		fmt.Println("\n检查事件日志...")
		blockNumber := receipt.BlockNumber.Uint64()
		filterOpts := &bind.FilterOpts{
			Start:   blockNumber,
			End:     &blockNumber,
			Context: context.Background(),
		}

		events, err := instance.FilterCountIncrement(filterOpts, nil)
		if err != nil {
			log.Printf("无法获取事件: %v", err)
		} else {
			for events.Next() {
				event := events.Event
				fmt.Printf("事件: 计数递增到 %d by %s\n",
					event.NewCount, event.Incrementer.Hex())
			}
		}
	} else {
		fmt.Println("✗ 交易失败")
	}

	// 4. 查询区块信息
	block, err := client.BlockByNumber(context.Background(), receipt.BlockNumber)
	if err != nil {
		log.Printf("无法获取区块信息: %v", err)
	} else {
		fmt.Printf("\n区块 #%d 信息:\n", block.Number().Int64())
		fmt.Printf("时间戳: %s\n", time.Unix(int64(block.Time()), 0))
		fmt.Printf("交易数量: %d\n", len(block.Transactions()))
	}
}

// 运行结果：
//连接到合约: 0x5E1bD45498136DD390DB9d5bb826d2b33BD5d8A3
//当前计数: 0
//合约所有者: 0x5a8F7793a036b4e5DffEbd05CBe93bFEE6E3FeCa
//
//正在递增计数...
//交易已发送: 0x8ae6751f0d2f98350b2f4eff436fc9a678bf12afc1c4a8bfb171f986085686a2
//等待交易确认...
//✓ 交易成功！
//新计数: 1
//
//检查事件日志...
//事件: 计数递增到 1 by 0x5a8F7793a036b4e5DffEbd05CBe93bFEE6E3FeCa
//
//区块 #9990600 信息:
//时间戳: 2026-01-06 22:09:24 +0800 CST
//交易数量: 114
//
//Process finished with the exit code 0
