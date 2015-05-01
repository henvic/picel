package main

var existsDependencyCases = []ExistsDependencyProvider{
	{"echo", true},
	{"unknown", false},
}

var CheckMissingDependencies = []CheckMissingDependenciesProvider{
	{},
	{[]string{}, true},
	{[]string{"echo"}, true},
	{[]string{"echo", "cd"}, true},
	{[]string{"unknown"}, true},
	{[]string{"unknown", "unknown2"}, true},
	{[]string{"unknown", "echo"}, true},
}
