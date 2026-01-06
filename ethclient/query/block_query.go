package query

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

/**
 * 任务 1：区块链读写 任务目标
	使用 Sepolia 测试网络实现基础的区块链交互，包括查询区块和发送交易。
 具体任务
环境搭建
	安装必要的开发工具，如 Go 语言环境、 go-ethereum 库。
	注册 Infura 账户，获取 Sepolia 测试网络的 API Key。
查询区块
	编写 Go 代码，使用 ethclient 连接到 Sepolia 测试网络。
	实现查询指定区块号的区块信息，包括区块的哈希、时间戳、交易数量等。
	输出查询结果到控制台。
发送交易
	准备一个 Sepolia 测试网络的以太坊账户，并获取其私钥。
	编写 Go 代码，使用 ethclient 连接到 Sepolia 测试网络。
	构造一笔简单的以太币转账交易，指定发送方、接收方和转账金额。
	对交易进行签名，并将签名后的交易发送到网络。
	输出交易的哈希值。
*/

// QueryBlock 查询指定区块的信息
func QueryBlock(client *ethclient.Client, blockNumber *big.Int) error {
	ctx := context.Background()
	// 获取区块信息
	block, err := client.BlockByNumber(ctx, blockNumber)
	if err != nil {
		return fmt.Errorf("QueryBlock err: %v", err)
	}
	// 输出区块信息
	fmt.Printf("区块高度：%d\n", block.Number().Int64())
	fmt.Printf("区块Hash：%s\n", block.Hash().Hex())
	fmt.Printf("区块时间戳: %d (%s)\n", block.Time(), formatTimestamp(int64(block.Time())))
	fmt.Printf("交易数量: %d\n", len(block.Transactions()))
	fmt.Printf("矿工地址: %s\n", block.Coinbase().Hex())
	fmt.Printf("区块难度: %d\n", block.Difficulty().Uint64())
	fmt.Printf("Gas 使用量: %d\n", block.GasUsed())
	fmt.Printf("Gas 限制: %d\n", block.GasLimit())
	fmt.Printf("父区块哈希: %s\n", block.ParentHash().Hex())

	// 如果又交易，显示前5笔交易
	if len(block.Transactions()) > 0 {
		fmt.Println("==== 前5比交易的哈希 ====")
		for idx, tx := range block.Transactions() {
			if idx >= 5 {
				break
			}
			fmt.Printf("第%d笔交易，哈希是：%s\n", idx+1, tx.Hash().Hex())
		}
	}
	return nil
}

func formatTimestamp(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return t.Format("2006-01-02 15:04:05")
}

// QueryLatestBlock 查询最新区块
func QueryLatestBlock(client *ethclient.Client) error {
	ctx := context.Background()
	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		return fmt.Errorf("QueryLatestBlock err: %v", err)
	}
	fmt.Printf("最新区块高度：%d\n", header.Number.Int64())
	return QueryBlock(client, header.Number)
}

func GetClient(rpcUrl string) (*ethclient.Client, error) {
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, fmt.Errorf("GetClient err: %v", err)
	}
	return client, nil
}

// GetTransactionReceipt 获取交易收据
func GetTransactionReceipt(client *ethclient.Client, txHash common.Hash) error {
	ctx := context.Background()
	receipt, err := client.TransactionReceipt(ctx, txHash)
	if err != nil {
		return fmt.Errorf("GetTransactionReceipt err: %v", err)
	}
	fmt.Println("=== Transaction Receipt ===")
	fmt.Printf("交易状态：%d（1=成功，0=失败）\n", receipt.Status)
	fmt.Printf("区块哈希: %s\n", receipt.BlockHash.Hex())
	fmt.Printf("区块号: %d\n", receipt.BlockNumber.Uint64())
	fmt.Printf("Gas 使用量: %d\n", receipt.GasUsed)
	fmt.Printf("交易索引: %d\n", receipt.TransactionIndex)

	return nil
}
