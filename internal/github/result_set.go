package github

import "fmt"

// I dont like this, im testing things out
type Resultset struct {
	Name    string
	Version string
	Change        string
}

var resultSet []Resultset

func NewDefaultResultSet(){
	resultSet = []Resultset{}
}

func AddResult(rs Resultset) {
	resultSet = append(resultSet, rs)
}

func GetDefaultResultSet() []Resultset {
	return resultSet
}

func Print() string {
	var output string
	for _, rs := range resultSet {
		output += fmt.Sprintf("| %s | %s | %s |\n", rs.Name, rs.Version, rs.Change)

	}
	return output
}