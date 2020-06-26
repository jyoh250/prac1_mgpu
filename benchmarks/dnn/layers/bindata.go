// Code generated by "esc -private -pkg=layers -o=bindata.go trans.hsaco gpu_gemm.hsaco relu.hsaco maxpooling.hsaco"; DO NOT EDIT.

package layers

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDirectory struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	once sync.Once
	data []byte
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDirectory) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	if !f.isDir {
		return nil, fmt.Errorf(" escFile.Readdir: '%s' is not directory", f.name)
	}

	fis, ok := _escDirs[f.local]
	if !ok {
		return nil, fmt.Errorf(" escFile.Readdir: '%s' is directory, but we have no info about content of this dir, local=%s", f.name, f.local)
	}
	limit := count
	if count <= 0 || limit > len(fis) {
		limit = len(fis)
	}

	if len(fis) == 0 && count > 0 {
		return nil, io.EOF
	}

	return fis[0:limit], nil
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// _escFS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func _escFS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// _escDir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func _escDir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// _escFSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func _escFSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		_ = f.Close()
		return b, err
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// _escFSMustByte is the same as _escFSByte, but panics if name is not present.
func _escFSMustByte(useLocal bool, name string) []byte {
	b, err := _escFSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// _escFSString is the string version of _escFSByte.
func _escFSString(useLocal bool, name string) (string, error) {
	b, err := _escFSByte(useLocal, name)
	return string(b), err
}

// _escFSMustString is the string version of _escFSMustByte.
func _escFSMustString(useLocal bool, name string) string {
	return string(_escFSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/gpu_gemm.hsaco": {
		name:    "gpu_gemm.hsaco",
		local:   "gpu_gemm.hsaco",
		size:    13800,
		modtime: 1591217294,
		compressed: `
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
/xcAAP//aVmkveg1AAA=
`,
	},

	"/maxpooling.hsaco": {
		name:    "maxpooling.hsaco",
		local:   "maxpooling.hsaco",
		size:    22672,
		modtime: 1593133381,
		compressed: `
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
`,
	},

	"/relu.hsaco": {
		name:    "relu.hsaco",
		local:   "relu.hsaco",
		size:    13904,
		modtime: 1591217294,
		compressed: `
H4sIAAAAAAAC/+xbzW/U1hY//shkMuG9x3ugJ8hb4AfvKYCenRlPPiZZvJKPJkEkIZDykVYI3dh3JpN4
7JHHE5KKDmkWBVWRSkvXbVddsehfkETqP1ChLll0w76q1C6Zyva94w9iAiWUUu5PGh/7fNx7ju+x5/he
+9bbU+M8x50BAgF+AM7dkfxjKpjN+/S0xytAGs7AAchACgDEkF6c7nBRmiZ8jtgl4evDUQoHA7u2kH9x
+r0QpWE711fIEn6MViFKqR3/nHY0vouPHF18Bruwfy4uPHL0FDw/RHo+eQgcD9N/RKkYskuT/oenxzx1
Ojb/8vLB54vQ3oqN8oanxyZmL/m6RwGgk/BRRS9ppowquvtbrCFZLhVXC9k8affW3wEyRFeW5cxlbNfK
ljkkUbwn5f4nZaVrmXPYNrFRG8pIkizNoAoOdCRJuoinLo1b9g1k6xn3eG6tsmAZIbXukMaZZb3b05pC
ZqmOSkFT56vYHJ2SRiPSlk+eL6p0zZMO2yXPGRe7OKRZddPJ0KN31qo4olIOCefK70dte1uiYaNcMod2
FV1GRh2fK5s6FY+seayogtsxVTibV4OGNe1CHRlB02O4iOqGkxxQ2UyOprtoWMg53Z0cUiE5pEJySBOG
tYCMkXqxiO3kuMbDcem6PVdFGqbR+U28QNxW3XkjAt8nzyfLuo5Nv/PzxWINO1efkpH9vS+///lX3P+7
r6D/Gct82o2g8Ixpw7x6ZV5N1w2nPGGX9bk1Uxu2Sy/s4ail41nbqrb+tNw/VGSX5nCpgk3Hd76QJcIJ
26pXiWi8vIp1X07Fs3Z5BTk4WSHaOAmfOnoFreCibdFOJal1GczUK3MTsxdrrVOV7wsklyMSKphGq+MG
cq5Y9rLvtNem2tefXCiMIG15j0qBqrBS4Tf8ZRb9Muv6m1oysFKJlUqsVGKlEvPqz1wqFV6bUql/z1JJ
UZRM4nwSBwBdQgrgOJlPIxNUf6V8Mt/2ORfo01+Xx/Hmkg66xzfvnb59T9+eF8jcEcSmplrzX+EiDEJz
N8DAwMDAwMDAwMDAsB/gSB3Oeau7QrAQnaTP3Ye73lpv9NlhNrR/wl+hD9amRTHVbDabf8T4N4DfEQHg
W+B3/LXt1M4RAOiAj7cy8NXWLdjc5pvcR673aS79WRpABeBVgb/zwZK0vnUI7mwLkNlxz946Dw2ATx98
ATy4dhyfagCIqiCkCgCzD3kAXuDTDR5AFQW+UIXN7Z9EMbMuiv8EqD70H4g+3AZY/91++zX+VyPjn3qt
x1/yxn9z60Di+IMqgD/+J9zx5+n4iw0QUg3BzQGRb+WS20YHn2mkUym1rSPt5YIAIHwJ/HE3B9ZTGzeX
pI2tdrizzcMnDzaAh04+3ciIoiq2B7nT4B/fBiL3cwtUgeRRNH8YGBgYGBgYGBgYGBieROtJjbxn30kO
jxBKn+Q3iZw+9XUR+vPjpuXSSSKn75Ubh3fvb3J0VNIMZJakFf+1KSmXVbJKVjrZoyMH9Sxhc7ls1uQb
lr1cqyIN92hWpVp3sGxbWkW2sSHnlWwPXnWwbSKjZ1HTZMdyegxjpSJXbWsJa06P30FvdqCvP9/fO5Ab
xIU8UnVd70MLOIuKRTyQzQ6oeU3XC33qKenkAqphXbJMyXUvr2SV3GDvYF4e6MNosE+VSUvSKZgqm8vY
HpKmpsZehuOGoT+723vO67ijO/bvKL+d8O/G+H+h+iei/ENUP8Y/Svj3Y/xjhP9djP9/d8O3B99BEPw3
YZ1ZTVhnBsW0HAyKvmbW1iqglMy6sohqi0C2Lt+xQXHwquMdoUpZA0WzKhVsOqDU1ioOWgCltlhzbH/P
pzAykr2e87YqjIzkvP0c2c97297kBerrY/Mzw9NnR/dlPq49vDae9B0HRK+38Ph2hu0PR+lk6LrmQt+r
0Ov9bwDwS7NpUXt6XVP6n5g7aXgyL9pCcnofoFSK2Ysxeoy8I8DH7juUdu2a5wG6w9/6QPL3QUkNyMS2
pZbw3U5bLH7yGQ/0kyZjaQ5Vwpjmdu+e0rdi70VQ7OR8+iME9+n0LuM3EXvHguIb1afze5y/Cwn2Hfmg
/afZ/xoAAP//bciEz1A2AAA=
`,
	},

	"/trans.hsaco": {
		name:    "trans.hsaco",
		local:   "trans.hsaco",
		size:    9656,
		modtime: 1591217294,
		compressed: `
H4sIAAAAAAAC/+xaTXPbxBt/JCuO4+T/B0605YAIh7QdJCtynDi50Lw0SadJmjbQUphOZyOtbTV6G2kd
YgbSF2inM3SGl+FePgCHHjm1OfABmJ576IHO8AXgiBlJu7bkWn2BdgpUv5n4sZ6XfX7PWqvsrvbi0ZVF
nuOOAEUO7gMXfHk1umaGH0cjeTjUVaEAR2AEipAHACHm1yv3uKQsUD1H49Lw9khSMj5B3EDsulf+yiVl
PC7gCiLV90gXkpLF8U8Zx+o79YDowhPExfkFOPmA6Hl4egisP3noEo/JW0NJKcTiCjT/7OpC6M5+mzfC
+yHSCzDYqY3pZlcXltbfj3z3A8Aw1SNLr2u2hCw9+Gv4SJLqtZ2qUqbtykMAReorSVLxNPZ8w7FnRIaP
xPF3REU8VzyOPRub/kxRFCVxDVm46yOK4nsesn3X8XExuNpoWZuOGXMa69iPbOljoc8KsutNVO82c8LF
9vyKOJ+wdviEPFTxXGid9eohkQB9yByzix1eLRcn7GM100Hk8FjHY8P4JBld7ZhmTaNuz/Q1nUZmEx83
bJ2Zl0xnE5lzzVoNe0mvgALzWiyr3dZ13dtwkYZPNpE502mia9c0ZomwgGuoaZL0uk80yctZuNMkbpOc
MXTSSO8AwybptU+k1z6RXvtcK1Sll30sXvZfLGsZG/UG+RfX9YxutGVD17Ed3SsnajUfkw8eQXBy4vnn
P/uC83/4AvKvOfaj7ovqE47yjNULY7XaNImx5Bn6RsvWZr3632Y47+h43XPczr/kYKqAvPoGrlvYJhH5
qkKNS57TdKlp0djBemRn5nXP2EYEpzskG6flM6Jn0DaueQ5LKoqdYbDWtDaW1k/5na4qV7qW0wkLM6yi
nUUTkTOOtxWRDttUK5NFWZaL6fO/YM52IJd/aB7Mxf4O0OmmSK+PTg9/x8emily8wc7kCTJkyPAfBUfH
Pxeu7nI9D4E+/twP8PVQsNZLPozW42v3aIXeXZsKQr7dbrf/ifXzwO8Fa9LPIb8XPBtzueLePgjWp1/e
Hobvb1+EG3egDdcC9gUofDtSKFwWrsGlgUvcFY7P7wIv7AIs3+MB+AKACvDTzzzw4MKNOzxcvTsk8PB/
QVD5HK/m+OufwYVL10D84vYBuB7aYRAgxxX3wr4dENSbQn6Uh6/uXhF4CHIP8IVdQRBUYTBfBVgP83Bw
9S43ADDI53c5DtSbwI9CEAM85AHUXI6vBvkB3HvRw/3ynexOz5AhQ4YMGTJkyJAhQ4aXG+xd863hSFIB
+6hkK/l99D08W/Uzv9/+aDuBvE8VnXf6I/3zLc/Pi5qJ7Lq4Hb1lFccVWZEV8WBJRwSVLmB7y7B96WPH
2/JdpOGS5lhuk2DJczRL8rAplWWlhHcI9mxklhqaJhGHlExz25Jcz7mANVKKEkwoU5XJ8uTE1Pg0rpaR
qut6BW1iBdVqeEpRptSypuvVinpIPLiJfKyLji0G9MqyIo9PT0yXpakKRtMVVaItiYdgxbC3sDcjrqws
PA/ipqk/Oe3H7usEv+6nryX1g1T/S4/+9XAzZLB7ToHifyn7yiDbDsEg6y3bb1kg1+2m3EB+A+hnoCce
yATvkPAKWYYGsuZYFrYJyH7LImgTZL/hEy/6FkmYm1POj4efamwT+vzC2bXZ1WPzz3LfazC29512XqKz
hwUP9+NwLIyNHyaV2PjhYudC2Lh6BQB+b7cdFs/GD5NiD61CT/79tG2+d7yNJPPwPfyZfJOe4+B7xjeT
w33vpy7G4mdqIP0cTloDEo3NMUXK+ZiBnvpZmknapNKTxqXxq1z/9Ey+G++7GPbeiuRO7Dkn9Pn9luLc
Y6jSc1JnH9N/J1PivxlNdkda/J8BAAD//5VIJrK4JQAA
`,
	},
}

var _escDirs = map[string][]os.FileInfo{}
