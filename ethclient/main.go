package main

import (
	"ethclient/config"
	"ethclient/query"
	"ethclient/transaction"
	"flag"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

func main() {
	// 定义命令行参数：
	action := flag.String("action", "query", "执行的操作：query,send, balance,queryWait")
	blockNumber := flag.Int64("block", -1, "要查询的区块号")
	amount := flag.String("amount", "0.01", "转账金额（ETH）")
	flag.Parse()

	// 加载配置
	cfg := config.LoadConfig()

	if cfg.InfuraURL == "" {
		log.Fatal("请设置 InfuraURL 环境变量")
	}

	// 创建客户端
	client, err := query.GetClient(cfg.InfuraURL)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	//区块操作
	switch *action {
	case "query":
		if *blockNumber == -1 {
			// 查询最新区块
			fmt.Print("查询最新区块：")
			err = query.QueryLatestBlock(client)
		} else {
			// 查询指定区块
			fmt.Println("查询指定区块#：")
			query.QueryBlock(client, big.NewInt(*blockNumber))
		}
	case "send":
		if cfg.PrivateKey == "" || cfg.ToAddr == "" {
			log.Fatal("发送交易需要设置privateKey和to_address")
		}
		// 转换金额为 wei
		amountFloat, ok := new(big.Float).SetString(*amount)
		if !ok {
			log.Fatal("无效的金额格式。。。")
		}
		amountFloat.Mul(amountFloat, big.NewFloat(1e18))

		amountWei := new(big.Int)
		amountFloat.Int(amountWei)

		// 发送交易
		toAddr := common.HexToAddress(cfg.ToAddr)
		txHash, err := transaction.SendTransaction(client, cfg.PrivateKey, common.HexToAddress(cfg.FromAddr), toAddr, amountWei)
		if err == nil {
			fmt.Printf("\n可以在以下链接查看交易状态:\n")
			fmt.Printf("https://sepolia.etherscan.io/tx/%s\n", txHash)
		}
	case "balance":
		if cfg.FromAddr == "" {
			log.Fatal("查询余额需要设置 from_address 环境变量")
		}
		fromAddr := common.HexToAddress(cfg.FromAddr)
		err = transaction.CheckBalance(client, fromAddr)
	case "queryWait":
		txHash := common.HexToHash("0x205aa295f80c8defe0c50e5820ea0efc252b6b36df9437dedd8d733af64f9e48")
		transaction.WaitForTransaction(client, txHash)
	default:
		fmt.Printf("未知操作，请使用：query, send, balance\n")
		return
	}

	if err != nil {
		log.Fatal(err)
	}

}
