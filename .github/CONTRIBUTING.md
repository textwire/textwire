# Textwire Contribution Guide

## Correct branch
When contributing to this repository, please make sure you are merging your changes into the correct branch matching the version you are working on. For example, if you are working on a new feature for the v2, you should merge your changes into the `v2` branch. If you are fixing a bug in the v1.*, you should merge your changes into the `main` branch, since first versions are corresponding to the `main` branch.

It's done this way because in Go packages, branch names like `v2`, `v3`, etc. are corresponding to the major versions of the package. For example, to use version 2 package (it's `v2` branch), you'll use `import "github.com/textwire/textwire/v2"` in your Go code.