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
	"sync"

	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
)

const (
	moveSpeed  = 20
	shipWidth  = 44
	shipHeight = 79
)

type ship struct {
	mu sync.RWMutex

	texture *sdl.Texture
	explode *sdl.Texture

	x, y     int32
	w, h     int32
	speed    float64
	dead     bool
	deadtick int32
}

func newShip(r *sdl.Renderer) (*ship, error) {
	texture, err := img.LoadTexture(r, "res/imgs/rocket.png")
	if err != nil {
		return nil, fmt.Errorf("could not load image: %v", err)
	}
	explode, err := img.LoadTexture(r, "res/imgs/explode.png")
	if err != nil {
		return nil, fmt.Errorf("could not load image: %v", err)
	}
	return &ship{texture: texture, explode: explode, x: 10, y: 500, w: shipWidth, h: shipHeight}, nil
}

func (s *ship) update() {
	s.mu.Lock()
	defer s.mu.Unlock()
}

func (s *ship) paint(r *sdl.Renderer) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.dead {
		rect := &sdl.Rect{X: s.x, Y: s.y, W: s.w, H: s.w}
		if err := r.CopyEx(s.explode, nil, rect, 0, nil, sdl.FLIP_NONE); err != nil {
			return fmt.Errorf("could not copy texture: %v", err)
		}
	} else {
		rect := &sdl.Rect{X: s.x, Y: s.y, W: s.w, H: s.h}
		if err := r.CopyEx(s.texture, nil, rect, 90, nil, sdl.FLIP_NONE); err != nil {
			return fmt.Errorf("could not copy texture: %v", err)
		}
	}
	return nil
}

func (s *ship) restart() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.y = 500
	s.x = 10
	s.speed = 0
	s.dead = false
	s.deadtick = 0
}

func (s *ship) destroy() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.texture.Destroy()
	s.explode.Destroy()
}

func (s *ship) isDead() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.dead {
		s.deadtick++
		if s.deadtick > 100 {
			return s.dead
		}
	}
	return false
}

func (s *ship) left() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.dead {
		s.x -= moveSpeed
		if s.x < 0 {
			s.x = 0
		}
	}
}

func (s *ship) right() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.dead {
		s.x += moveSpeed
		if s.x > 740 {
			s.x = 740
		}
	}
}

func (s *ship) up() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.dead {
		s.y -= moveSpeed
		if s.y < 0 {
			s.y = 0
		}
	}
}

func (s *ship) down() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.dead {
		s.y += moveSpeed
		if s.y > 550 {
			s.y = 550
		}
	}
}

func (s *ship) touch(a *alien) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if a.x > s.x+s.h { // alien too far right
		return
	}
	if a.x+a.w < s.x { // alien too far left
		return
	}
	if a.y+(a.h/2) < s.y+4 { // alien too high
		return
	}
	if a.y > s.y+s.w { // alien too low
		return
	}
	s.dead = true
}
