package base

type Translate interface {
	GetTranslate(text string)(error, string, string,[]string)
	GetName()string
}
