/*
Copyright © 2023 Dimitris Dalianis <dimitris@dalianis.gr>
This file is part of CLI application keyvalDetector
*/
package cmd

import "fmt"

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
