package events

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/YumeMichi/yohane/utils"
)

var (
	NormalImagePath = ImagePath + "normal/"
)

func Round(ch chan string) {
	// 全局设定
	effectTime := time.Now().Add(-time.Hour * 2) // 初始道具使用时间

	for {
		// 开始时间
		startTime := time.Now()

		// 初始设定
		isAlive := true      // 存活状态
		needCheckBox := true // 检查开箱

		// 特效道具
		if time.Since(effectTime).Hours() > 1 {
			fmt.Println("检查特效道具")
			utils.FindClick(CommonImagePath+"effects.png", 1)
			for {
				_, err := utils.FindAllTemplates(CommonImagePath + "close.png")
				if err == nil {
					break
				}
				time.Sleep(time.Second)
			}
			_, err := utils.FindAllTemplates(CommonImagePath + "exp.png")
			if err == nil {
				utils.FindClick(CommonImagePath+"exp.png", 2)
				time.Sleep(time.Second)

				utils.FindClick(CommonImagePath+"ok.png", 1)

				// 更新特效道具使用时间
				fmt.Println("更新特效道具使用时间")
				effectTime = time.Now()

				time.Sleep(time.Second * 5)
			}
			utils.FindClick(CommonImagePath+"close.png", 1)
			time.Sleep(time.Second)
		}

		// 歌曲入口
		utils.FindClick(NormalImagePath+"song.png", 1)
		fmt.Println("歌曲入口")
		time.Sleep(time.Second * 3)

		// 检查体力
		fmt.Println("检查体力")
		_, err := utils.FindAllTemplates(CommonImagePath + "stamina.png")
		if err == nil {
			// 体力不足
			fmt.Println("体力不足")
			utils.FindClick(CommonImagePath+"ok.png", 1)
			time.Sleep(time.Second)

			// 使用 100% 糖罐
			utils.FindClick(CommonImagePath+"candy.png", 2)
			fmt.Println("使用 100% 糖罐")
			time.Sleep(time.Second)

			// 使用确认
			utils.FindClick(CommonImagePath+"ok.png", 1)
			time.Sleep(time.Second)
			utils.FindClick(CommonImagePath+"ok.png", 1)
			time.Sleep(time.Second * 5)
		}
		fmt.Println("检查体力结束")

		// 嘉宾选择
		fmt.Println("嘉宾选择")
		_, err = utils.FindAllTemplates(NormalImagePath + "friend_10.png")
		if err == nil {
			utils.FindClick(NormalImagePath+"friend_10.png", 2)
			time.Sleep(time.Second)
		} else {
			utils.FindClick(NormalImagePath+"friend_0.png", 2)
			time.Sleep(time.Second)
		}
		fmt.Println("嘉宾选择结束")
		utils.FindClick(CommonImagePath+"ok.png", 1)

		// 开始演唱会
		fmt.Println("开始演唱会")
		utils.FindClick(NormalImagePath+"ok.png", 1)

		// 检查演唱会状态 完成状态/存活状态
		for {
			fmt.Println("演唱会中，3秒后重试")
			time.Sleep(time.Second * 3)

			_, err = utils.FindAllTemplates(CommonImagePath + "failed.png")
			if err == nil {
				// 演唱会失败
				utils.FindClick(CommonImagePath+"cancel.png", 1)
				time.Sleep(time.Second)
				utils.FindClick(CommonImagePath+"ok.png", 1)
				fmt.Println("演唱会失败")
				time.Sleep(time.Second)
				isAlive = false
				break
			}

			_, err = utils.FindAllTemplates(CommonImagePath + "error_1.png")
			if err == nil {
				// 网络连接失败/重试
				utils.FindClick(CommonImagePath+"ok.png", 1)
				fmt.Println("网络连接失败/重试")
				continue
			}

			_, err = utils.FindAllTemplates(CommonImagePath + "ok.png")
			if err == nil {
				// 演唱会结束/开箱
				utils.FindClick(CommonImagePath+"ok.png", 1)
				fmt.Println("演唱会结束")
				break
			}

			_, err = utils.FindAllTemplates(CommonImagePath + "click.png")
			if err == nil {
				// 演唱会结束/fallback
				utils.FindClick(CommonImagePath+"click.png", 1)
				fmt.Println("演唱会结束/fallback")
				needCheckBox = false
				break
			}
		}

		if isAlive {
			// 检查开箱子
			if needCheckBox {
				fmt.Println("检查开箱子")
				for {
					_, err = utils.FindAllTemplates(CommonImagePath + "click.png")
					if err == nil {
						utils.FindClick(CommonImagePath+"click.png", 1)
						fmt.Println("检查开箱子结束")
						utils.FindClick(NormalImagePath+"complete_1.png", 1)
						// time.Sleep(time.Second * 3)
						break
					}
					time.Sleep(time.Second)
					utils.FindClick(CommonImagePath+"ok.png", 1)
				}
			}

			// 等待结算/歌曲数据
			fmt.Println("结算/歌曲数据")
			for {
				_, err = utils.FindAllTemplates(NormalImagePath + "complete_2.png")
				if err == nil {
					fmt.Println("结算/歌曲数据结束")
					utils.FindClick(NormalImagePath+"complete_2.png", 1)
					time.Sleep(time.Second)
					break
				}
				time.Sleep(time.Second)
			}

			// 检查社员上限
			fmt.Println("检查社员上限")
			for {
				_, err = utils.FindAllTemplates(CommonImagePath + "click.png")
				if err == nil {
					fmt.Println("检查社员上限结束")
					utils.FindClick(NormalImagePath+"complete_3.png", 1)
					time.Sleep(time.Second)
					break
				}
				utils.FindClick(CommonImagePath+"ok.png", 1)
				time.Sleep(time.Second)
			}

			// 等待结算/绊数据
			fmt.Println("结算/绊数据")
			for {
				_, err = utils.FindAllTemplates(CommonImagePath + "kitsuna.png")
				if err == nil {
					utils.FindClick(CommonImagePath+"kitsuna.png", 1)
					time.Sleep(time.Second)
					break
				}
				time.Sleep(time.Second)
			}

			// 歌曲目标（如果存在的话）
			_, err = utils.FindAllTemplates(CommonImagePath + "ok.png")
			if err == nil {
				utils.FindClick(CommonImagePath+"ok.png", 1)
				time.Sleep(time.Second)
				utils.FindClick(CommonImagePath+"close.png", 1)
				time.Sleep(time.Second)
				fmt.Println("结算/歌曲目标结束")
				time.Sleep(time.Second * 2)
				utils.FindClick(CommonImagePath+"kitsuna.png", 1)
			}

			// 等待结算/绊数据
			fmt.Println("结算/绊数据结束")
			time.Sleep(time.Second * 2)

			// 爱心课题（如果存在的话）
			_, err = utils.FindAllTemplates(CommonImagePath + "keti_1.png")
			if err == nil {
				utils.FindClick(CommonImagePath+"keti_1.png", 1)
				time.Sleep(time.Second)
				utils.FindClick(CommonImagePath+"keti_2.png", 1)
				time.Sleep(time.Second)
				fmt.Println("结算/爱心课题结束")
				time.Sleep(time.Second * 2)
			}

			// 每周课题（如果存在的话）
			_, err = utils.FindAllTemplates(CommonImagePath + "week.png")
			if err == nil {
				utils.FindClick(CommonImagePath+"week.png", 1)
				fmt.Println("结算/每周课题结束")
				time.Sleep(time.Second * 2)
			}
		}

		// 结束时间
		endTime := time.Now()
		timeDiff := endTime.Sub(startTime).String()
		ch <- "回合结束, 耗时: " + timeDiff

		time.Sleep(time.Duration(rand.Intn(3)+3) * time.Second)
	}
}
