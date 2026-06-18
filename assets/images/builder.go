// Package images: builder_dreamborn. Describe URI Building interface
package images

type ImageURIBuilder interface {
	Card(ID string) string
}
