package jsync

import "reflect"

type tagProp struct {
	name string
	t    reflect.Type
}

func (tp tagProp) Name() string {
	return tp.name
}

func (tp tagProp) Type() reflect.Type {
	return tp.t
}

// extractTagMap retorna um map sendo a chave o nome do objeto json (valor da tag `json`) e o valor os dados da variável
// que detém a tag.
func extractTagMap(obj any) map[string]tagProp {
	t := reflect.TypeOf(obj)
	m := make(map[string]tagProp)

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		jsonV, ok := f.Tag.Lookup("json")
		if !ok {
			continue
		}

		m[jsonV] = tagProp{
			name: f.Name,
			t:    f.Type,
		}
	}

	return m
}
