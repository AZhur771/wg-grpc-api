// Package usecase implements application business logic. Each logic group in own file.
package usecase

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type Peer interface {
}

type PeerRepo interface {
}
