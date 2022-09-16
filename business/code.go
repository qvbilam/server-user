package business

import (
	"log"
	"strconv"
	"user/cache"
	"user/utils"
)

type UserCodeBusiness struct {
}

const maxContinuity = 4
const maxRepeat = 3

func (b *UserCodeBusiness) RandomCode(IsSpecial bool) (int64, error) {
	codes, err := b.RandomCodes(1, IsSpecial)
	if err != nil {
		return 0, err
	}
	if len(codes) != 1 {
		return 0, err
	}

	return codes[0], nil
}

func (b *UserCodeBusiness) RandomCodes(count int64, IsSpecial bool) ([]int64, error) {

	r := cache.RedisServer{}
	var err error
	var userCodes []int64
	var randoms []string
	if IsSpecial {
		randoms, err = r.RandomUserSpecialCodes(count)
	} else {
		randoms, err = r.RandomUserCodes(count)
	}

	for _, code := range randoms {
		codeInt, _ := strconv.Atoi(code)
		userCodes = append(userCodes, int64(codeInt))
	}

	if err != nil {
		return nil, err
	}

	userCodesCount := int64(len(userCodes))
	if userCodesCount < count {
		if err := b.Generate(0); err != nil {
			log.Printf("自动生成用户code失败:%v", err)
			return nil, err
		}
		// 不足递归生成
		supplementCount := count - userCodesCount

		newCodes, err := b.RandomCodes(supplementCount, IsSpecial)
		if err != nil {
			return nil, err
		}
		for _, c := range newCodes {
			userCodes = append(userCodes, c)
		}
	}

	return userCodes, nil
}

func (b *UserCodeBusiness) Generate(digit int) error {
	r := cache.RedisServer{}
	if digit == 0 {
		nowDigit := r.GetUserCodeDigit()
		digit = int(nowDigit) + 1
	}

	var count int
	if digit < 1 {
		count = 1
	} else {
		count = digit - 1
	}

	min := int(utils.Pow(10, count))
	max := int(utils.Pow(10, count+1) - 1)

	var codes, specialCodes []interface{}

	for n := min; n <= max; n++ {
		if isContinuity(n, maxContinuity) || isRepeat(n, maxRepeat) {
			specialCodes = append(specialCodes, n)
			continue
		}
		codes = append(codes, n)
	}

	log.Printf("Generate3: 生产用户code: %v\n", codes)

	if _, err := r.GenerateUserCodes(int64(digit), codes); err != nil {
		return err
	}
	if _, err := r.GenerateUserSpecialCodes(int64(digit), specialCodes); err != nil {
		return err
	}

	return nil
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
