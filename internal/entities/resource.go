package entities

type Resource struct {
	URI        string
	Filename   string
	Downloaded bool
	Archived   bool
	Error      error
}

func NewResource(url string) Resource {
	return Resource{
		URI:        url,
		Filename:   "",
		Downloaded: false,
		Error:      nil,
	}
}
