/*
Copyright © 2023 Dimitris Dalianis <dimitris@dalianis.gr>
This file is part of CLI application keyvalDetector
*/
package cmd

import "testing"

func TestIsSystemConfigMap(t *testing.T) {
	testCases := []struct {
		name     string
		element  string
		expected bool
	}{
		{
			name:     "Empty string",
			element:  "",
			expected: false,
		},
		{
			name:     "Non system configmap",
			element:  "my-configmap",
			expected: false,
		},
		{
			name:     "System configmap",
			element:  "cluster-info",
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := isSystemConfigMap(tc.element)
			if result != tc.expected {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}

}

func TestContains(t *testing.T) {
	testCases := []struct {
		name     string
		input    []string
		element  string
		expected bool
	}{
		{
			name:     "Empty slice",
			input:    []string{},
			element:  "abc",
			expected: false,
		},
		{
			name:     "Element present in slice",
			input:    []string{"abc", "def", "ghi"},
			element:  "def",
			expected: true,
		},
		{
			name:     "Element not present in slice",
			input:    []string{"abc", "def", "ghi"},
			element:  "xyz",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := contains(tc.input, tc.element)
			if result != tc.expected {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}
