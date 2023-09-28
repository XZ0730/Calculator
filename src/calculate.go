package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Knetic/govaluate"
)

var (
	entry           *widget.Entry
	filter          = []string{"+", "-", "*", "/", "%", ".", "^"}
	lastInputIsFlag bool
)

type RecType int

const (
	Sin RecType = iota
	Cos
)

func main() {
	a := app.New()
	w := a.NewWindow("Calculator")
	w.CenterOnScreen()
	w.Resize(fyne.NewSize(450, 300))
	w.SetIcon(theme.FyneLogo())
	entry = widget.NewEntry()
	entry.MultiLine = true
	entry.Resize(fyne.NewSize(150, 150))

	// 布局设置
	digits := []string{
		"^", "%",
		"7", "8", "9", "*",
		"4", "5", "6", "-",
		"1", "2", "3", "+",
	}
	var digitBtns []fyne.CanvasObject
	digitBtns = append(digitBtns, widget.NewButton("c", func() {
		entry.SetText("")
		entry.Refresh()
	}))

	// 正负切换
	digitBtns = append(digitBtns, widget.NewButton("+/-", sign()))
	// 括号
	digitBtns = append(digitBtns, container.New(layout.NewGridLayout(2),
		widget.NewButton("(", input("(")), widget.NewButton(")", input(")"))))
	// 除法
	digitBtns = append(digitBtns, widget.NewButton("/", input("/")))

	digitBtns = append(digitBtns, container.New(layout.NewGridLayout(2),
		widget.NewButton("tan", input("tan(")), widget.NewButton("rtan", input("rtan("))))
	digitBtns = append(digitBtns, widget.NewButton("rsin", input("rsin(")))
	digitBtns = append(digitBtns, widget.NewButton("rcos", input("rcos(")))
	// back
	digitBtns = append(digitBtns, widget.NewButton("back", back()))

	digitBtns = append(digitBtns, widget.NewButton("sin", input("sin(")))
	digitBtns = append(digitBtns, widget.NewButton("cos", input("cos(")))
	// 其余数字和运算符
	for _, v := range digits {
		val := v
		digitBtns = append(digitBtns, widget.NewButton(val, input(val)))
	}

	buts := container.New(layout.NewGridLayout(4), digitBtns...)

	equal := widget.NewButton("=", equals())

	lastLine := container.New(
		layout.NewGridLayout(2),
		widget.NewButton("0", input("0")),
		container.New(layout.NewGridLayout(1), widget.NewButton(".", input(".")), equal))

	w.SetContent(container.New(layout.NewVBoxLayout(), entry, buts, lastLine))
	w.ShowAndRun()
}

func percent() func() {
	return func() {

	}
}

func sign() func() {
	return func() {
		defer entry.Refresh()
		lines := strings.Split(entry.Text, "\n")
		text := lines[len(lines)-1]
		if strings.Contains(text, ".") {
			value, err := strconv.ParseFloat(text, 64)
			if err != nil {
				fmt.Println("err parse :", err)
				entry.Text = fmt.Sprint("parse error please clear\n")
				return
			}
			value = -value
			entry.Text = fmt.Sprint(value)
		} else {
			value, err := strconv.ParseInt(text, 10, 64)
			if err != nil {
				fmt.Println("err parse :", err)
				entry.Text = fmt.Sprint("parse error please clear\n")
				return
			}
			value = -value
			entry.Text = fmt.Sprint(value)
		}
	}
}

func equals() func() {
	return func() {
		// 切割换行
		lines := strings.Split(entry.Text, "\n")
		fmt.Println("lines:", len(lines))
		// 空表达式不变
		if len(lines) == 0 || (lines[0] == "" && len(lines) == 1) {
			entry.Text = ""
			entry.Refresh()
			log.Println("empty expression")
			return
		}

		line := lines[len(lines)-1]
		// 错误切除
		if len(lines) >= 3 || strings.Contains(entry.Text, "error") {
			entry.Text = line
			if strings.Contains(lines[0], "error") {
				entry.SetText("")
				log.Println("continue error calculate")
				return
			}
			entry.Refresh()
		}

		// 溢出切除
		if strings.Contains(line, "Inf") {
			entry.SetText("error:inf calculate\n")
			entry.Refresh()
			log.Println("error:inf calculate")
			return
		}

		for {
			if strings.Contains(line, "sin") || strings.Contains(line, "cos") || strings.Contains(line, "tan") {
				res, err := Hsin(line)
				if err != nil {
					entry.SetText("error:sin or cos calculate expression\n")
					entry.Refresh()
					log.Println(err)
					return
				}
				log.Println("--------------------")
				_, _, trim_exp, err := TrimTec(line)
				if err != nil {
					entry.SetText("error:sin or cos trim expression\n")
					entry.Refresh()
					log.Println(err)
					return
				}
				log.Println("--------------------")
				line = strings.Replace(line, trim_exp, fmt.Sprint(float32(res)), 1)
				log.Println("line 三角函数切割 :", line)
			} else {
				break
			}
		}
		line = strings.Trim(line, "+/x")
		// 次方替换
		if strings.Contains(line, "^") {
			line = strings.ReplaceAll(line, "^", "**")
		}

		expr, err := govaluate.NewEvaluableExpression(line)
		if err != nil || expr == nil {
			if strings.Contains(err.Error(), "transition") {
				entry.Text = fmt.Sprint("transition error")
			}
			entry.Text = fmt.Sprint("error:wrong expression\n")
			entry.Refresh()
			return
		}
		result, err := expr.Evaluate(nil)
		if err != nil || result == nil {
			if strings.Contains(err.Error(), "transition") {
				entry.Text = fmt.Sprint("transition error")
			}
			entry.Text = fmt.Sprint("error:unespected error please clear\n")
		}

		entry.Text += "=\n"
		entry.Text += fmt.Sprint(result)
		entry.Refresh()
	}
}

func input(val string) func() {
	return func() {
		errorRefresh()
		var thisFlag bool
		for _, v := range filter {
			if v == val {
				thisFlag = true
			}
		}

		if thisFlag && lastInputIsFlag {
			return
		}

		lastInputIsFlag = thisFlag
		entry.SetText(entry.Text + val)
		entry.Refresh()
	}
}

func back() func() {
	return func() {
		lines := strings.Split(entry.Text, "\n")
		line := lines[len(lines)-1]
		if strings.Contains(line, "error") {
			entry.SetText("")
			entry.Refresh()
			return
		}
		entry.SetText("")

		for i := 0; i < len(lines); i++ {
			if i != len(lines)-1 {
				entry.Text += lines[i]
				entry.Text += "\n"
			}
		}
		if len(line) == 0 {
			entry.SetText("")
			entry.Refresh()
			return
		}
		line = line[:len(line)-1]
		log.Println(line)
		entry.Text += line
		log.Println(entry.Text)
		entry.Refresh()
	}
}

func errorRefresh() {
	if strings.Contains(entry.Text, "error") || strings.Contains(entry.Text, "Inf") {
		entry.SetText("")
		entry.Refresh()
		return
	}
}

// 求解表达式的第一个三角函数值
func Hsin(exp string) (float64, error) {
	restr := ""
	preindex := 0
	lastindex := 0
	pkcnt := 0
	lkcnt := 0
	begin := false
	sin := false
	tan := false
	rjudge := false
	if strings.Contains(exp, "sin") || strings.Contains(exp, "cos") || strings.Contains(exp, "tan") {
		for index, v := range exp {
			ch := string(v)
			if begin && ch == "(" {
				pkcnt++
			} else if begin && ch == ")" {
				lkcnt++
				if lkcnt == pkcnt { //括号匹配
					lastindex = index
					restr = exp[preindex+3 : lastindex]
					break
				}
			}
			if ch == "r" {
				rjudge = true
			}
			if (ch == "i" || ch == "o" || ch == "a") && !begin {
				preindex = index
				log.Println("--------------")
				begin = true
			}
			if ch == "s" && !begin {
				sin = true
			} else if ch == "t" && !begin {
				tan = true
			}
		}
		if pkcnt != lkcnt {
			log.Println("error expresiion")
			return 0, errors.New("error expresiion")
		}
		exp = restr
	} else {
		expr, err := govaluate.NewEvaluableExpression(exp)
		if err != nil || expr == nil {
			log.Println("expression:", exp)
			log.Println("error expresiion1")
			return 0, err
		}
		result, err := expr.Evaluate(nil)
		if err != nil {
			log.Println("error expresiion2")
			return 0, err
		}
		log.Println("output:", result.(float64))
		return result.(float64), nil
	}
	log.Println("exp:", exp)
	result, err := Hsin(exp)
	if err != nil {
		return 0, err
	}
	if sin {
		log.Println("sin------:", result)
		if rjudge {
			return math.Asin(result) * 180 / math.Pi, nil
		}
		return math.Sin(result * math.Pi / 180), nil
	} else if tan {
		log.Println("tan------:", result)
		if rjudge {
			return math.Atan(result) * 180 / math.Pi, nil
		}
		return math.Tan(result * math.Pi / 180), nil
	}

	if rjudge {
		return math.Acos(result) * 180 / math.Pi, nil
	}
	log.Println("cos------:", result)
	return math.Cos(result * math.Pi / 180), nil
}

// 获得第一个cos或者sin表达式
func TrimTec(exp string) (preindex, lastindex int, res string, err error) {

	begin := false
	pkcnt := 0
	lkcnt := 0
	rjudge := false
	for index, v := range exp {
		if string(v) == "r" && !begin {
			rjudge = true
		}
		if (string(v) == "s" || string(v) == "c" || string(v) == "t") && !begin {
			preindex = index
			begin = true
		}
		if begin {
			if string(v) == "(" {
				pkcnt++
			} else if string(v) == ")" {
				lkcnt++
				if lkcnt == pkcnt {
					lastindex = index
					log.Println("preindex:", preindex)
					log.Println("lastindex:", lastindex)
					if rjudge {
						res = exp[preindex-1 : lastindex+1]
						return
					}
					res = exp[preindex : lastindex+1]
					return
				}
			}
		}
	}
	err = errors.New("error:expression error")

	return
}
