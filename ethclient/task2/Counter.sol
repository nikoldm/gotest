// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

/**
## 任务 2：合约代码生成 任务目标
使用 abigen 工具自动生成 Go 绑定代码，用于与 Sepolia 测试网络上的智能合约进行交互。
 具体任务
1. 编写智能合约
   - 使用 Solidity 编写一个简单的智能合约，例如一个计数器合约。
   - 编译智能合约，生成 ABI 和字节码文件。
2. 使用 abigen 生成 Go 绑定代码
   - 安装 abigen 工具。
   - 使用 abigen 工具根据 ABI 和字节码文件生成 Go 绑定代码。
3. 使用生成的 Go 绑定代码与合约交互
   - 编写 Go 代码，使用生成的 Go 绑定代码连接到 Sepolia 测试网络上的智能合约。
   - 调用合约的方法，例如增加计数器的值。
   - 输出调用结果。
 */
contract Counter {
    uint256 private count;
    address public owner;

    event CountIncrement(uint256 newCount, address indexed incrementer);
    constructor(){
        count = 0;
        owner = msg.sender;
    }

    function Increment() public returns (uint256) {
        count ++;
        emit CountIncrement(count, msg.sender);
        return count;
    }

    function getCount() public view returns (uint256){
        return count;
    }
    function setCount(uint256 newCount) public {
        require(msg.sender == owner, "Only owner can set");
        count = newCount;
    }

//    function getOwner()public view returns (address){
//        return owner;
//    }
}
