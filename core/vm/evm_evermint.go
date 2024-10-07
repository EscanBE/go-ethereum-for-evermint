package vm

import "github.com/ethereum/go-ethereum/common"

func (evm *EVM) WithCustomPrecompiledContracts(precompiles ...PrecompiledContract) *EVM {
	evm.customPrecompiledContracts = make(map[common.Address]PrecompiledContract, len(precompiles))
	for _, precompile := range precompiles {
		evm.customPrecompiledContracts[precompile.(*CustomPrecompiledContract).Address()] = precompile
	}
	return evm
}

func (evm *EVM) GetCustomPrecompiledContractsAddress() []common.Address {
	addrs := make([]common.Address, len(evm.customPrecompiledContracts))
	for addr := range evm.customPrecompiledContracts {
		addrs = append(addrs, addr)
	}
	return addrs
}
