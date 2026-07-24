package nix

const NixCliName = "nix"

// NixBuildExpressionArgs builds arguments for building a Nix expression.
func NixBuildExpressionArgs(expression string) []string {
	return []string{
		"--extra-experimental-features",
		"nix-command flakes",
		"build",
		"--impure",
		"--no-link",
		"--print-out-paths",
		"--expr",
		expression,
	}
}
