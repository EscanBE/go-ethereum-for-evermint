package vm

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
)

var _ PrecompiledContract = &CustomPrecompiledContract{}

// CustomPrecompiledContract is a precompiled contract that externally provided.
type CustomPrecompiledContract struct {
	name     string                            // name of the contract
	address  common.Address                    // address of the contract
	methods  []CustomPrecompiledContractMethod // methods of the contract
	disabled bool                              // disabled flag indicate the contract is disabled
}

func NewCustomPrecompiledContract(
	address common.Address, methods []CustomPrecompiledContractMethod, name string,
) PrecompiledContract {
	if address == (common.Address{}) {
		panic("invalid address")
	}

	if len(methods) == 0 {
		panic("no methods registered")
	}

	uniqueSig := make(map[string]struct{})
	for _, m := range methods {
		hexSig := hex.EncodeToString(m.Method4BytesSignatures)
		if err := m.Validate(); err != nil {
			panic(fmt.Sprintf("invalid method %s: %s", hexSig, err))
		}
		if _, exists := uniqueSig[hexSig]; exists {
			panic(fmt.Sprintf("duplicate method %s", hexSig))
		}
	}

	return &CustomPrecompiledContract{
		address: address,
		methods: methods,
		name:    name,
	}
}

func (s CustomPrecompiledContract) RequiredGas(input []byte) uint64 {
	sig := input[:4]

	for _, method := range s.methods {
		if bytes.Equal(method.Method4BytesSignatures, sig) {
			return method.RequireGas
		}
	}

	return 0
}

// Run runs the precompiled contract.
func (s CustomPrecompiledContract) Run(_ []byte) ([]byte, error) {
	panic("use RunCustom instead")
}

// RunCustom runs the custom precompiled contract with extra arguments supports state-modifying contracts.
func (s CustomPrecompiledContract) RunCustom(caller ContractRef, input []byte, readOnly bool, evm *EVM) ([]byte, error) {
	sig := input[:4]

	for _, method := range s.methods {
		if bytes.Equal(method.Method4BytesSignatures, sig) {
			if readOnly && !method.ReadOnly {
				return nil, ErrWriteProtection
			}

			return method.Executor.Execute(caller, s.address, input, evm)
		}
	}

	return nil, ErrExecutionReverted
}

func (s CustomPrecompiledContract) Name() string {
	return s.name
}

func (s CustomPrecompiledContract) Address() common.Address {
	return s.address
}

func (s *CustomPrecompiledContract) WithDisabled(disabled bool) *CustomPrecompiledContract {
	s.disabled = disabled
	return s
}

type CustomPrecompiledContractMethod struct {
	Method4BytesSignatures []byte
	RequireGas             uint64
	ReadOnly               bool
	Executor               CustomPrecompiledContractMethodExecutorI
}

func (m CustomPrecompiledContractMethod) Validate() error {
	if len(m.Method4BytesSignatures) != 4 {
		return fmt.Errorf("invalid method signature, expected 4 bytes, got %d", len(m.Method4BytesSignatures))
	}

	if m.ReadOnly {
		// allow any gas requirement for read-only methods
	} else {
		if m.RequireGas == 0 {
			return fmt.Errorf("invalid gas requirement, expected non-zero value")
		}
	}

	if m.Executor == nil {
		return fmt.Errorf("missing executor")
	}

	return nil
}

type CustomPrecompiledContractMethodExecutorI interface {
	// Execute executes the method with the given input and returns the output.
	Execute(caller ContractRef, contractAddress common.Address, input []byte, evm *EVM) ([]byte, error)
}
