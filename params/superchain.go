package params

import (
	"fmt"
	"math/big"

	"github.com/ethereum-optimism/superchain-registry/superchain"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
)

func LoadOPStackChainConfig(chainID uint64) (*ChainConfig, error) {
	chConfig, ok := superchain.OPChains[chainID]
	if !ok {
		return nil, fmt.Errorf("unknown chain ID: %d", chainID)
	}
	superchainConfig, ok := superchain.Superchains[chConfig.Superchain]
	if !ok {
		return nil, fmt.Errorf("unknown superchain %q: %w", chConfig.Superchain)
	}

	genesisActivation := uint64(0)
	out := &ChainConfig{
		ChainID:                       new(big.Int).SetUint64(chainID),
		HomesteadBlock:                common.Big0,
		DAOForkBlock:                  nil,
		DAOForkSupport:                false,
		EIP150Block:                   common.Big0,
		EIP155Block:                   common.Big0,
		EIP158Block:                   common.Big0,
		ByzantiumBlock:                common.Big0,
		ConstantinopleBlock:           common.Big0,
		PetersburgBlock:               common.Big0,
		IstanbulBlock:                 common.Big0,
		MuirGlacierBlock:              common.Big0,
		BerlinBlock:                   common.Big0,
		LondonBlock:                   common.Big0,
		ArrowGlacierBlock:             common.Big0,
		GrayGlacierBlock:              common.Big0,
		MergeNetsplitBlock:            common.Big0,
		ShanghaiTime:                  nil,
		CancunTime:                    nil,
		PragueTime:                    nil,
		BedrockBlock:                  common.Big0,
		RegolithTime:                  &genesisActivation,
		TerminalTotalDifficulty:       common.Big0,
		TerminalTotalDifficultyPassed: true,
		Ethash:                        nil,
		Clique:                        nil,
		Optimism:                      &OptimismConfig{
			EIP1559Elasticity:  6,
			EIP1559Denominator: 50,
		},
	}

	// note: no actual parameters are being loaded, yet.
	// Future superchain upgrades are loaded from the superchain chConfig and applied to the geth ChainConfig here.
	_ = superchainConfig.Config

	// special overrides for OP-Stack chains with pre-Regolith upgrade history
	switch chainID {
	case 420:
		out.LondonBlock = big.NewInt(4061224)
		out.ArrowGlacierBlock = big.NewInt(4061224)
		out.GrayGlacierBlock = big.NewInt(4061224)
		out.MergeNetsplitBlock = big.NewInt(4061224)
		out.BedrockBlock = big.NewInt(4061224)
		out.RegolithTime = &OptimismGoerliRegolithTime
		out.Optimism.EIP1559Elasticity = 10
	case 10:
		out.BerlinBlock =                   big.NewInt(3950000)
		out.LondonBlock =                   big.NewInt(105235063)
		out.ArrowGlacierBlock =             big.NewInt(105235063)
		out.GrayGlacierBlock =              big.NewInt(105235063)
		out.MergeNetsplitBlock =            big.NewInt(105235063)
		out.BedrockBlock =                  big.NewInt(105235063)
	case 84531:
		out.RegolithTime = &BaseGoerliRegolithTime
	}

	return out, nil
}

func LoadOPStackGenesis(chainID uint64) (*core.Genesis, error) {
	chConfig, ok := superchain.OPChains[chainID]
	if !ok {
		return nil, fmt.Errorf("unknown chain ID: %d", chainID)
	}

	cfg, err := LoadOPStackChainConfig(chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to load params.ChainConfig for chain %d: %w", chainID, err)
	}

	genesis := &core.Genesis{
		Config:     cfg,
		Nonce:      0,
		Timestamp:  chConfig.Genesis.L2Time,
		ExtraData:  []byte("BEDROCK"),
		GasLimit:   30_000_000,
		Difficulty: nil,
		Mixhash:    common.Hash{},
		Coinbase:   common.Address{},
		Alloc:      nil,
		Number:     0,
		GasUsed:    0,
		ParentHash: common.Hash{},
		BaseFee:    nil,
	}

	// TODO: load state allocations

	// TODO: exceptions for OP-Mainnet and OP-Goerli to handle pre-Bedrock history

	if chConfig.Genesis.ExtraData != nil {
		genesis.ExtraData = *chConfig.Genesis.ExtraData
		if len(genesis.ExtraData) > 32 {
			return nil, fmt.Errorf("chain must have 32 bytes or less extra-data in genesis, got %d", len(genesis.ExtraData))
		}
	}
	// TODO: apply all genesis block values

	// Verify we correctly produced the genesis config by recomputing the genesis-block-hash
	genesisBlock := genesis.ToBlock()
	genesisBlockHash := genesisBlock.Hash()
	if [32]byte(chConfig.Genesis.L2.Hash) != genesisBlockHash {
		return nil, fmt.Errorf("produced genesis with hash %s but expected %s", genesisBlockHash, chConfig.Genesis.L2.Hash)
	}
	return genesis, nil
}

func SystemConfigAddr(chainID uint64) (common.Address, error) {
	// TODO(proto): when we move to CREATE-2 proxy addresses
	// for SystemConfig contracts we can deterministically compute the system config addr,
	// and do not have to load it from the superchain configs.
	chConfig, ok := superchain.OPChains[chainID]
	if !ok {
		return common.Address{}, fmt.Errorf("unknown chain ID: %d", chainID)
	}
	return common.Address(chConfig.SystemConfigAddr), nil
}
