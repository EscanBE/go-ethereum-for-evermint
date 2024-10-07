package vm

import (
	"errors"
)

var (
	ErrDisabledPrecompile = errors.New("the precompile contract is disabled")
)

// RunPrecompiledContract runs a precompiled contract with the given input and supplied gas.
func (in *EVMInterpreter) RunPrecompiledContract(caller ContractRef, p PrecompiledContract, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	cpc, isCustomPrecompiledContract := p.(*CustomPrecompiledContract)
	if !isCustomPrecompiledContract {
		return RunPrecompiledContract(p, input, suppliedGas)
	}

	if cpc.disabled {
		return nil, suppliedGas, ErrDisabledPrecompile
	}

	if len(input) < 4 {
		return nil, 0, ErrExecutionReverted
	}

	gasCost := cpc.RequiredGas(input)
	if suppliedGas < gasCost {
		return nil, 0, ErrOutOfGas
	}
	suppliedGas -= gasCost

	{
		// We do not increase the call depth for precompiled contracts.
	}

	// Make sure the readOnly is only set if we aren't in readOnly yet.
	// This also makes sure that the readOnly flag isn't removed for child calls.
	if readOnly && !in.readOnly {
		in.readOnly = true
		defer func() { in.readOnly = false }()
	}

	output, err := cpc.RunCustom(caller, input, readOnly, in.evm)
	return output, suppliedGas, err
}
