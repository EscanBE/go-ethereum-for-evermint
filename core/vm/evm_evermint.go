package vm

import "github.com/ethereum/go-ethereum/common"

func (evm *EVM) WithCustomPrecompiledContract(precompiles ...PrecompiledContract) *EVM {
	evm.customPrecompiledContracts = make(map[common.Address]PrecompiledContract, len(precompiles))
	for _, precompile := range precompiles {
		evm.customPrecompiledContracts[precompile.(*CustomPrecompiledContract).Address()] = precompile
	}
	return evm
}
