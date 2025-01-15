package generator

import "testing"

func TestThatACommonUrlCanBeCreated(t *testing.T) {
	g := NewImgproxyUrlGenerator("https://example.com", "s3")
	filename := g.From("file.jpg").Width(500).Height(300).Quality(92).Crop().Get()

	// "{$host}/insecure/rs:{$crop}:{$width}:{$height}/g:{$this->gravity}/plain/s3://youshouldsurf/{$this->filenameWithExtension}";
	expected := "https://example.com/insecure/rs:fill:500:300/g:ce/q:92/plain/s3://file.jpg"
	if filename != expected {
		t.Errorf("Expected %s, got %s", expected, filename)
	}
}
