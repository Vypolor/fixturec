//go:generate mockgen -destination=mock/mock_gen.go -package=mock -source=./interface.go
package mypackage1

type MyType1 interface {
	Call1() int
}
