package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/mojocn/base64Captcha"
)

// 滑块验证码尺寸常量
const (
	SliderBgWidth   = 320 // 背景图宽度
	SliderBgHeight  = 160 // 背景图高度
	SliderPieceSize = 40  // 拼块大小
	SliderTolerance = 6   // 拖动位置容差(像素)
)

var store = base64Captcha.DefaultMemStore

// GetCaptcha 获取滑块验证码，返回 (验证码ID, 背景图base64, error)
// 缺口位置随机生成并存入内存存储，前端渲染背景图后通过拖动拼块对齐缺口，
// 提交时将拼块的最终X位置传给后端，由后端校验是否命中缺口。
func GetCaptcha() (string, string, error) {
	// 生成随机缺口位置，保证拼块完整落在背景内
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	gapX := SliderPieceSize + 20 + rnd.Intn(SliderBgWidth-2*SliderPieceSize-40)
	if gapX < SliderPieceSize {
		gapX = SliderPieceSize
	}
	if gapX > SliderBgWidth-SliderPieceSize {
		gapX = SliderBgWidth - SliderPieceSize
	}

	bgBase64, err := drawSliderImage(gapX)
	if err != nil {
		return "", "", err
	}

	// 生成验证码ID并保存缺口位置答案(一次性)
	captchaID := RandString(24)
	store.Set(captchaID, strconv.Itoa(gapX))
	return captchaID, bgBase64, nil
}

// VerifyCaptcha 验证滑块验证码
// answer 为前端提交的拼块最终X位置(字符串形式的像素值)
func VerifyCaptcha(id string, answer string) bool {
	if id == "" || answer == "" {
		return false
	}
	stored := store.Get(id, true) // 取出并清除，验证码一次性使用
	if stored == "" {
		return false
	}
	gapX, err1 := strconv.Atoi(stored)
	submitX, err2 := strconv.Atoi(answer)
	if err1 != nil || err2 != nil {
		return false
	}
	diff := gapX - submitX
	if diff < 0 {
		diff = -diff
	}
	return diff <= SliderTolerance
}

// TrajectoryPoint 滑块拖动轨迹点
type TrajectoryPoint struct {
	X int   `json:"x"`
	Y int   `json:"y"`
	T int64 `json:"t"` // 毫秒时间戳
}

// VerifyTrajectory 滑块拖动轨迹检测
// 采用宽松的启发式规则拦截明显非人类的拖动轨迹(直接计算位置提交的脚本)：
//  1. 轨迹点数量足够多(真实拖动会产生大量采样点)
//  2. 拖动总时长在合理范围
//  3. 时间戳单调递增
//  4. 无单步瞬移(单步位移过大且耗时过短)
//  5. 平均速度不超过人类拖动上限
func VerifyTrajectory(data string) bool {
	if data == "" {
		return false
	}
	var points []TrajectoryPoint
	if err := json.Unmarshal([]byte(data), &points); err != nil {
		return false
	}
	if len(points) < 6 {
		return false
	}
	first, last := points[0], points[len(points)-1]
	duration := last.T - first.T
	if duration < 150 || duration > 15000 {
		return false
	}

	prevX, prevY, prevT := first.X, first.Y, first.T
	totalDist := 0.0
	for _, p := range points[1:] {
		if p.T < prevT {
			return false // 时间倒流，异常
		}
		dx := p.X - prevX
		dy := p.Y - prevY
		dist := math.Sqrt(float64(dx*dx + dy*dy))
		dt := p.T - prevT
		// 单步瞬移：极短时间间隔内位移过大(模拟器直接跳到位)
		if dt < 50 && dist > 60 {
			return false
		}
		totalDist += dist
		prevX, prevY, prevT = p.X, p.Y, p.T
	}
	// 平均速度上限：3000px/s(正常人类拖动远低于此)
	if duration > 0 && totalDist/float64(duration)*1000 > 3000 {
		return false
	}
	return true
}

// drawSliderImage 绘制滑块验证码背景图：浅色底 + 随机噪点 + 深色缺口槽位
func drawSliderImage(gapX int) (string, error) {
	img := image.NewRGBA(image.Rect(0, 0, SliderBgWidth, SliderBgHeight))

	// 浅色背景
	draw.Draw(img, img.Bounds(), &image.Uniform{C: color.RGBA{238, 241, 246, 255}}, image.Point{}, draw.Src)

	// 随机噪点，增加视觉辨识成本
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 500; i++ {
		x, y := rnd.Intn(SliderBgWidth), rnd.Intn(SliderBgHeight)
		g := uint8(rnd.Intn(80))
		img.Set(x, y, color.RGBA{g, g, g + 10, 255})
	}

	// 缺口目标槽位：深色填充块 + 深色边框（用户可见的对齐目标）
	top := (SliderBgHeight - SliderPieceSize) / 2
	slot := image.Rect(gapX, top, gapX+SliderPieceSize, top+SliderPieceSize)
	draw.Draw(img, slot, &image.Uniform{C: color.RGBA{110, 122, 148, 255}}, image.Point{}, draw.Src)
	borderColor := color.RGBA{60, 70, 90, 255}
	for x := slot.Min.X; x < slot.Max.X; x++ {
		img.Set(x, slot.Min.Y, borderColor)
		img.Set(x, slot.Max.Y-1, borderColor)
	}
	for y := slot.Min.Y; y < slot.Max.Y; y++ {
		img.Set(slot.Min.X, y, borderColor)
		img.Set(slot.Max.X-1, y, borderColor)
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", err
	}
	// 返回完整的 Data URI（与 base64Captcha 旧行为一致，可直接用于 <img src>）
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
