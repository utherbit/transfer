package main

import (
	"go/parser"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatExpr(t *testing.T) {
	tests := []struct {
		expr       string
		importsMap map[string]string

		wantImportsSet map[string]struct{}
		wantFormatExpr string
	}{
		{
			expr:           "string",
			importsMap:     map[string]string{},
			wantImportsSet: map[string]struct{}{},
			wantFormatExpr: "string",
		},
		{
			expr:           "*string",
			importsMap:     map[string]string{},
			wantImportsSet: map[string]struct{}{},
			wantFormatExpr: "*string",
		},
		{
			expr: "uuid.UUID",
			importsMap: map[string]string{
				"uuid": "github.com/google/uuid",
			},
			wantImportsSet: map[string]struct{}{
				"github.com/google/uuid": {},
			},
			wantFormatExpr: "uuid.UUID",
		},
		{
			expr: "map[uuid.UUID]int",
			importsMap: map[string]string{
				"uuid": "github.com/google/uuid",
			},
			wantImportsSet: map[string]struct{}{
				"github.com/google/uuid": {},
			},
			wantFormatExpr: "map[uuid.UUID]int",
		},
	}

	for _, test := range tests {
		t.Run(test.expr, func(t *testing.T) {
			expr, err := parser.ParseExpr(test.expr)
			require.NoError(t, err)

			importsSet := make(map[string]struct{})
			formated := formatExpr(expr, test.importsMap, importsSet)

			require.Equal(t, test.wantImportsSet, importsSet)
			require.Equal(t, test.wantFormatExpr, formated)
		})
	}
	// fmt.Printf("%#v", expr)
}

func TestCapitalize(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{
			in:  "foo",
			out: "Foo",
		},
		{
			in:  "foo_bar",
			out: "Foo_bar",
		},
		{
			in:  "id",
			out: "ID",
		},
		{
			in:  "idUser",
			out: "IDUser",
		},
		{
			in:  "idid",
			out: "IDid",
		},
	}

	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			require.Equal(t, test.out, capitalize(test.in))
		})
	}
}

func Test_findFileByRef(t *testing.T) {
	_, err := findStructByRef("./../../../internal/integrations/iiko/internal/domain/terminal/brand_settings.go:5")

	require.NoError(t, err)
}
