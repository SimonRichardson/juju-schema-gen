package ast

type Location struct {
	File     string
	Line     int
	Position int
}

type Package struct {
	Location
	Name    string
	Facades []Facade
}

type Facade struct {
	Location
	Name    string
	Version int
	Params  []Data
	Results []Data
	Methods []Methods
}

type Data struct {
	Location
	Version int
	Name    string
	Type    string
}

type Methods struct {
	Location
	Name   string
	Inputs []Type
	Output []Type
}

type Type struct {
	Location
	Name string
}
