/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/gabemanfroi/go_crud/app/generate"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate [modelName]",
	Short: "generates necessary files for a basic model CRUD",
	Long: `Generates Necessary Files for the model passed as argument.
	
	Files:
		-model_controller.go
		-model_controller_interface.go
		-model_service.go		
		-model_service_interface.go
		-model_repository.go
		-model_repository_interface.go
`,
	Run: func(cmd *cobra.Command, args []string) {
		model, _ := cmd.Flags().GetString("model")

		if model != "" {
			generate.Generate(model)
		}

	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	generateCmd.PersistentFlags().String("model", "", "Model to generate files for")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
