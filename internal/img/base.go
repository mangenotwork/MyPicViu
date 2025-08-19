package img

import "math"

// RGBToHSV 将 RGB 颜色转换为 HSV 颜色
func RGBToHSV(r, g, b uint8) (float64, float64, float64) {
	rNorm := float64(r) / 255.0
	gNorm := float64(g) / 255.0
	bNorm := float64(b) / 255.0
	maxVal := math.Max(rNorm, math.Max(gNorm, bNorm))
	minVal := math.Min(rNorm, math.Min(gNorm, bNorm))
	delta := maxVal - minVal

	var h, s, v float64
	v = maxVal

	if delta == 0 {
		h = 0
	} else {
		s = delta / maxVal
		if maxVal == rNorm {
			h = math.Mod((gNorm-bNorm)/delta, 6)
		} else if maxVal == gNorm {
			h = (bNorm-rNorm)/delta + 2
		} else {
			h = (rNorm-gNorm)/delta + 4
		}
		h *= 60
		if h < 0 {
			h += 360
		}
	}
	return h, s, v
}

// HSVToRGB 将 HSV 颜色转换为 RGB 颜色
func HSVToRGB(h, s, v float64) (uint8, uint8, uint8) {
	c := v * s
	hPrime := h / 60
	x := c * (1 - math.Abs(math.Mod(hPrime, 2)-1))
	var r1, g1, b1 float64
	switch {
	case 0 <= hPrime && hPrime < 1:
		r1 = c
		g1 = x
		b1 = 0
	case 1 <= hPrime && hPrime < 2:
		r1 = x
		g1 = c
		b1 = 0
	case 2 <= hPrime && hPrime < 3:
		r1 = 0
		g1 = c
		b1 = x
	case 3 <= hPrime && hPrime < 4:
		r1 = 0
		g1 = x
		b1 = c
	case 4 <= hPrime && hPrime < 5:
		r1 = x
		g1 = 0
		b1 = c
	case 5 <= hPrime && hPrime < 6:
		r1 = c
		g1 = 0
		b1 = x
	}
	m := v - c
	r := uint8((r1 + m) * 255)
	g := uint8((g1 + m) * 255)
	b := uint8((b1 + m) * 255)
	return r, g, b
}
