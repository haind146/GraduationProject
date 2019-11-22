package blockchain

type FacadeInterface interface {
	GetBalance(address string) (float64, error)
	GetBlockHash(blockHeight int64) (string, error)
	GetBestBlock() (string, int64, error)
	NextBlock(blockHeight int64, blockHash string) (string, bool, error)
	ApplyNextBlock(blockHash string, blockHeight int64) error
	RevertBlock(blockHeight int64) error
}

func GetFacadeInterface(paymentMethodId uint) FacadeInterface {
	facade := &BtcFacade{}
	return facade
}