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
	"log"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
)

type scene struct {
	r      *sdl.Renderer
	bg     *sdl.Texture
	ship   *ship
	aliens *aliens
	blast  *blast
}

func newScene(r *sdl.Renderer) (*scene, error) {
	bg, err := img.LoadTexture(r, "res/imgs/moon.png")
	if err != nil {
		return nil, fmt.Errorf("could not load background image: %v", err)
	}

	s, err := newShip(r)
	if err != nil {
		return nil, err
	}

	as, err := newAliens(r)
	if err != nil {
		return nil, err
	}

	return &scene{r: r, bg: bg, ship: s, aliens: as}, nil
}

func (s *scene) run(events <-chan sdl.Event, r *sdl.Renderer) <-chan error {
	errc := make(chan error)

	go func() {
		defer close(errc)
		tick := time.Tick(10 * time.Millisecond)
		for {
			select {
			case e := <-events:
				if done := s.handleEvent(e); done {
					return
				}
			case <-tick:
				s.update()

				if s.ship.isDead() {
					s.restart()
				}

				if err := s.paint(r); err != nil {
					errc <- err
				}
			}
		}
	}()

	return errc
}

func (s *scene) handleEvent(event sdl.Event) bool {
	switch event.(type) {
	case *sdl.QuitEvent:
		return true
	case *sdl.MouseButtonEvent:
		s.fire()
	case *sdl.KeyDownEvent:
		switch event.(*sdl.KeyDownEvent).Keysym.Sym {
		case sdl.GetKeyFromName("Space"):
			s.fire()
		case sdl.GetKeyFromName("Escape"):
			return true
		case sdl.GetKeyFromName("Left"):
			s.ship.left()
		case sdl.GetKeyFromName("Right"):
			s.ship.right()
		case sdl.GetKeyFromName("Up"):
			s.ship.up()
		case sdl.GetKeyFromName("Down"):
			s.ship.down()

		}
	case *sdl.MouseMotionEvent, *sdl.WindowEvent, *sdl.TouchFingerEvent, *sdl.CommonEvent, *sdl.KeyUpEvent, *sdl.TextInputEvent:
	default:
		log.Printf("unknown event %T", event)
	}
	return false
}

func (s *scene) update() {
	s.ship.update()
	s.aliens.update()
	if s.blast != nil {
		s.blast = s.blast.update()
	}
	s.aliens.touch(s.ship)
	if s.blast != nil {
		s.aliens.splode(s.blast)
	}
}

func (s *scene) restart() {
	s.ship.restart()
	s.aliens.restart()
}

func (s *scene) fire() {
	if s.blast != nil {
		return
	}
	if s.ship.dead {
		return
	}
	var err error
	s.blast, err = newBlast(s.r, s.ship.x+shipHeight, s.ship.y+shipWidth/2)
	if err != nil {
		fmt.Println("new blast err", err)
	}
}

func (s *scene) paint(r *sdl.Renderer) error {
	r.Clear()
	rect := &sdl.Rect{X: 0, Y: 0, W: 800, H: 600}
	moonrect := &sdl.Rect{X: 100, Y: 0, W: 600, H: 600}
	if err := r.DrawRect(rect); err != nil {
		return fmt.Errorf("Could not draw background rect: %v", err)
	}
	if err := r.SetDrawColor(uint8(0), uint8(0), uint8(0), uint8(255)); err != nil {
		return fmt.Errorf("Could not set draw color: %v", err)
	}
	if err := r.FillRect(rect); err != nil {
		return fmt.Errorf("Could not fill rect: %v", err)
	}
	if err := r.CopyEx(s.bg, nil, moonrect, 0, nil, sdl.FLIP_NONE); err != nil {
		return fmt.Errorf("could not copy background: %v", err)
	}
	if err := s.ship.paint(r); err != nil {
		return err
	}
	if err := s.aliens.paint(r); err != nil {
		return err
	}
	if s.blast != nil {
		if err := s.blast.paint(r); err != nil {
			return err
		}
	}
	r.Present()
	return nil
}

func (s *scene) destroy() {
	s.ship.destroy()
	s.aliens.destroy()
	if s.blast != nil {
		s.blast.destroy()
	}
}
