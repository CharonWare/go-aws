package ui

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

func CreatePrompt(items []string, label string) (int, string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
		Size:  30,
		Searcher: func(input string, index int) bool {
			item := items[index]
			return containsIgnoreCase(item, input)
		},
	}

	index, output, err := prompt.Run()
	if err != nil {
		return 0, "", fmt.Errorf("prompt failed: %w", err)
	}
	return index, output, nil
}

func containsIgnoreCase(str, substr string) bool {
	return strings.Contains(strings.ToLower(str), strings.ToLower(substr))
}
