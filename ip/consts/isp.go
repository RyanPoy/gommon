package consts

// ISP 运营商
// English Name : Chinese Name
var ISPDict = map[string]string{
	"AIPU":     "爱普",
	"ALI":      "阿里云",
	"BGCTV":    "歌华有线",
	"BGCTVNET": "歌华有线",
	"CATV":     "视讯宽带",
	"CERNET":   "教育网",
	"CITIC":    "中信网络",
	"CMNET":    "移动",
	"CNC":      "联通",
	"CRTC":     "铁通",
	"CT":       "电信",
	"DXT":      "鹏博士",
	"EOC":      "广电",
	"FWBN":     "方正宽带",
	"HKBN":     "香港宽频",
	"OTHER":    "其他",
	"SCC":      "有线通",
	"TECH":     "科技网",
	"TWNET":    "天威视讯",
	"WASU":     "华数宽带",
	"JD":       "京东",
}

// IspName 运营商名称
func IspName(ispCode string) string {
	if name, ok := ISPDict[ispCode]; ok {
		return name
	}
	return ""
}
