package encoding

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math"
	"math/rand"
	"obfuscator/config"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	dispersionPercent = config.GetConfig().Obfuscator.DispersionPercent
)

func ObfuscateValue(rawValue *interface{}, dbType string) (interface{}, error) {
	if rawValue == nil || *rawValue == nil {
		return nil, nil
	}
	if strings.HasPrefix(dbType, CharType) || strings.HasPrefix(dbType, VarcharType) {
		value, err := obfuscateString(*rawValue, dbType)
		if err != nil {
			return nil, err
		}
		return value, nil
	}
	if strings.HasPrefix(dbType, DecimalType) {
		value, err := obfuscateFloat(*rawValue, dbType)
		if err != nil {
			return nil, err
		}
		return value, nil
	}

	switch dbType {
	case TinyintType, SmallintType, MediumintType, IntType, BigintType:
		value, err := obfuscateInt(*rawValue, dbType)
		if err != nil {
			return nil, err
		}
		return value, nil
	case UTinyintType, USmallintType, UMediumintType, UIntType, UBigintType:
		value, err := obfuscateUint(*rawValue, dbType)
		if err != nil {
			return nil, err
		}
		return value, nil
	case FloatType, DoubleType:
		value, err := obfuscateFloat(*rawValue, dbType)
		if err != nil {
			return nil, err
		}
		return value, nil
	case TinytextType, TextType, MediumtextType, LongtextType:
		value, err := obfuscateString(*rawValue, dbType)
		if err != nil {
			return nil, err
		}
		return value, nil
	default:
		return *rawValue, nil
	}
}

func obfuscateInt(rawValue interface{}, dbType string) (int64, error) {
	value, err := getInt(rawValue, dbType)
	if err != nil {
		return 0, err
	}
	var lowerBound, upperBound int64
	switch dbType {
	case TinyintType:
		lowerBound = LowerBoundTinyint
		upperBound = UpperBoundTinyint
	case SmallintType:
		lowerBound = LowerBoundSmallint
		upperBound = UpperBoundSmallint
	case MediumintType:
		lowerBound = LowerBoundMediumint
		upperBound = UpperBoundMediumint
	case IntType:
		lowerBound = LowerBoundInt
		upperBound = UpperBoundInt
	case BigintType:
		lowerBound = LowerBoundBigint
		upperBound = UpperBoundBigint
	}

	operation := getRandBool()
	maxDispersionValue := absInt(value / 100 * dispersionPercent)
	dispersion := getIntDispersion(maxDispersionValue)
	if operation {
		if upperBound >= value+dispersion {
			value += dispersion
		} else {
			value -= dispersion
		}
	} else {
		if lowerBound <= value-dispersion {
			value -= dispersion
		} else {
			value += dispersion
		}
	}
	return value, nil
}

func obfuscateUint(rawValue interface{}, dbType string) (uint64, error) {
	value, err := getUint(rawValue, dbType)
	if err != nil {
		return 0, err
	}
	var upperBound uint64
	switch dbType {
	case UTinyintType:
		upperBound = UpperBoundUTinyint
	case USmallintType:
		upperBound = UpperBoundUSmallint
	case UMediumintType:
		upperBound = UpperBoundUMediumint
	case UIntType:
		upperBound = UpperBoundUInt
	case UBigintType:
		upperBound = UpperBoundUBigint
	}

	operation := getRandBool()
	maxDispersionValue := int64(value / 100 * uint64(dispersionPercent))
	dispersion := uint64(getIntDispersion(maxDispersionValue))
	if operation {
		if upperBound >= value+dispersion {
			value += dispersion
		} else {
			value -= dispersion
		}
	} else {
		value -= dispersion
	}
	return value, nil
}

func obfuscateFloat(rawValue interface{}, dbType string) (float64, error) {
	value, err := getFloat(rawValue, dbType)
	if err != nil {
		return 0, err
	}
	operation := getRandBool()
	dispersion := getFloatDispersion(value, int32(dispersionPercent))

	var upperBound, lowerBound float64
	if strings.HasPrefix(dbType, DecimalType) {
		upperBound, err = getDecimalAbsBound(dbType)
		if err != nil {
			return 0, err
		}
		lowerBound = -1 * upperBound
	}
	//probability of going out of bounds is negligible if it's float or double
	if operation {
		if strings.HasPrefix(dbType, DecimalType) && upperBound <= value+dispersion {
			value -= dispersion
		} else {
			value += dispersion
		}
	} else {
		if strings.HasPrefix(dbType, DecimalType) && lowerBound >= value-dispersion {
			value += dispersion
		} else {
			value -= dispersion
		}
	}
	//data will be truncated automatically by mysql if need
	return value, nil
}

func obfuscateString(rawValue interface{}, dbType string) (string, error) {
	value := asString(rawValue)
	value = getMD5Hash(value)
	if strings.HasPrefix(dbType, CharType) || strings.HasPrefix(dbType, VarcharType) {
		sizeStr := getSubstringInSingleLastBrackets(dbType)
		size, err := strconv.Atoi(sizeStr)
		if err != nil {
			return "", err
		}
		if utf8.RuneCountInString(value) > size {
			value = trimStr(value, size)
		}
	}
	return value, nil
}

func getMD5Hash(value string) string {
	hash := md5.Sum([]byte(value))
	return hex.EncodeToString(hash[:])
}

func getIntDispersion(maxValue int64) int64 {
	return rand.Int63n(maxValue + 1)
}

func getFloatDispersion(number float64, dispersionPercent int32) float64 {
	//step by 0.1%
	percent := rand.Int31n(dispersionPercent*100 + 1)
	return number * float64(percent) / (100 * 100)
}

func getRandBool() bool {
	return rand.Intn(2) == 1
}

func getInt(rawValue interface{}, dbType string) (int64, error) {
	s := asString(rawValue)
	switch dbType {
	case TinyintType:
		value, err := strconv.ParseInt(s, 10, 1*8)
		return value, err
	case SmallintType:
		value, err := strconv.ParseInt(s, 10, 2*8)
		return value, err
	case MediumintType:
		value, err := strconv.ParseInt(s, 10, 3*8)
		return value, err
	case IntType:
		value, err := strconv.ParseInt(s, 10, 4*8)
		return value, err
	case BigintType:
		value, err := strconv.ParseInt(s, 10, 8*8)
		return value, err
	default:
		return 0, fmt.Errorf("unknown int type: " + dbType)
	}
}

func getUint(rawValue interface{}, dbType string) (uint64, error) {
	s := asString(rawValue)
	switch dbType {
	case UTinyintType:
		value, err := strconv.ParseUint(s, 10, 1*8)
		return value, err
	case USmallintType:
		value, err := strconv.ParseUint(s, 10, 2*8)
		return value, err
	case UMediumintType:
		value, err := strconv.ParseUint(s, 10, 3*8)
		return value, err
	case UIntType:
		value, err := strconv.ParseUint(s, 10, 4*8)
		return value, err
	case UBigintType:
		value, err := strconv.ParseUint(s, 10, 8*8)
		return value, err
	default:
		return 0, fmt.Errorf("unknown uint type: " + dbType)
	}
}

func getFloat(rawValue interface{}, dbType string) (float64, error) {
	s := asString(rawValue)
	if strings.HasPrefix(dbType, DecimalType) {
		value, err := strconv.ParseFloat(s, 8*8)
		return value, err
	}

	switch dbType {
	case FloatType:
		value, err := strconv.ParseFloat(s, 4*8)
		return value, err
	case DoubleType:
		value, err := strconv.ParseFloat(s, 8*8)
		return value, err
	default:
		return 0, fmt.Errorf("unknown float type: " + dbType)
	}
}

func asString(rawValue interface{}) string {
	switch t := rawValue.(type) {
	case string:
		return t
	case []byte:
		return string(t)
	}
	rv := reflect.ValueOf(rawValue)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(rv.Uint(), 10)
	case reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 64)
	case reflect.Float32:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 32)
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool())
	}
	return fmt.Sprintf("%v", rawValue)
}

func absInt(value int64) int64 {
	if value < 0 {
		return value * -1
	}
	return value
}

func trimStr(value string, size int) string {
	runes := []rune(value)
	var sb strings.Builder
	for i, r := range runes {
		if i >= size {
			break
		}
		sb.WriteRune(r)
	}
	return sb.String()
}

func getSubstringInSingleLastBrackets(value string) string {
	return strings.Split(strings.Split(value, "(")[1], ")")[0]
}

//not included
func getDecimalAbsBound(dbType string) (float64, error) {
	sizeStr := getSubstringInSingleLastBrackets(dbType)
	leftPartStr := strings.Split(sizeStr, ",")[0]
	rightPartStr := strings.Split(sizeStr, ",")[1]
	leftPartSize, err := strconv.Atoi(leftPartStr)
	if err != nil {
		return 0, err
	}
	rightPartSize, err := strconv.Atoi(rightPartStr)
	if err != nil {
		return 0, err
	}
	intSize := leftPartSize - rightPartSize
	return math.Pow10(intSize), nil
}
