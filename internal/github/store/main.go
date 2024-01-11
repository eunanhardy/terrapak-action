package store

import "fmt"

// I dont like this, im testing things out
type ResultStore struct {
	Name    string
	Version string
	Change  string
}

var resultSet []ResultStore

func New(){
	resultSet = []ResultStore{}
}

func (rs ResultStore) Add() {
	resultSet = append(resultSet, rs)
}

func GetDefaultResultSet() []ResultStore {
	return resultSet
}

func Print() string {
	var output string
	for _, rs := range resultSet {
		output += fmt.Sprintf("| %s | %s | %s | <br>", rs.Name, rs.Version, rs.Change)

	}
	return output
}