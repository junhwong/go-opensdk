package wechat

import (
	"fmt"
	"testing"
	"time"
)

func Test(t *testing.T) {
	backoff := &Backoff{Max: time.Second * 5}
	for backoff.Wait() {
		fmt.Println(backoff)
	}
}

func TestParsePrivateKey(t *testing.T) {
	var pemData = []byte(`
	-----BEGIN PRIVATE KEY-----
	MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC1yhh6LNB8nXmO
	SxdGKWmDh0OxAM/wnGyHKSD9tcEhMQTe+wabce0POXzejCmwFBzZa7ZmxH5LoAey
	T7Fpwb7pptbbDx58CxCYhNEdQ2XrFILUCq3daMj++KQlyDp8U0NspFKsO57gSlih
	AJ49DzcXQb7Vs5daIvtLapIouPyixAE5uDL+afmJ+bXC11xP5sPWw1RfXynW3vbE
	yfRol9hQyQWfmO15GSZi6TTAhTKaW31yKaQNChy06K+LsE9JAU+ESxihthtGiMbY
	3fFRyhF9Ka2e0wIOz6UdcfwMjxXWRV4OLD1uFG9IYbUiugmYtDyIYZaFDPYdi/+R
	jm10Ps5lAgMBAAECggEAb19kRZ2lEWOM8D9S//opGZrKPuvneVrsJpZtDuLGcqZM
	fKvALYXLnZMzzEiE1cpMrmuOMUHaukxNytGGOOupIg7D/SszGv3QahCc6Ne83hwP
	1wa/5DDpS0RblIYqRrbgTPQTbk+Mk48Y43K0f2YN82KlHtnLNT7PRDIDX42Nwc1X
	8f4JcfyKUE/pOSn+YUlu5Edu6QYbWJWS7mlojEZ/wuWbSymbs6mVVkKeSWGTIh1v
	4n2F3Gj6ckUDlt4aZWTVcBa2+ZvSE2h5frSH0snpdGV1bW44IqE3NkwfTQ7JI34C
	VJdhb3goIyoTmiz6NGEZuiyr8gP9IOjqPfeP7GO5YQKBgQDuB1CT8ksO4SqR3skR
	kdCQW7kOogZgDThei+3HUMOsHr8L42oYkJDmk2res1ow/mz6SoIV4w6mvvUSnACx
	dtYA1AzUEs3jvltv8cQ1HAuDhLRslWrhSoxrQQh20yrVxxGN0J4DdCAGURSUwypz
	UHR+mlfcjacPyxKUsT41+8zG+QKBgQDDg8ZGivuV794RuA3cfpitUFG+0nA0ZS3q
	AZqlA3ufnCudHQixFIsf83Q7sX7pBob5PNONqsbv0OKpC3/xJRSPIwjWTBUPlDLX
	rsGajKMhUPtkWo4zkfrSa8XaUpUVDU0qTzS71f9Aab3SkPH1d1o4cQxO08axGLbm
	TV/46QCBzQKBgDd7ZQDXPT+epHmT4HJD9sVvW9dZVPsWmckP/MC0xqdcE1QGEjjf
	mablPcfjLma1J1m//Ep1vniHkkBgNJkpBgDzbHoSWAN5335ccEug2d4yFIwq19rj
	sY9efUaVOirSV/kiY3KSotRWGeIDC+YNHtpTx58VNZes0gvutH2Iz9ahAoGAUcoW
	b/xEMv0dURxF8C+lfxtSlxlBhymsg3AYWV+Tn7mdJSS4Nhv592vI/A/Mn37zh+BC
	P8lpX3lq2HzPEPoKF7b4Q22ggdvlSQT6SMT8mTtfbyPSyRAQdWZQZnyVkTD3TvPD
	g7CKD1As8KFiFuXPAD2KgI9nVz6XhNBpjZ8rbyECgYEAsOrm1hbNZbvlNhnuUjw5
	DTgTuJ3B0j1aK/7C2EQWR+mIG2q5TKDC6xNdszV0gK1/TbJk4RNgQo0JLkuZ2Xk2
	Q8KhaNe+X8SYP9CFKIsXuhGrYI5ICjipov5oJqjESV4wle575eWwdPgF1ICabpIq
	dnX2MxS9tkk830uXxPrXpRA=
	-----END PRIVATE KEY-----`)

	key, err := ParsePrivateKey(pemData)
	t.Log(err)
	t.Log(key)
}
