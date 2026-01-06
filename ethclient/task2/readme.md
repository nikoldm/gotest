1、安装solc
    npm install -g solc
2、编译合约代码
    solcjs --bin .\Counter.sol
3、生成合约 abi 文件
    solcjs --abi Counter.so
4、安装abigen:
    go install github.com/ethereum/go-ethereum/cmd/abigen@latest
5、使用 abigen 工具根据这两个生成 bin 文件和 abi 文件，生成 go 代码：
    abigen --bin=Counter_sol_Counter.bin --abi=Counter_sol_Counter.abi --pkg=counter --out=counter/counter.go
6、部署合约
