/*
rasterization code extracted and converted to go from https://github.com/petrkutalek/png2pos/blob/master/png2pos.c
*/

package escpos

const (
	GS8L_MAX_Y = 1662
)

func (e *Escpos) Raster(width, height, bytesWidth int, img_bw []byte) {
	flushCmd := []byte{
		/* GS ( L, Print the graphics data in the print buffer,
		   p. 241 Moves print position to the left side of the
		   print area after printing of graphics data is
		   completed */
		0x1d, 0x28, 0x4c, 0x02, 0x00, 0x30,
		/* Fn 50 */
		0x32,
	}

	for l := 0; l < height; {
		n_lines := GS8L_MAX_Y
		if n_lines > height-l {
			n_lines = height - l
		}

		f112_p := 10 + n_lines*bytesWidth
		storeCmd := []byte{
			/* GS 8 L, Store the graphics data in the print buffer
			   (raster format), p. 252 */
			0x1d, 0x38, 0x4c,
			/* p1 p2 p3 p4 */
			byte(f112_p), byte(f112_p >> 8),
			byte(f112_p >> 16), byte(f112_p >> 24),
			/* Function 112 */
			0x30, 0x70, 0x30,
			/* bx by, zoom */
			0x01, 0x01,
			/* c, single-color printing model */
			0x31,
			/* xl, xh, number of dots in the horizontal direction */
			byte(width), byte(width >> 8),
			/* yl, yh, number of dots in the vertical direction */
			byte(n_lines), byte(n_lines >> 8),
		}

		e.WriteRaw(storeCmd)
		e.WriteRaw(img_bw[l*bytesWidth : (l+n_lines)*bytesWidth])
		e.WriteRaw(flushCmd)

		l += n_lines
	}
}
