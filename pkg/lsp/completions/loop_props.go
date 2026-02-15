package completions

import "fmt"

func GetLoopProperties() []Completion {
	return []Completion{
		{
			Label: "index",
			Documentation: fmt.Sprintf("%s\n%s",
				"(property) index: int",
				"The current iteration of the loop. Starts with 0."),
		},
		{
			Label: "first",
			Documentation: fmt.Sprintf("%s\n%s",
				"(property) first: bool",
				"Returns `true` if this is the first iteration of the loop."),
		},
		{
			Label: "last",
			Documentation: fmt.Sprintf("%s\n%s",
				"(property) last: bool",
				"Returns `true` if this is the last iteration of the loop."),
		},
		{
			Label: "iter",
			Documentation: fmt.Sprintf("%s\n%s",
				"(property) iter: int",
				"The current iteration of the loop. Starts with 1."),
		},
	}
}
