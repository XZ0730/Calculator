package test

import (
	"errors"
	"fmt"
	"log"
	"math"
	"strings"
	"testing"

	"github.com/Knetic/govaluate"
)

func TestEvaluate(t *testing.T) {

	expr, err := govaluate.NewEvaluableExpression("2!")
	if err != nil || expr == nil {
		fmt.Println("err:", err)
		t.Log("error")
		t.Fail()
	}
	result, err := expr.Evaluate(nil)
	t.Log("res:", result)
}

func TestString(t *testing.T) {
	str := "helpo"
	for _, v := range str {
		t.Log(string(v))

	}
	fmt.Println(str[2:])
}

func TestPow(t *testing.T) {
	f := math.Pow(10, 8)
	t.Log("------------<f:", f)
}

func TestTrim(t *testing.T) {
	s := strings.SplitAfter("(1+2!)!+2!", "!")
	t.Log(s)
}

func TestSin(t *testing.T) {
	exp := "1+rsin(0.5)"

	f, err := Hsin(exp)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	t.Log("res=====>", float32(f))

}

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
			return math.Asin(result * math.Pi / 180), nil
		}
		return math.Sin(result * math.Pi / 180), nil
	} else if tan {
		log.Println("tan------:", result)
		if rjudge {
			return math.Atan(result * math.Pi / 180), nil
		}
		return math.Tan(result * math.Pi / 180), nil
	}

	if rjudge {
		return math.Acos(result * math.Pi / 180), nil
	}
	log.Println("cos------:", result)
	return math.Cos(result * math.Pi / 180), nil
}

func TestTrimTec(t *testing.T) {
	exp := "1+sin(sin(10+110))+rsin"
	_, _, res, err := TrimTec(exp)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	s := strings.Replace(exp, res, "xxx", 1)
	t.Log(s)
}

func TestCalculateSC(t *testing.T) {
	// 1.判断表达式是否有三角函数
	// 2.有则放入Hsin求得第一个函数的值
	// 3.然后将值替换原本的三角函数
	// 4.回到1
	// 5.没有三角函数了在将整个表达式求解
	exp := "1+sin(sin(10+110))+rsin(0.5)"
	for {
		if strings.Contains(exp, "sin") || strings.Contains(exp, "cos") {
			res, err := Hsin(exp)
			if err != nil {
				t.Log(err)
				t.Fail()
			}
			_, _, trim_exp, err := TrimTec(exp)
			if err != nil {
				t.Log(err)
				t.Fail()
			}
			exp = strings.Replace(exp, trim_exp, fmt.Sprint(float32(res)), 1)
			t.Log("change", exp)
		} else {
			break
		}
	}

}

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
