package file_updater

type FileStore interface {
	Save() error
}
