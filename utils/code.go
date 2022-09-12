package utils

import (
	"user/cache"
)

type Random struct {
	Count     int64
	IsSpecial bool
}

type Generator struct {
	Digit                  int64
	MaxRepeat              int64
	MaxContinuity          int64
	IsGenerateCodes        bool
	IsGenerateSpecialCodes bool
}

type UserCoder struct {
	Random    *Random
	Generator *Generator
}

func NewDefaultCoder() *UserCoder {
	return &UserCoder{
		Random: &Random{
			Count:     1,
			IsSpecial: false,
		},
		Generator: &Generator{
			Digit:                  4,
			MaxRepeat:              3,
			MaxContinuity:          3,
			IsGenerateCodes:        true,
			IsGenerateSpecialCodes: true,
		},
	}
}

func (r *Random) RandomUserCode() ([]string, error) {
	s := cache.RedisServer{}
	var err error
	var userCodes []string
	if r.IsSpecial {
		userCodes, err = s.RandomUserCodes(r.Count)
	} else {
		userCodes, err = s.RandomUserSpecialCodes(r.Count)
	}

	if err != nil {
		return nil, err
	}

	userCodesCount := len(userCodes)
	if int64(userCodesCount) < r.Count {
		if err := GenerateUserCode(0); err != nil {
			return nil, err
		}
	}
}

func GenerateUserCode(digit int) error {
	s := cache.RedisServer{}
	if digit == 0 {
		nowDigit := s.GetUserCodeDigit()
		digit = int(nowDigit) + 1
	}

	var count int
	if digit < 1 {
		count = 1
	} else {
		count = digit - 1
	}

	min := int(pow(10, count))
	max := int(pow(10, count+1) - 1)

	maxContinuity := 4
	maxRepeat := 3

	var codes, specialCodes []interface{}

	for n := min - 1; n <= max; n++ {
		if isContinuity(n, maxContinuity) || isRepeat(n, maxRepeat) {
			specialCodes = append(specialCodes, n)
			continue
		}
		codes = append(codes, n)
	}
	if _, err := s.GenerateUserCodes(int64(digit), codes); err != nil {
		return err
	}
	if _, err := s.GenerateUserSpecialCodes(int64(digit), specialCodes); err != nil {
		return err
	}

	return nil
}

func pow(x float64, n int) float64 {
	if n == 0 {
		return 1
	}

	return x * pow(x, n-1)
}

// 是否连续
func isContinuity(n, lens int) bool {
	//统计正顺次数 12345
	z := 0
	//统计反顺次数  654321
	f := 0
	//判断3个数字是否是顺子，只需要判断2次
	lens = lens - 1
	for {
		// 个位数
		g := n % 10
		n = n / 10
		// 十位数
		s := n % 10

		if s-g == 1 {
			f = f + 1
		} else {
			f = 0
		}

		if g-s == 1 {
			z = z + 1
		} else {
			z = 0
		}

		if f == lens || z == lens {
			return true
		}

		if n == 0 {
			return false
		}
	}
}

// 是否重复数字
func isRepeat(n, lens int) bool {
	c := 0
	lens = lens - 1
	var g, s int
	for {
		g = n % 10
		n = n / 10
		s = n % 10

		if s == g {
			c = c + 1
		} else {
			c = 0
		}

		if c == lens {
			return true
		}

		if n == 0 {
			return false
		}
	}
}
