package options

type TaskOptions struct {
	MaxNumResources       int
	MaxNumTasks           int
	AllowedFileExtensions map[string]struct{}
	FileDir               string
}
