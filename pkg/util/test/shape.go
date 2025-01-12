// SPDX-License-Identifier: AGPL-3.0-only

package test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xlab/treeprint"
)

// We use unsafe casting to convert between some Mimir and Prometheus types
// For this to be safe, the two types need to have the same shape (same fields
// in the same order). This function requires that this property is maintained.
// The fields do not need to have the same names to make the conversion safe,
// but we also check the names are the same here to ensure there's no confusion
// (eg. two bool fields swapped).
func RequireSameShape(t *testing.T, expectedType reflect.Type, actualType reflect.Type) {
	expectedFormatted := prettyPrintType(expectedType)
	actualFormatted := prettyPrintType(actualType)

	require.Equal(t, expectedFormatted, actualFormatted)
}

func prettyPrintType(t reflect.Type) string {
	if t.Kind() != reflect.Struct {
		panic(fmt.Sprintf("expected %s to be a struct but is %s", t.Name(), t.Kind()))
	}

	tree := treeprint.NewWithRoot("<root>")
	addTypeToTree(t, tree)

	return tree.String()
}

func addTypeToTree(t reflect.Type, tree treeprint.Tree) {
	if t.Kind() == reflect.Pointer {
		name := fmt.Sprintf("%s: *%s", t.Name(), t.Elem().Kind())
		addTypeToTree(t.Elem(), tree.AddBranch(name))
		return
	}

	if t.Kind() != reflect.Struct {
		panic(fmt.Sprintf("unexpected kind %s", t.Kind()))
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		switch f.Type.Kind() {
		case reflect.Pointer:
			name := fmt.Sprintf("+%v %s: *%s", f.Offset, f.Name, f.Type.Elem().Kind())
			addTypeToTree(f.Type.Elem(), tree.AddBranch(name))
		case reflect.Slice:
			name := fmt.Sprintf("+%v %s: []%s", f.Offset, f.Name, f.Type.Elem().Kind())

			if isPrimitive(f.Type.Elem().Kind()) {
				tree.AddNode(name)
			} else {
				addTypeToTree(f.Type.Elem(), tree.AddBranch(name))
			}
		default:
			name := fmt.Sprintf("+%v %s: %s", f.Offset, f.Name, f.Type.Kind())
			tree.AddNode(name)
		}
	}
}

func isPrimitive(k reflect.Kind) bool {
	switch k {
	case reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr,
		reflect.Float32,
		reflect.Float64,
		reflect.Complex64,
		reflect.Complex128,
		reflect.String,
		reflect.UnsafePointer:
		return true
	default:
		return false
	}
}
