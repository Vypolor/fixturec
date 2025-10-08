package fixturec

import (
	"log"

	"github.com/Vypolor/fixturec/internal/generator"
)

type Config struct {
	// TypeName - name of the structure for which the fixture should be generated
	TypeName string
	// External - flag that points to generate mock for external struct dependencies
	External bool
}

func GenerateFixture(cfg Config) error {
	g := generator.NewGenerator(".", cfg.TypeName)

	err := g.Generate()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
