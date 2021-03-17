package DanmuCenter

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/dlclark/regexp2"
)

const replaceDictText = "ÀÁÂÃÄÅàáâãäåĀāĂăĄąȀȁȂȃȦȧɑΆΑάαАаӐӑӒӓ:a;ƀƁƂƃƄƅɃʙΒβВЬвЪъьѢѣҌҍ:b;ÇçĆćĈĉĊċČčƇƈϲϹСсҪҫ:c;ÐĎďĐđƉƊƋƌȡɖɗ:d;ÈÉÊËèéêëĒēĔĕĖėĘęĚěȄȅȆȇȨȩɐΈΕЀЁЕеѐёҼҽҾҿӖӗ:e;Ƒƒƭ:f;ĜĝĞğĠġĢģƓɠɡɢʛԌԍ:g;ĤĥĦħȞȟʜɦʰʱΉΗНнћҢңҤҺһӇӈӉӊԊԋ:h;ÌÍÎÏìíîïĨĩĪīĬĭĮįİıƗȈȉȊȋɪΊΙΪϊії:i;ĴĵʲͿϳ:j;ĶķĸƘƙΚκϏЌКкќҚқҜҝҞҟҠҡԞԟ:k;ĹĺĻļĽľĿŀŁłȴɭʟӏ:l;ɱʍΜϺϻМмӍӎ:m;ÑñŃńŅņŇňŉŊŋƝƞȵɴΝηПп:n;ÒÓÔÕÖòóôõöŌōŎŏŐőơƢȌȍȎȏȪȫȬȭȮȯȰȱΌΟοόОоӦӧ:o;ƤΡρϼРр:p;ɊɋԚԛ:q;ŔŕŖŗŘřƦȐȑȒȓɌɍʀʳг:r;ŚśŜŝŞşŠšȘșȿЅѕ:s;ŢţŤťŦŧƫƬƮȚțͲͳΤТтҬҭ:t;ÙÚÛÜùúûŨũŪūŬŭŮůŰűŲųƯưƱȔȕȖȗ:u;ƔƲʋνυϋύΰѴѵѶѷ:v;ŴŵƜɯɰʷωώϢϣШЩшщѡѿԜԝ:w;ΧχХхҲҳӼӽ:x;ÝýÿŶŷŸƳƴȲȳɎɏʏʸΎΥΫϒϓϔЎУуўҮүӮӯӰӱӲӳ:y;ŹźŻżŽžƵƶȤȥʐʑΖ:z;o:0;∃э:3;➏:6;┑┐┓┑:7;╬╪:+"

var replaceDict = getReplaceDict()
var replaceSpecial = regexp.MustCompile(`[.|/\\@\-*&^ +-]`)

func getReplaceDict() map[rune]rune {
	result := map[rune]rune{}
	for _, list := range strings.Split(replaceDictText, ";") {
		item := strings.Split(list, ":")
		replaceResult := []rune(item[1])[0]
		for _, raw := range item[0] {
			result[raw] = replaceResult
		}
	}
	return result
}

func ReplaceSimilarAndNumberRune(rawData string) string {
	result := strings.Map(func(code rune) rune {
		if result, ok := replaceDict[code]; ok {
			code = result
		}
		if code == 12288 {
			return -1
		}
		if code > 65280 && code < 65375 {
			return code - 65248
		}
		if code >= 0x0030 && code <= 0x0039 ||
			code >= 0x2460 && code <= 0x249b ||
			code >= 0x3220 && code <= 0x3229 ||
			code >= 0x3248 && code <= 0x324f ||
			code >= 0x3251 && code <= 0x325f ||
			code >= 0x3280 && code <= 0x3289 ||
			code >= 0x32b1 && code <= 0x32bf ||
			code >= 0xff10 && code <= 0xff19 {
			return '#'
		}
		return code
	}, rawData)
	result = replaceSpecial.ReplaceAllString(strings.ToLower(result), "")
	return result
}

func CompressRepeatGroup(repeatGroupMinLen int) func(string) string {
	replaceGroup := regexp2.MustCompile(`(.{`+strconv.Itoa(repeatGroupMinLen)+`,}?)\1+`, 0)
	return func(s string) string {
		if result, err := replaceGroup.Replace(s, "$1", -1, -1); err == nil {
			return result
		}
		return s
	}
}
