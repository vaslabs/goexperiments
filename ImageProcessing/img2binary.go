package main

import (
    "image/png"
    "image/color"
    "image"
    "runtime"
    "os"
    "fmt"
    "math"
    "sync"
)

func rgbV2linear(v float64) float64 {
    if v < 0.04045 {
        return v/12.92
    }
    return math.Pow((v+0.055)/1.055, 2.4)
}

func rgb2lineargrayscale(R,G,B uint32) float64 {
    R_linear := rgbV2linear(float64(R)/65535.0)
    G_linear := rgbV2linear(float64(G)/65535.0)
    B_linear := rgbV2linear(float64(B)/65535.0)
    return 0.299 * R_linear + 0.587 * G_linear + 0.114 * B_linear
}

func grayscale2blackandwhite(v float64) color.Color {
    if (v < 0.5) {
        return color.Black
    } else {
        return color.White
    }
}

func saveImage(fn string, bwImage image.Image) {
    
    toimg, _ := os.Create(fn)
    defer toimg.Close()

    png.Encode(toimg, bwImage)
}

func convertImagePart(img image.Image, bwImage *image.RGBA, startX, endX, startY, endY int) {
    for y := startY; y < endY ; y++ {
        for x := startX; x < endX; x++ {
            col := img.At(x, y)
            r,g,b,a := col.RGBA()
            _ = a
            grayscale := rgb2lineargrayscale(r,g,b)
            binaryColor := grayscale2blackandwhite(grayscale)
            bwImage.Set(x,y, binaryColor)
        }
    }
}


func img2binary(img image.Image) image.Image {
    var wg sync.WaitGroup
    bounds := img.Bounds()
    bwImage := image.NewRGBA(image.Rect(0, 0 ,bounds.Max.X, bounds.Max.Y))
    startY := 0
    endY := bounds.Max.Y / 2
    
    startYOther := endY;
    endYOther := bounds.Max.Y;
    
    wg.Add(2)
    
    go func() {
        defer wg.Done()
        convertImagePart(img, bwImage, 0, bounds.Max.X, startY, endY)
    }()
    go func() {
        defer wg.Done()
        convertImagePart(img, bwImage, 0, bounds.Max.X, startYOther, endYOther)
    }()
    wg.Wait()
    
    return bwImage
}

func main() {
    fmt.Println("CPUs: ", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())
    file2convert := os.Args[1]
    save2 := os.Args[2]
    reader, err := os.Open(file2convert)
    img, err := png.Decode(reader)
    if err != nil {
        fmt.Printf("%s\n", err)
    }
    
    reader.Close()
    
    newImg := img2binary(img)
    saveImage(save2, newImg)
    
}
