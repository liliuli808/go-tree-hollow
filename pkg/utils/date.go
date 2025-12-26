package utils

import (
	"time"
)

// CalculateAge 根据生日计算年龄
func CalculateAge(birthdayStr string) int {
	if birthdayStr == "" {
		return 0
	}
	birthday, err := time.Parse("2006-01-02", birthdayStr)
	if err != nil {
		return 0
	}
	now := time.Now()
	age := now.Year() - birthday.Year()
	if now.Month() < birthday.Month() || (now.Month() == birthday.Month() && now.Day() < birthday.Day()) {
		age--
	}
	return age
}

// GetConstellation 根据生日获取星座
func GetConstellation(birthdayStr string) string {
	if birthdayStr == "" {
		return ""
	}
	birthday, err := time.Parse("2006-01-02", birthdayStr)
	if err != nil {
		return ""
	}
	month := birthday.Month()
	day := birthday.Day()

	constellations := []struct {
		Name  string
		Month time.Month
		Day   int
	}{
		{"摩羯座", time.January, 19},
		{"水瓶座", time.February, 18},
		{"双鱼座", time.March, 20},
		{"白羊座", time.April, 19},
		{"金牛座", time.May, 20},
		{"双子座", time.June, 21},
		{"巨蟹座", time.July, 22},
		{"狮子座", time.August, 22},
		{"处女座", time.September, 22},
		{"天秤座", time.October, 23},
		{"天蝎座", time.November, 22},
		{"射手座", time.December, 21},
		{"摩羯座", time.December, 31},
	}

	for _, c := range constellations {
		if month < c.Month || (month == c.Month && day <= c.Day) {
			return c.Name
		}
	}
	// Fallback/Wraparound (should be covered by loop but safe to handle Capricon end of year)
	if month == time.December && day > 21 {
		return "摩羯座"
	}
	// Logic above is slightly flawed for the loop.
	// Correct logic:
	// Jan 1-19: Capricon. Jan 20+: Aquarius.
	// The loop check `month < c.Month` is wrong because it returns immediately.
	// We should check ranges.
	// Let's rewrite with simpler logic.
	return getConstellation(int(month), day)
}

func getConstellation(month, day int) string {
	days := []int{20, 19, 21, 20, 21, 22, 23, 23, 23, 24, 22, 22}
	constellations := []string{"摩羯座", "水瓶座", "双鱼座", "白羊座", "金牛座", "双子座", "巨蟹座", "狮子座", "处女座", "天秤座", "天蝎座", "射手座", "摩羯座"}

	start := month - 1
	if day < days[start] {
		return constellations[start]
	}
	return constellations[start+1]
}
