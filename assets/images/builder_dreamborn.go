// Pacakge images: builder_dreamborn. Contains all dreamboard URI building logic
package images

import (
	"fmt"
	"os"
)

type DreambornURIBuilder struct {
	BaseURI string
}

func New() ImageURIBuilder {
	return &DreambornURIBuilder{
		BaseURI: os.Getenv("CDN_URL"),
	}
}

func (b *DreambornURIBuilder) Card(ID string) string {
	return fmt.Sprintf("%s/%s", b.BaseURI, ID)
}
