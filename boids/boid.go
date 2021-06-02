package main

import (
	"math"
	"math/rand"
	"time"
)

type Boid struct {
	position Vector2D
	velocity Vector2D
	id       int
}

func (b *Boid) calcAcceleration() Vector2D {
	uppper, lower := b.position.AddV(viewRadius), b.position.AddV(-viewRadius)
	avgPosition, avgVelocity, separation := Vector2D{0, 0}, Vector2D{0, 0}, Vector2D{0, 0}

	// count of boids in view radius
	count := 0.0

	rWlock.RLock()
	for i := math.Max(lower.x, 0); i <= math.Min(uppper.x, screenWidth); i++ {
		for j := math.Max(lower.y, 0); j <= math.Min(uppper.y, screenHeight); j++ {
			if otherBoidId := boidMap[int(i)][int(j)]; otherBoidId != -1 && otherBoidId != b.id {
				if dist := boids[otherBoidId].position.Distance(b.position); dist < viewRadius {
					count++
					avgVelocity = avgVelocity.Add(boids[otherBoidId].velocity)
					avgPosition = avgPosition.Add(boids[otherBoidId].position)
					separation = separation.Add(b.position.Subtract(boids[otherBoidId].position).DivisionV(dist))
				}

			}
		}
	}
	rWlock.RUnlock()
	acceleration := Vector2D{0, 0}
	if count > 0 {
		avgPosition, avgVelocity = avgPosition.DivisionV(count), avgVelocity.DivisionV(count)
		accelAlignment := avgVelocity.Subtract(b.velocity).MultiplyV(adjRate)
		accelCohesion := avgPosition.Subtract(b.position).MultiplyV(adjRate)
		accelSeparation := separation.MultiplyV(adjRate)
		acceleration = acceleration.Add(accelAlignment).Add(accelCohesion).Add(accelSeparation)
	}
	return acceleration
}

func (b *Boid) moveOne() {
	//calcAcceleration() has lock in, so must move out of this lock
	acceleration := b.calcAcceleration()
	rWlock.Lock()
	// limit 1 pixel, don't jump more than 1 pixel. more smooth for simulation
	b.velocity = b.velocity.Add(acceleration).limit(-1, 1)
	// update new position for boid after move
	boidMap[int(b.position.x)][int(b.position.y)] = -1
	b.position = b.position.Add(b.velocity)
	boidMap[int(b.position.x)][int(b.position.y)] = b.id
	next := b.position.Add(b.velocity)
	if next.x >= screenWidth || next.x < 0 {
		b.velocity = Vector2D{-b.velocity.x, b.velocity.y}
	}
	if next.y >= screenHeight || next.y < 0 {
		b.velocity = Vector2D{b.velocity.x, -b.velocity.y}
	}
	rWlock.Unlock()
}

func (b *Boid) start() {
	for {
		b.moveOne()
		time.Sleep(5 * time.Millisecond)
	}
}

func createBoid(bid int) {
	b := Boid{
		position: Vector2D{rand.Float64() * screenWidth, rand.Float64() * screenHeight},
		velocity: Vector2D{(rand.Float64() * 2) - 1.0, (rand.Float64() * 2) - 1.0},
		id:       bid,
	}
	boids[bid] = &b
	boidMap[int(b.position.x)][int(b.position.y)] = b.id
	go b.start()
}
