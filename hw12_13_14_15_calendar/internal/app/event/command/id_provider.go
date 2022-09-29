package command

//go:generate mockgen -destination=./mock/id_provider_gen.go -package=mock . IDProvider

type IDProvider interface {
	GetID() (string, error)
}
