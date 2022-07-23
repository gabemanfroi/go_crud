/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"embed"
	"github.com/gabemanfroi/go_crud/cmd"
	"github.com/gabemanfroi/go_crud/internal/generate"
)

//go:embed templates/*
var templateFs embed.FS

func main() {
	generate.TemplateFs = templateFs
	cmd.Execute()
}
