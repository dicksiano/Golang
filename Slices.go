package main

import "golang.org/x/tour/pic"

func Pic(dx, dy int) [][]uint8 {
  pic := make([][]uint8, dy) // columns
  for i := range pic {
    pic[i] = make([]uint8, dx) // lines
    for j := range pic[i] {
      pic[i][j] = uint8(i+j)
	    //pic[i][j] = uint8(i*j)
    }
  }
  return pic
}
func main() {
	pic.Show(Pic)
}
