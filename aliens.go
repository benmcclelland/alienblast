// Copyright 2017 Google Inc. All rights reserved.
// Copyright 2017 Ben McClelland
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to writing, software distributed
// under the License is distributed on a "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
)

const (
	alienWidth  = 48
	alienHeight = 44
)

type aliens struct {
	mu sync.RWMutex

	texture *sdl.Texture
	explode *sdl.Texture
	speed   int32

	aliens []*alien
}

func newAliens(r *sdl.Renderer) (*aliens, error) {
	texture, err := img.LoadTexture(r, "res/imgs/si.png")
	if err != nil {
		return nil, fmt.Errorf("could not load alien image: %v", err)
	}
	explode, err := img.LoadTexture(r, "res/imgs/explode.png")
	if err != nil {
		return nil, fmt.Errorf("could not load image: %v", err)
	}

	as := &aliens{
		explode: explode,
		texture: texture,
		speed:   2,
	}

	go func() {
		for {
			as.mu.Lock()
			as.aliens = append(as.aliens, newAlien())
			as.mu.Unlock()
			time.Sleep(time.Second * 2)
		}
	}()

	return as, nil
}

func (as *aliens) paint(r *sdl.Renderer) error {
	as.mu.RLock()
	defer as.mu.RUnlock()

	for _, a := range as.aliens {
		if err := a.paint(r, as.texture, as.explode); err != nil {
			return err
		}
	}
	return nil
}

func (as *aliens) touch(s *ship) {
	as.mu.RLock()
	defer as.mu.RUnlock()
	for _, a := range as.aliens {
		a.touch(s)
	}
}

func (as *aliens) splode(b *blast) {
	as.mu.RLock()
	defer as.mu.RUnlock()
	for _, a := range as.aliens {
		a.splode(b)
	}
}

func (as *aliens) restart() {
	as.mu.Lock()
	defer as.mu.Unlock()

	as.aliens = nil
}

func (as *aliens) update() {
	as.mu.Lock()
	defer as.mu.Unlock()

	var rem []*alien
	for _, a := range as.aliens {
		a.mu.Lock()
		if !a.dead {
			a.x -= as.speed
		} else {
			a.deadtick++
			if a.deadtick > 50 {
				a.remove = true
			}
		}
		a.mu.Unlock()
		if a.x+a.w > 0 && !a.remove {
			rem = append(rem, a)
		}
	}
	as.aliens = rem
}

func (as *aliens) destroy() {
	as.mu.Lock()
	defer as.mu.Unlock()

	as.texture.Destroy()
}

type alien struct {
	mu sync.RWMutex

	x        int32
	y        int32
	h        int32
	w        int32
	xloc     int32
	yloc     int32
	dead     bool
	remove   bool
	deadtick int32
}

func newAlien() *alien {
	return &alien{
		x:    800,
		y:    int32(rand.Intn(500)),
		xloc: int32(rand.Intn(4)),
		yloc: int32(rand.Intn(6)),
		h:    alienWidth,
		w:    alienHeight,
	}
}

func (a *alien) touch(s *ship) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	s.touch(a)
}

func (a *alien) splode(b *blast) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if b.x > a.x+a.w { // blast too far right
		return
	}
	if b.x+b.w < a.x { // blast too far left
		return
	}
	if b.y+b.h < a.y { // blast too high
		return
	}
	if b.y > a.y+a.h { // bast too low
		return
	}
	a.dead = true
}

func (a *alien) paint(r *sdl.Renderer, texture *sdl.Texture, explode *sdl.Texture) error {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.dead {
		rect := &sdl.Rect{X: a.x, Y: a.y, W: a.w, H: a.w}
		if err := r.CopyEx(explode, nil, rect, 0, nil, sdl.FLIP_NONE); err != nil {
			return fmt.Errorf("could not copy texture: %v", err)
		}
	} else {
		rect := &sdl.Rect{X: a.x, Y: a.y, W: alienWidth, H: alienHeight}
		alien := &sdl.Rect{X: a.yloc * alienWidth, Y: a.xloc * alienHeight, W: a.w, H: a.h}
		if err := r.CopyEx(texture, alien, rect, 0, nil, sdl.FLIP_NONE); err != nil {
			return fmt.Errorf("could not copy background: %v", err)
		}
	}
	return nil
}
