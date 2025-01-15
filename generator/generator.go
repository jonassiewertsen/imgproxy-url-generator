package generator

import "fmt"

type Repository struct {
	Salt             []byte
	saltBin          []byte
	Key              []byte
	keyBin           []byte
	Host             string
	Disk             string // Cloud bucket
	QualityDefault   int
	EncryptionKey    *string
	encryptionKeyBin []byte
	PlainUrl         bool
}

type Generator struct {
	disk     string
	host     string
	fileName string
	width    int
	height   int
	quality  int
	crop     bool
	gravity  string
	format   Format
}

func NewImgproxyUrlGenerator(host, disk string) *Repository {
	return &Repository{
		Host: host,
		Disk: disk,
	}
}

func (g *Repository) From(fileName string) *Generator {
	if g.QualityDefault == 0 {
		g.QualityDefault = 85
	}

	return &Generator{
		disk:     g.Disk,
		host:     g.Host,
		fileName: fileName,
		width:    0,
		height:   0,
		quality:  g.QualityDefault,
		crop:     false,
		gravity:  "ce",
		format:   WEBP,
	}
}

func (g *Generator) Size(sizeName string) *Generator {
	panic("implement me")

	return g
}

func (g *Generator) Width(width int) *Generator {
	g.width = width

	return g
}

func (g *Generator) Height(height int) *Generator {
	g.height = height

	return g
}

func (g *Generator) Quality(quality int) *Generator {
	g.quality = quality

	return g
}

func (g *Generator) Crop() *Generator {
	g.crop = true

	return g
}

func (g *Generator) Gravity(gravity string) *Generator {
	g.gravity = gravity

	return g
}

func (g *Generator) Format(format Format) *Generator {
	g.format = format

	return g
}

func (g *Generator) Get() string {
	// Decide if using plain URL or encrypted URL
	// Decide if using disk or cloud bucket

	signature := "insecure"

	sizeAndCrop := fmt.Sprintf("rs:%s:%d:%d", "fill", g.width, g.height)
	gravity := fmt.Sprintf("g:%s", g.gravity)
	quality := fmt.Sprintf("q:%d", g.quality)
	parameters := fmt.Sprintf("%s/%s/%s/%s", sizeAndCrop, gravity, quality, "plain")

	path := fmt.Sprintf("%s://%s", g.disk, g.fileName)

	return fmt.Sprintf("%s/%s/%s/%s", g.host, signature, parameters, path)
}
