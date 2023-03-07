package utils

import (
	"errors"
	"fmt"
	"image"

	"github.com/go-vgo/robotgo"
	"gocv.io/x/gocv"
)

type ImageSize struct {
	H int
	W int
}

type Result struct {
	imgSize ImageSize
	maxVal  float32
	maxLoc  image.Point
}

func FindAllTemplates(src string) (Result, error) {
	var res Result
	var err error

	img := gocv.IMRead(src, gocv.IMReadGrayScale)
	if img.Empty() {
		panic(errors.New("failed to open image"))
	}
	defer img.Close()

	res.imgSize.H = img.Rows()
	res.imgSize.W = img.Cols()

	cap := robotgo.CaptureImg()
	capRGB, err := gocv.ImageToMatRGBA(cap)
	if err != nil {
		panic(err)
	}
	defer capRGB.Close()

	capGray := gocv.NewMat()
	defer capGray.Close()
	gocv.CvtColor(capRGB, &capGray, gocv.ColorRGBAToGray)

	result := gocv.NewMat()
	defer result.Close()

	mask := gocv.NewMat()
	defer mask.Close()

	gocv.MatchTemplate(img, capGray, &result, gocv.TmCcoeffNormed, mask)
	_, maxVal, _, maxLoc := gocv.MinMaxLoc(result)
	if maxVal < 0.95 {
		errs := fmt.Sprintf("maxVal: %f, no match result.\n", maxVal)
		return res, errors.New(errs)
	}
	res.maxVal = maxVal
	res.maxLoc = maxLoc

	return res, nil
}

func LocateCenterClick(res Result) {
	s_h, s_w := res.imgSize.H/2, res.imgSize.W/2
	x, y := res.maxLoc.X, res.maxLoc.Y

	scale := robotgo.ScaleF()
	new_x := int(float64(x) / scale)
	new_y := int(float64(y) / scale)
	new_w := int(float64(s_w) / scale)
	new_h := int(float64(s_h) / scale)

	robotgo.MouseSleep = 0
	robotgo.MoveClick(new_x+new_w, new_y+new_h)
}

func LocateRightBottomClick(res Result) {
	s_h, s_w := res.imgSize.H, res.imgSize.W
	x, y := res.maxLoc.X, res.maxLoc.Y

	scale := robotgo.ScaleF()
	new_x := int(float64(x) / scale)
	new_y := int(float64(y) / scale)
	new_w := int(float64(s_w) / scale)
	new_h := int(float64(s_h) / scale)

	robotgo.MouseSleep = 0
	robotgo.MoveClick(new_x+new_w, new_y+new_h)
}

func FindClick(img string, click int) {
	var res Result
	var err error

	res, err = FindAllTemplates(img)
	if err != nil {
		// fmt.Println(err.Error())
		return
	}

	if click == 1 {
		// 居中
		LocateCenterClick(res)
	} else if click == 2 {
		// 右下角
		LocateRightBottomClick(res)
	} else {
		panic("Invalid click type!")
	}
}
