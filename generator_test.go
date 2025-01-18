package main

import (
	"testing"
)

func TestThatACommonUrlCanBeCreated(t *testing.T) {
	g, err := NewImgproxyUrlGenerator(Config{
		Host: "https://example.com",
		Disk: "s3://bucket",
		//Key:  "",
		//Salt: "",
	})
	if err != nil {
		t.Error(err)
	}

	filename := g.File("file.jpg").Width(500).Height(300).Quality(92).Crop().Get()

	// "{$host}/insecure/rs:{$crop}:{$width}:{$height}/g:{$this->gravity}/plain/s3://bucket/{$this->filenameWithExtension}";
	expected := "https://example.com/insecure/rs:fill:500:300/g:ce/q:92/plain/s3://bucket/file.jpg.webp"
	if filename != expected {
		t.Errorf("Expected %s, got %s", expected, filename)
	}
}

func TestThatACommonUrlCanBeEncrypted(t *testing.T) {
	g, err := NewImgproxyUrlGenerator(Config{
		Host: "http://localhost:9050",
		Disk: "s3://bucket",
		// Key:  "943b421c9eb07c830af81030552c86009268de4e532ba2ee2eab8247c6da0881",
		// Salt: "520f986b998545b4785e0defbc4f3c1203f22de2374a3d53cb7a7fe9fea309c5",
		EncodeUrl: true,
	})
	if err != nil {
		t.Error(err)
	}

	filename := g.File("file.jpg").Width(222).Height(333).Quality(92).Crop().Format(WEBP).Get()

	// "{$host}/insecure/rs:{$crop}:{$width}:{$height}/g:{$this->gravity}/plain/s3://bucket/{$this->filenameWithExtension}";
	expected := "http://localhost:9050/insecure/rs:fill:222:333/g:ce/q:92/czM6Ly9idWNrZXQv/ZmlsZS5qcGc.webp"
	if filename != expected {
		t.Errorf("Expected %s, got %s", expected, filename)
	}
}

func TestEncryptingUrls(t *testing.T) {
	g, err := NewImgproxyUrlGenerator(Config{
		Host:      "https://example.com",
		Disk:      "s3://bucket",
		Key:       "11111111111111",
		Salt:      "22222222222222",
		EncodeUrl: true,
	})
	if err != nil {
		t.Error(err)
	}

	filename := g.File("file.jpg").Width(222).Height(333).Quality(92).Crop().Format(WEBP).Get()

	// "{$host}/insecure/rs:{$crop}:{$width}:{$height}/g:{$this->gravity}/plain/s3://bucket/{$this->filenameWithExtension}";
	expected := "https://example.com/ssPyruH8et370vhYRwrvX3f-dlC2hHbSZPHwRICJEhA/rs:fill:222:333/g:ce/q:92/czM6Ly9idWNrZXQv/ZmlsZS5qcGc.webp"
	if filename != expected {
		t.Errorf("Expected %s, got %s", expected, filename)
	}
}
