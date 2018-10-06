package main

import "github.com/tkrajina/typescriptify-golang-structs/typescriptify"

type Address struct {
	// Used in html
	City    string  `json:"city"`
	Number  float64 `json:"number"`
	Country string  `json:"country,omitempty"`
}

type PersonalInfo struct {
	Hobbies []string `json:"hobby"`
	PetName string   `json:"pet_name"`
}

type Person struct {
	Name         string            `json:"name"`
	PersonalInfo PersonalInfo      `json:"personal_info"`
	Nicknames    []string          `json:"nicknames"`
	Addresses    []Address         `json:"addresses"`
	Address      *Address          `json:"address"`
	Children     map[string]Person `json:"children"`
	ChildrenAge  map[string]int    `json:"children_age"`
}

func main() {
	converter := typescriptify.New()
	converter.CreateFromMethod = true
	converter.Indent = "    "
	converter.BackupDir = ""

	converter.Add(Person{})

	err := converter.ConvertToFile("ts_tests/example_output.ts")
	if err != nil {
		panic(err.Error())
	}
}