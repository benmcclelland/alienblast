# Alien Blaster

Alien Blaster uses a lot of boilerplate code from Francesc's Flappy Gopher
JustForFunc demo.  You can find it [here](https://github.com/campoy/flappy-gopher).

## Installation

You need to install first SDL2 and the SDL2 bindings for Go. To do so follow the instructions [here](https://github.com/veandco/go-sdl2).
It is quite easy to install on basically any platform.

You will also need to install [pkg-config](https://en.wikipedia.org/wiki/Pkg-config).

For mac:
```
brew install sdl2{,_image,_ttf,_mixer}
brew install pkg-config
```

After that you should be able to simply run:

    go get github.com/benmcclelland/alienblast

And run the binary generated in `$GOPATH/bin`.

## Images, fonts, and licenses

All the images used in this game are obtained from [Clipart](https://openclipart.org/).

This project is licensed under Apache v2 license. Read more in the [LICENSE](LICENSE) file.
