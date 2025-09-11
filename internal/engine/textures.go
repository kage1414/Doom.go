package engine

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) initTextures() {
	if g.wallTex != nil {
		return
	}
	const tw, th = 256, 256
	img := image.NewRGBA(image.Rect(0, 0, tw, th))
	makeRuggedRock(img)
	g.wallTex = ebiten.NewImageFromImage(img)
	g.texW, g.texH = tw, th
}

// makeRuggedRock fills an image with a dark, menacing rock-like texture using value noise,
// dents, and rusty stains—no external assets.
func makeRuggedRock(dst *image.RGBA) {
	w, h := dst.Bounds().Dx(), dst.Bounds().Dy()
	// base colors
	baseA := color.RGBA{24, 24, 28, 255}
	baseB := color.RGBA{36, 36, 42, 255}
	draw.Draw(dst, dst.Bounds(), &image.Uniform{baseA}, image.Point{}, draw.Src)
	// subtle tiles/strata break-up
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if ((x>>4)^(y>>4))&1 == 0 {
				dst.SetRGBA(x, y, baseB)
			}
		}
	}
	rng := rand.New(rand.NewSource(1337))
	// value noise (3 octaves)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			n := 0.0
			amp := 1.0
			freq := 1.0 / 16.0
			for o := 0; o < 3; o++ {
				n += amp * valueNoise2D(float64(x)*freq, float64(y)*freq, rng)
				amp *= 0.5
				freq *= 2.0
			}
			n = (n + 1) * 0.5 // 0..1
			// dark rock with slight blue tint
			r := uint8(clamp01(n*0.35+0.15) * 255)
			g := uint8(clamp01(n*0.35+0.16) * 255)
			b := uint8(clamp01(n*0.45+0.18) * 255)
			// embed into existing pixel to keep some strata
			orig := dst.RGBAAt(x, y)
			dst.SetRGBA(x, y, color.RGBA{
				R: uint8((int(orig.R)*3 + int(r)) / 4),
				G: uint8((int(orig.G)*3 + int(g)) / 4),
				B: uint8((int(orig.B)*3 + int(b)) / 4),
				A: 255,
			})
		}
	}
	// cracks: random thin dark lines
	for i := 0; i < 60; i++ {
		x := rng.Intn(w)
		y := rng.Intn(h)
		len := rng.Intn(w/2) + w/6
		ang := rng.Float64() * math.Pi * 2
		dx := math.Cos(ang)
		dy := math.Sin(ang)
		for t := 0; t < len; t++ {
			px := int(float64(x) + float64(t)*dx)
			py := int(float64(y) + float64(t)*dy)
			if px < 0 || py < 0 || px >= w || py >= h {
				break
			}
			c := dst.RGBAAt(px, py)
			dark := uint8(float64(c.R) * 0.6)
			dst.SetRGBA(px, py, color.RGBA{dark, dark, dark, 255})
			// slight thickness
			if px+1 < w {
				c2 := dst.RGBAAt(px+1, py)
				dark2 := uint8(float64(c2.R) * 0.75)
				dst.SetRGBA(px+1, py, color.RGBA{dark2, dark2, dark2, 255})
			}
		}
	}
	// stains: rusty red blots
	for i := 0; i < 10; i++ {
		cx := rng.Intn(w)
		cy := rng.Intn(h)
		rad := rng.Intn(w/6) + w/12
		for y := -rad; y <= rad; y++ {
			for x := -rad; x <= rad; x++ {
				px := cx + x
				py := cy + y
				if px < 0 || py < 0 || px >= w || py >= h {
					continue
				}
				d := math.Hypot(float64(x), float64(y)) / float64(rad)
				if d > 1 {
					continue
				}
				s := clamp01(1 - d*d)
				base := dst.RGBAAt(px, py)
				r := uint8(clamp01(float64(base.R)/255.0*0.9+0.25*s) * 255)
				g := uint8(clamp01(float64(base.G)/255.0*0.8+0.05*s) * 255)
				b := uint8(clamp01(float64(base.B)/255.0*0.7+0.00*s) * 255)
				dst.SetRGBA(px, py, color.RGBA{r, g, b, 255})
			}
		}
	}
}

func valueNoise2D(x, y float64, rng *rand.Rand) float64 {
	// simple hash-based gradientless noise
	x0 := math.Floor(x)
	y0 := math.Floor(y)
	tx := x - x0
	ty := y - y0
	v00 := hash01(int(x0), int(y0), rng)
	v10 := hash01(int(x0+1), int(y0), rng)
	v01 := hash01(int(x0), int(y0+1), rng)
	v11 := hash01(int(x0+1), int(y0+1), rng)
	sx := smoothstep(tx)
	sy := smoothstep(ty)
	a := lerpF(v00, v10, sx)
	b := lerpF(v01, v11, sx)
	return lerpF(a, b, sy)*2 - 1 // -1..1
}

func hash01(x, y int, rng *rand.Rand) float64 {
	// deterministic integer hash (avoid rng use to keep tiling consistent)
	h := uint32(x*73856093 ^ y*19349663)
	h ^= h >> 13
	h *= 1274126177
	return float64(h&0xFFFF) / 65535.0
}

func lerpF(a, b, t float64) float64 { return a + (b-a)*t }
func smoothstep(t float64) float64  { return t * t * (3 - 2*t) }
