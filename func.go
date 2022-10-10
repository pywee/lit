package goExpr

func checkFunctionName(name string) string {
	switch name {
	case "print":
		return ""
	}
	return "-"
}
