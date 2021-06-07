package deadlock

import (
	"time"

	"github.com/iasonliu/multithreading/deadlocks_train/common"
)

func MoveTrain(train *common.Train, distance int, crossings []*common.Crossing) {
	for train.Front < distance {
		train.Front += 1
		for _, crossing := range crossings {
			if train.Front == crossing.Position {
				crossing.Intersection.Mutex.Lock()
				crossing.Intersection.LockedBy = train.Id
			}
			back := train.Front - train.TrainLength
			if back == crossing.Position {
				crossing.Intersection.LockedBy = -1
				crossing.Intersection.Mutex.Unlock()
			}
		}
		time.Sleep(30 * time.Millisecond)
	}
}
