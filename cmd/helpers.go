/*
Copyright © 2023 Dimitris Dalianis <dimitris@dalianis.gr>
This file is part of CLI application keyvalDetector
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

func printSlice(s []string) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Color codes: https://github.com/fatih/color/blob/main/color.go#L69
func colorPrint(colorCode int, input string) {
	colored := fmt.Sprintf("\x1b[%dm%s\x1b[0m", colorCode, input)
	fmt.Print(colored)
}

func printTable(tabledata [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Namespace"})

	for _, v := range tabledata {
		table.Append(v)
	}
	table.Render() // Send output
}
