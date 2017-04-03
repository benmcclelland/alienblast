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
	blastSpeed  = 10
	blastWidth  = 44
	blastHeight = 10
)

type blast struct {
	mu sync.RWMutex

	texture *sdl.Texture

	x, y int32
	w, h int32
	time int32
}

func newBlast(r *sdl.Renderer, x, y int32) (*blast, error) {
	path := fmt.Sprintf("res/imgs/blast.png")
	texture, err := img.LoadTexture(r, path)
	if err != nil {
		return nil, fmt.Errorf("could not load image: %v", err)
	}
	return &blast{texture: texture, x: x, y: y, w: blastWidth, h: blastHeight}, nil
}

func (b *blast) update() *blast {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.time++
	b.x += blastSpeed
	if b.x > 800 {
		b.texture.Destroy()
		return nil
	}
	return b
}

func (b *blast) paint(r *sdl.Renderer) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	rect := &sdl.Rect{X: b.x, Y: b.y, W: blastWidth, H: blastHeight}
	i := b.time / 10 % 11
	blast := &sdl.Rect{X: 0, Y: i * 128, W: 512, H: 128}

	if err := r.CopyEx(b.texture, blast, rect, 0, nil, sdl.FLIP_NONE); err != nil {
		return fmt.Errorf("could not copy texture: %v", err)
	}
	return nil
}

func (b *blast) destroy() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.texture.Destroy()
}
