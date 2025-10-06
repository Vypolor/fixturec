//go:generate mockgen -destination=mock/mock_gen.go -package=mock -source=./interface.go
package mypackage2

type MyType2 interface {
	Call2() int
}
