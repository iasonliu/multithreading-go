package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const numberOfThreads int = 8

type Point2D struct {
	x int
	y int
}

var (
	r         = regexp.MustCompile(`\((\d*),(\d*)\)`)
	waitGroup = sync.WaitGroup{}
)

func findArea(inputChannel chan string) {
	for pointsStr := range inputChannel {
		var points []Point2D
		for _, p := range r.FindAllStringSubmatch(pointsStr, -1) {
			x, _ := strconv.Atoi(p[1])
			y, _ := strconv.Atoi(p[2])
			points = append(points, Point2D{x, y})
		}
		area := 0.0
		for i := 0; i < len(points); i++ {
			a, b := points[i], points[(i+1)%len(points)]
			area += float64(a.x*b.y) - float64(a.y*b.x)
		}

		fmt.Println(math.Abs(area) / 2.0)
	}
	waitGroup.Done()
}
func main() {
	absPath, _ := filepath.Abs("../")
	data, _ := ioutil.ReadFile(filepath.Join(absPath, "polygons.txt"))
	text := string(data)

	inputChannel := make(chan string, 1000)
	for i := 0; i < numberOfThreads; i++ {
		go findArea(inputChannel)
	}
	waitGroup.Add(numberOfThreads)
	start := time.Now()
	for _, line := range strings.Split(text, "\n") {
		// line := "(4,10),(12,8),(10,3),(2,2),(7,5)"
		inputChannel <- line
	}
	close(inputChannel)
	waitGroup.Wait()
	elapsed := time.Since(start)
	fmt.Printf("Processing took %s \n", elapsed)
}
