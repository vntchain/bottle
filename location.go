package main

type Location struct {
	Path   string
	Offset int
	Size   int
}

func NewLocation(path string, offset int, size int) Location {
	return Location{
		Path:   path,
		Offset: offset,
		Size:   size,
	}
}
