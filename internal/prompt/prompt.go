package prompt

import (
	"errors"
	"fmt"
	"github.com/gabemanfroi/go_crud/internal/utils"
	"github.com/manifoldco/promptui"
	"strconv"
)

var templates = &promptui.PromptTemplates{
	Prompt:  "{{ . }} : ",
	Valid:   "{{ . | green }} : ",
	Invalid: "{{ . | red }} : ",
	Success: "{{ . | bold }} : ",
}

func GetInput(c Content) string {
	validate := func(input string) error {
		if len(input) <= 0 {
			return errors.New(c.ErrorMsg)
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:     c.Label,
		Templates: templates,
		Validate:  validate,
		Pointer:   promptui.PipeCursor,
	}

	result, err := prompt.Run()
	utils.Check(err, "Prompt failed")

	fmt.Printf("You have entered: %s\n", result)

	return result
}

func GetSelection(c Content, items []string) string {

	prompt := promptui.Select{
		Label: c.Label,
		Items: items,
	}

	_, result, err := prompt.Run()
	utils.Check(err, "Prompt failed")
	return result
}

func GetBoolean(c Content) bool {
	items := []string{"Yes", "No"}

	prompt := promptui.Select{
		Label: c.Label,
		Items: items,
	}

	_, result, err := prompt.Run()
	utils.Check(err, "Prompt failed")
	return result == "Yes"
}

func GetNumber(c Content) uint {
	validate := func(input string) error {
		if len(input) <= 0 {
			return errors.New(c.ErrorMsg)
		}
		parsedInput, err := strconv.Atoi(input)
		if err != nil {
			return errors.New(c.ErrorMsg)
		}
		if parsedInput <= 0 {
			return errors.New(c.ErrorMsg)
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:     c.Label,
		Templates: templates,
		Validate:  validate,
		Pointer:   promptui.PipeCursor,
	}

	result, err := prompt.Run()
	utils.Check(err, "Prompt failed")

	fmt.Printf("You have entered: %s\n", result)

	parsedResult, err := strconv.Atoi(result)

	utils.Check(err, "Failed to parse result")

	return uint(parsedResult)
}
