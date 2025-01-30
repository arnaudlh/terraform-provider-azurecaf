package testutils

type ResourceTestData struct {
	ResourceType    string
	MinLength      int
	MaxLength      int
	ValidationRegex string
	Scope          string
	Slug           string
	Dashes         bool
	LowerCase      bool
	Regex          string
}
