package fixture

const templateName = "fixture"

const (
	mockAlias = "mockAlias"
	typeShort = "typeShort"
)

const fixtureTemplate = `
package {{ .PackageName }}

import (
	"context"
	"testing"

	ctrl_mock "go.uber.org/mock/gomock"
{{- range .Fields }}
	{{ mockAlias .PkgPath }} "{{ .PkgPath }}/mock"
{{- end }}
)

type fixture struct {
	ctrl *ctrl_mock.Controller
	ctx  context.Context

{{ range .Fields }}
	{{ .FieldName }} *{{ mockAlias .PkgPath }}.{{ typeShort .TypeName }}
{{- end }}

	{{ .ImplAlias }} *{{ .StructName }}
}

func setUp(t *testing.T) *fixture {
	f := &fixture{
		ctrl: ctrl_mock.NewController(t),
		ctx:  context.Background(),

{{ range .Fields }}
		{{ .FieldName }}: {{ mockAlias .PkgPath }}.New{{ typeShort .TypeName }}(ctrl_mock.NewController(t)),
{{- end }}
	}

	f.{{ .ImplAlias }} = &{{ .StructName }}{
{{- range .Fields }}
		{{ .FieldName }}: f.{{ .FieldName }},
{{- end }}
	}

	return f
}

func (f *fixture) tearDown() {
	f.ctrl.Finish()
}
`
