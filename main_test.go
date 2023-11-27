package main

import (
	"testing"
)

type Test struct {
	testName   string
	mockRecipe UpdateRecipeRequest
	mockId     int
	want       string
}

func TestBuildUpdateRecipeQuery(t *testing.T) {
	title := "Hello"
	mockUpdateRecipe := &UpdateRecipeRequest{
		Title: &title,
	}

	tests := []Test{
		{"I expect title to be included as value", *mockUpdateRecipe, 1, "UPDATE recipes SET title = ? WHERE id = ?"},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			query, _ := buildUpdateRecipeQuery(&test.mockRecipe, 1)

			if query != test.want {
				t.Error("Result query string does not match expected value!")
			}
		})
	}
}
