package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

type Config struct {
	Salt           string
	saltBin        []byte
	Key            string
	keyBin         []byte
	Host           string
	Disk           string
	QualityDefault int
	EncodeUrl      bool
}

type Repository struct {
	Config                Config
	qualityDefault        int
	plainUrl              string
	encryptedAndSignedUrl string
}

type Generator struct {
	Config   Config
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

func NewImgproxyUrlGenerator(config Config) (*Repository, error) {
	if config.Host == "" || config.Disk == "" {
		return nil, errors.New("host and disk are required")
	}

	// If no salt or key is provided, we won't encrypt the URL
	if config.Salt == "" || config.Key == "" {
		return &Repository{
			Config:         config,
			qualityDefault: 0, // Defaults to img proxy default if 0
		}, nil
	}

	var err error
	if config.keyBin, err = hex.DecodeString(config.Key); err != nil {
		return nil, fmt.Errorf("error decoding key: %w", err)
	}

	if config.saltBin, err = hex.DecodeString(config.Salt); err != nil {
		return nil, fmt.Errorf("error decoding salt: %w", err)
	}

	return &Repository{
		Config:         config,
		qualityDefault: 0, // Defaults to img proxy default if 0
	}, nil
}

func (g *Repository) File(fileName string) *Generator {
	return &Generator{
		Config:   g.Config,
		disk:     g.Config.Disk,
		host:     g.Config.Host,
		fileName: fileName,
		width:    0,
		height:   0,
		quality:  g.Config.QualityDefault,
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
	path := g.generatePath()

	sizeAndCrop := fmt.Sprintf("rs:%s:%d:%d", "fill", g.width, g.height)
	gravity := fmt.Sprintf("g:%s", g.gravity)
	quality := fmt.Sprintf("q:%d", g.quality)
	parameters := fmt.Sprintf("%s/%s/%s", sizeAndCrop, gravity, quality)

	fileUrl := strings.Join([]string{string(g.Config.saltBin), parameters, path}, "/")
	fileUrl = strings.Join([]string{fileUrl, string(g.format)}, ".") // Add the format to the URL

	signature := g.generateSignature(fileUrl)

	return fmt.Sprintf("%s/%s/%s/%s.%s", g.host, signature, parameters, path, g.format)
}

func (g *Generator) generateSignature(fileUrl string) string {
	if g.Config.Salt == "" && g.Config.Key == "" {
		return "insecure" // Has to be defined if no signature is provided
	}

	mac := hmac.New(sha256.New, g.Config.keyBin)
	mac.Write([]byte(fileUrl))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

}

func (g *Generator) generatePath() string {
	path := fmt.Sprintf("%s/%s", g.disk, g.fileName)

	if g.Config.EncodeUrl {
		return encodeFilePath(path)
	}

	return "plain/" + path
}

func encodeFilePath(filePath string) string {
	// Encode the file path to Base64
	encoded := base64.StdEncoding.EncodeToString([]byte(filePath))

	// Replace characters as per the custom mapping
	encoded = strings.NewReplacer("+", "-", "/", "_", "=", "").Replace(encoded)

	// Add a slash ("/") after every 16 characters
	var result strings.Builder
	for i, r := range encoded {
		result.WriteRune(r)
		if (i+1)%16 == 0 { // Add '/' after every 16th character
			result.WriteRune('/')
		}
	}

	return result.String()
}
