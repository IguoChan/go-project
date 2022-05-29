package appx

type Worker interface {
	Start() error
	Stop() error
}
