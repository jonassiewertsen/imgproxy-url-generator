# Imgproxy URL Generator

This is a simple tool to generate [imgproxy](https://docs.imgproxy.net/) URLs for images. 

## Example

```go
g := NewImgproxyUrlGenerator(Config{
    Host: "https://example.com",
    Disk: "s3://bucket",
    //Key:  "",
    //Salt: "",
})
filename := g.File("file.jpg").Width(500).Height(300).Quality(92).Crop().Get()
```
