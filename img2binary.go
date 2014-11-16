package main

import (
 "image/png"
 "image/color"
 "image"
 "os"
 "fmt"
 "math"
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


func main() {
    reader, err := os.Open("test.png")
    img, err := png.Decode(reader)
    if err != nil {
        fmt.Printf("%s\n", err)
    }
    
    reader.Close()
    
    bounds := img.Bounds()
    bwImage := image.NewRGBA(image.Rect(0, 0, bounds.Max.Y, bounds.Max.X))
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            col := img.At(x, y)
            r,g,b,a := col.RGBA()
            _ = a
            grayscale := rgb2lineargrayscale(r,g,b)
            binaryColor := grayscale2blackandwhite(grayscale)
            bwImage.Set(x,y, binaryColor)
        }
    }
    
    toimg, _ := os.Create("bwtest.png")
    defer toimg.Close()

    png.Encode(toimg, bwImage)
    
}
