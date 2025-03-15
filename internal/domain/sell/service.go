package sell

type SellService struct {
	repository SellRepository
}

func NewService(repository SellRepository) *SellService {
	return &SellService{repository: repository}
}
