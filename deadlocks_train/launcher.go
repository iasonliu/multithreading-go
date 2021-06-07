package main

import (
	"log"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/iasonliu/multithreading/deadlocks_train/arbitrator"
	"github.com/iasonliu/multithreading/deadlocks_train/common"
)

var (
	trains        [4]*common.Train
	intersections [4]*common.Intersection
)

const trainLength = 70

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	DrawTracks(screen)
	DrawIntersections(screen)
	DrawTrains(screen)
}

func (g *Game) Layout(_, _ int) (w, h int) {
	return 320, 320
}

func main() {
	for i := 0; i < 4; i++ {
		trains[i] = &common.Train{Id: i, TrainLength: trainLength, Front: 0}
	}

	for i := 0; i < 4; i++ {
		intersections[i] = &common.Intersection{Id: i, Mutex: sync.Mutex{}, LockedBy: -1}
	}

	go arbitrator.MoveTrain(trains[0], 300, []*common.Crossing{{Position: 125, Intersection: intersections[0]},
		{Position: 175, Intersection: intersections[1]}})

	go arbitrator.MoveTrain(trains[1], 300, []*common.Crossing{{Position: 125, Intersection: intersections[1]},
		{Position: 175, Intersection: intersections[2]}})

	go arbitrator.MoveTrain(trains[2], 300, []*common.Crossing{{Position: 125, Intersection: intersections[2]},
		{Position: 175, Intersection: intersections[3]}})

	go arbitrator.MoveTrain(trains[3], 300, []*common.Crossing{{Position: 125, Intersection: intersections[3]},
		{Position: 175, Intersection: intersections[0]}})

	ebiten.SetWindowSize(320*3, 320*3)
	ebiten.SetWindowTitle("Trains in a box")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
