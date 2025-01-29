package usecase

type DataProcessorUsecase interface {
	ProcessData(flag bool, links []string) ([][]byte, error)
}
