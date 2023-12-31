//
//    Lollipops diagram generation framework for genetic variations.
//    Copyright (C) 2015 Jeremy Jay <jeremy@pbnjay.com>
//
//    This program is free software: you can redistribute it and/or modify
//    it under the terms of the GNU General Public License as published by
//    the Free Software Foundation, either version 3 of the License, or
//    (at your option) any later version.
//
//    This program is distributed in the hope that it will be useful,
//    but WITHOUT ANY WARRANTY; without even the implied warranty of
//    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//    GNU General Public License for more details.
//
//    You should have received a copy of the GNU General Public License
//    along with this program.  If not, see <http://www.gnu.org/licenses/>.

package drawing

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

var (
	FontName string
	theFont  *truetype.Font
)

var defaultFontURL = "https://github.com/googlefonts/opensans/raw/main/fonts/ttf/OpenSans-Regular.ttf"

// we try to have sane defaults wrt font usage
//
// 1) auto-load Arial if found as the default font.
// 2) allow users to set a different font if desired
//

func LoadDefaultFont() error {
	// try to find Arial in the most common locations
	commonPaths := []string{
		// OS X path
		"/Library/Fonts/Arial.ttf",

		// Windows paths
		"C:\\WINDOWS\\Fonts\\arial.ttf",
		"C:\\WINNT\\Fonts\\arial.ttf",

		// Ubuntu with multiverse msttcorefonts package
		"/usr/share/fonts/truetype/msttcorefonts/arial.ttf",
	}
	for _, path := range commonPaths {
		err := LoadFont("Arial", path)
		if err == nil {
			return nil
		}
	}
	err := LoadFont("OpenSans", "OpenSans-Regular.ttf")
	if err == nil {
		return nil
	}

	resp, err := http.Get(defaultFontURL)
	if err == nil {
		fmt.Fprintln(os.Stderr, "Downloading Font: OpenSans-Regular.ttf (Apache Licensed)")
		fmt.Fprintln(os.Stderr, "Learn more about usage at https://github.com/googlefonts/opensans")
		fmt.Fprint(os.Stderr, "Or provide a different truetype font with -f fontname.ttf\n\n")

		f, err := os.Create("OpenSans-Regular.ttf")
		if err == nil {
			n, err := io.Copy(f, resp.Body)
			if n == 0 || err != nil {
				return fmt.Errorf("failed to download a font to use")
			}
		}
		f.Close()
		resp.Body.Close()
		err = LoadFont("OpenSans", "OpenSans-Regular.ttf")
		if err == nil {
			return nil
		}
	}

	return fmt.Errorf("unable to find Arial.ttf")
}

func LoadFont(name, path string) error {
	fontBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	theFont, err = truetype.Parse(fontBytes)
	if err != nil {
		return err
	}
	FontName = name
	return nil
}

// MeasureFont returns the pixel width of the string s at font size sz.
// It tries to use system Arial font if possible, but falls back to a
// conservative ballpark estimate otherwise.
func (x *Settings) MeasureFont(s string, sz int) int {
	// use actual TTF font metrics if available
	if theFont != nil {
		myFace := truetype.NewFace(theFont, &truetype.Options{
			Size: float64(sz),
			DPI:  float64(x.dpi),
		})
		d := &font.Drawer{Face: myFace}
		w := d.MeasureString(s)

		// convert from 26.6 fixed point to pixels
		return int(w >> 6)
	}

	return len(s) * (sz - 2)
}
