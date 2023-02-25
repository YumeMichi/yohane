package main

import (
	"errors"
	"fmt"
	"image"
	"math/rand"
	"time"

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

func CheckExpInUse() bool {
	_, err := FindAllTemplates("images/coin.png")
	return err != nil
}

func FindClick(img string, click int) {
	var res Result
	var err error

	res, err = FindAllTemplates(img)
	if err != nil {
		fmt.Println(err.Error())
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

func main() {
	for {
		// 存活状态
		isAlive := true

		// 特效道具
		FindClick("images/effects.png", 1)
		time.Sleep(time.Second * 3)
		_, err := FindAllTemplates("images/exp.png")
		if err == nil {
			FindClick("images/exp.png", 2)
			time.Sleep(time.Second)
			FindClick("images/ok.png", 1)
			time.Sleep(time.Second * 5)
		}
		FindClick("images/close.png", 1)
		time.Sleep(time.Second)

		// 活动入口
		FindClick("images/event.png", 1)
		fmt.Println("活动入口")
		time.Sleep(time.Second)

		// 难度选择
		FindClick("images/master.png", 1)
		fmt.Println("难度选择")
		time.Sleep(time.Second)

		// 难度确认
		FindClick("images/ok.png", 1)
		fmt.Println("难度确认")
		time.Sleep(time.Second * 3)

		// 检查体力
		_, err = FindAllTemplates("images/st.png")
		if err == nil {
			// 体力不足
			fmt.Println("体力不足")
			FindClick("images/ok.png", 1)
			time.Sleep(time.Second)

			// 使用 100% 糖罐
			FindClick("images/candy.png", 2)
			fmt.Println("使用 100% 糖罐")
			time.Sleep(time.Second)

			// 使用确认
			FindClick("images/ok.png", 1)
			time.Sleep(time.Second)
			FindClick("images/ok.png", 1)
			time.Sleep(time.Second * 5)
		}

		// 开始演唱会
		fmt.Println("开始演唱会")
		FindClick("images/ok.png", 1)

		// 检查演唱会状态 完成状态/存活状态
		for {
			fmt.Println("演唱会中，5秒后重试")
			time.Sleep(time.Second * 5)

			_, err = FindAllTemplates("images/rank.png")
			if err == nil {
				// 演唱会失败
				FindClick("images/rank.png", 1)
				time.Sleep(time.Second)
				FindClick("images/confirm.png", 1)
				fmt.Println("演唱会失败")
				time.Sleep(time.Second)
				isAlive = false
				break
			}

			_, err = FindAllTemplates("images/ok.png")
			if err == nil {
				// 演唱会结束
				FindClick("images/ok.png", 1)
				fmt.Println("演唱会结束")
				time.Sleep(time.Second * 10)
				break
			}
		}

		if isAlive {
			// 检查开箱子
			for {
				_, err = FindAllTemplates("images/ok.png")
				if err != nil {
					fmt.Println("开箱子结束")
					time.Sleep(time.Second)
					break
				}
				FindClick("images/ok.png", 1)
				fmt.Println("检查开箱子")
				time.Sleep(time.Second * 10)
			}

			// 等待结算/歌曲数据
			FindClick("images/complete.png", 1)
			time.Sleep(time.Second * 5)

			// 检查社员上限
			for {
				_, err = FindAllTemplates("images/ok.png")
				if err != nil {
					fmt.Println("检查社员上限结束")
					time.Sleep(time.Second)
					break
				}
				FindClick("images/ok.png", 1)
				fmt.Println("检查社员上限")
				time.Sleep(time.Second * 3)
			}

			// 等待结算/社员数据
			FindClick("images/complete.png", 1)
			fmt.Println("结算/社员数据结束")
			time.Sleep(time.Second * 5)

			// 等待结算/活动数据
			FindClick("images/rank.png", 1)
			fmt.Println("结算/活动数据结束")
			time.Sleep(time.Second * 2)

			// 等待结算/绊数据
			FindClick("images/kitsuna.png", 1)
			fmt.Println("结算/绊数据结束")
			time.Sleep(time.Second * 2)

			// 每周课题
			FindClick("images/effects.png", 1)
			time.Sleep(time.Second * 2)
		}

		time.Sleep(time.Duration(rand.Intn(3)+3) * time.Second)
	}
}
