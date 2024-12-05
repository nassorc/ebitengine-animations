package animations

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Spritesheet struct{ Rows, Cols, TileWidth, TileHeight int }

func NewSpritesheet(row, col, tw, th int) *Spritesheet {
	return &Spritesheet{
		Rows:       row,
		Cols:       col,
		TileWidth:  tw,
		TileHeight: th,
	}
}

func (sheet *Spritesheet) Cell(cell int) []int {
	return []int{cell}
}

func (sheet *Spritesheet) Cells(cells []int) []int {
	for _, cell := range cells {
		if cell < 0 || cell >= (sheet.Cols*sheet.Rows) {
			panic("out of bounds")
		}
	}
	return cells
}

func (sheet *Spritesheet) Row(r int) []int {
	var clip = []int{}
	var rowStart = r * sheet.Cols

	for i := rowStart; i < (rowStart + sheet.Cols); i++ {
		clip = append(clip, i)
	}

	return clip
}

func (sheet *Spritesheet) Col(c int) []int {
	var clip = []int{}
	var width = sheet.Cols

	for i := range sheet.Rows {
		clip = append(clip, c*(i*width))
	}

	return clip
}

type Animation struct {
	Name    string
	Frames  []int
	Layer   int
	OffsetX int
	OffsetY int
  FPS     int
}

func NewAnimation(name string, frames []int) *Animation {
	return &Animation{name, frames, 0, 0, 0, 1}
}

func (a *Animation) SetOffset(ox, oy int) {
	a.OffsetX = ox
	a.OffsetY = oy
}

type AnimationMap struct {
	State           string
	Frame           int
	Texture         *ebiten.Image
	Animations      map[string]*Animation
	Spritesheet     *Spritesheet
	Tick            float64
	CurrentFrameIdx int
}

func NewAnimationMap(txr *ebiten.Image, sheet *Spritesheet) *AnimationMap {
	return &AnimationMap{
		Texture:     txr,
		Animations:  make(map[string]*Animation),
		Spritesheet: sheet,
	}
}

type AnimationOpt func(anim *Animation)

func WithFPS(fps int) AnimationOpt {
	return func(anim *Animation) {
    anim.FPS = fps
  }
}

func WithOffset(x, y int) AnimationOpt {
	return func(anim *Animation) {
		anim.OffsetX = x
		anim.OffsetY = y
	}
}

func WithLayer(l int) AnimationOpt {
	return func(anim *Animation) {
		anim.Layer = l
	}
}

func (a *AnimationMap) CreateClip(name string, frames []int, opts ...AnimationOpt) {
	var anim = NewAnimation(name, frames)
	a.Animations[name] = anim

	for _, f := range opts {
		f(anim)
	}
}

func (a *AnimationMap) Switch(name string) {
	if _, ok := a.Animations[name]; !ok {
		return
	}

	a.State = name
}

func (a *AnimationMap) IsPlaying(name string) bool {
	return a.State == name
}

func (a *AnimationMap) Update() {
  var animation = a.Animations[a.State]
	a.CurrentFrameIdx = int(a.Tick) % len(animation.Frames)
  a.Tick += (float64(animation.FPS)/float64(ebiten.TPS()))
}

func (a *AnimationMap) Animation() *Animation {
	return a.Animations[a.State]
}

func (a *AnimationMap) Sprite() *ebiten.Image {
	var animation = a.Animations[a.State]
	var frameIdx = animation.Frames[a.CurrentFrameIdx]
	var sheet = a.Spritesheet
	var row = frameIdx / sheet.Cols
	var col = frameIdx % sheet.Cols
	var x = (col * sheet.TileWidth)
	var y = (row * sheet.TileHeight)

	return a.Texture.SubImage(image.Rect(x, y, x+sheet.TileWidth, y+sheet.TileHeight)).(*ebiten.Image)
}
