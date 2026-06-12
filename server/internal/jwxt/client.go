package jwxt

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// Client wraps the JWXT educational system HTTP interactions.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a JWXT client with cookie jar for session persistence.
func NewClient(baseURL string) *Client {
	jar, _ := cookiejar.New(nil)
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Jar:     jar,
			Timeout: 30 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

type LoginResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ScheduleItem struct {
	Name       string `json:"name"`
	Teacher    string `json:"teacher"`
	Location   string `json:"location"`
	CourseType string `json:"course_type"`
	Credit     string `json:"credit"`
	DayOfWeek  int    `json:"day_of_week"`
	StartSec   int    `json:"start_section"`
	EndSec     int    `json:"end_section"`
	Weeks      string `json:"weeks"`
	ExamType   string `json:"exam_type"`
}

func (c *Client) Login(studentID, password string) (*LoginResult, error) {
	loginPageURL := fmt.Sprintf("%s/jwglxt/xtgl/login_slogin.html", c.BaseURL)
	req, err := http.NewRequest("GET", loginPageURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create login page request: %w", err)
	}
	c.setCommonHeaders(req)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch login page: %w", err)
	}
	resp.Body.Close()

	pubKeyURL := fmt.Sprintf("%s/jwglxt/xtgl/login_getPublicKey.html", c.BaseURL)
	req, err = http.NewRequest("GET", pubKeyURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create pubkey request: %w", err)
	}
	c.setCommonHeaders(req)

	resp, err = c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch public key: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read public key response: %w", err)
	}

	var keyResp struct {
		Modulus  string `json:"modulus"`
		Exponent string `json:"exponent"`
	}
	if err := json.Unmarshal(body, &keyResp); err != nil {
		return nil, fmt.Errorf("parse public key JSON: %w", err)
	}

	encryptedPwd, err := rsaEncrypt(password, keyResp.Modulus, keyResp.Exponent)
	if err != nil {
		return nil, fmt.Errorf("RSA encrypt password: %w", err)
	}

	loginURL := fmt.Sprintf("%s/jwglxt/xtgn/login_slogin.html", c.BaseURL)
	formData := url.Values{
		"yhm": {studentID},
		"mm":  {encryptedPwd},
	}
	req, err = http.NewRequest("POST", loginURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create login request: %w", err)
	}
	c.setCommonHeaders(req)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err = c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send login request: %w", err)
	}
	defer resp.Body.Close()

	loginBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read login response: %w", err)
	}

	bodyStr := string(loginBody)

	if resp.StatusCode == 302 || strings.Contains(bodyStr, "index_cxAreaFive") || strings.Contains(bodyStr, "xtgn/index_initMenu.html") {
		return &LoginResult{Success: true, Message: "登录成功"}, nil
	}

	if strings.Contains(bodyStr, "密码不正确") {
		return &LoginResult{Success: false, Message: "密码不正确"}, nil
	}
	if strings.Contains(bodyStr, "用户名不存在") {
		return &LoginResult{Success: false, Message: "用户名不存在"}, nil
	}
	if strings.Contains(bodyStr, "验证码") {
		return &LoginResult{Success: false, Message: "需要验证码，请稍后重试"}, nil
	}

	return &LoginResult{Success: false, Message: "登录失败，请检查学号和密码"}, nil
}

func (c *Client) GetSchedule(studentID, semester string) ([]ScheduleItem, error) {
	apiURL := fmt.Sprintf("%s/jwglxt/kbc/xsgrkbcx_cxXsgrkb.html", c.BaseURL)

	formData := url.Values{
		"xh":  {studentID},
		"xnm": {semester[:9]},
		"xqm": {semester[10:]},
	}

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create schedule request: %w", err)
	}
	c.setCommonHeaders(req)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch schedule: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read schedule response: %w", err)
	}

	return parseScheduleResponse(body)
}

func parseScheduleResponse(body []byte) ([]ScheduleItem, error) {
	var raw struct {
		KBList []map[string]interface{} `json:"kbList"`
	}

	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("parse schedule JSON: %w", err)
	}

	var items []ScheduleItem
	for _, entry := range raw.KBList {
		item := ScheduleItem{
			Name:       getString(entry, "kcmc"),
			Teacher:    getString(entry, "xm"),
			Location:   getString(entry, "cdmc"),
			CourseType: getString(entry, "kcxszc"),
			Credit:     getString(entry, "xf"),
			ExamType:   getString(entry, "khfsmc"),
		}

		item.DayOfWeek = getInt(entry, "xqj")
		item.StartSec = getInt(entry, "jcs")
		item.EndSec = getInt(entry, "jcmc")
		if item.EndSec == 0 {
			item.EndSec = item.StartSec
		}

		weekStr := getString(entry, "zcd")
		item.Weeks = parseWeeks(weekStr)

		items = append(items, item)
	}

	return items, nil
}

func parseWeeks(s string) string {
	if s == "" {
		return ""
	}

	isOdd := strings.Contains(s, "单")
	isEven := strings.Contains(s, "偶")

	s = strings.ReplaceAll(s, "周", "")
	s = strings.ReplaceAll(s, "(单)", "")
	s = strings.ReplaceAll(s, "(偶)", "")

	var weeks []string

	parts := strings.Split(s, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) == 2 {
				var start, end int
				fmt.Sscanf(rangeParts[0], "%d", &start)
				fmt.Sscanf(rangeParts[1], "%d", &end)
				for w := start; w <= end; w++ {
					if isOdd && w%2 == 0 {
						continue
					}
					if isEven && w%2 != 0 {
						continue
					}
					weeks = append(weeks, fmt.Sprintf("%d", w))
				}
			}
		} else {
			weeks = append(weeks, part)
		}
	}

	return strings.Join(weeks, ",")
}

func rsaEncrypt(password, modulusHex, exponentHex string) (string, error) {
	modulus := new(big.Int)
	exponent := new(big.Int)

	modulus.SetString(modulusHex, 16)
	exponent.SetString(exponentHex, 16)

	if modulus.Sign() == 0 || exponent.Sign() == 0 {
		return "", fmt.Errorf("invalid RSA key parameters")
	}

	pubKey := &rsa.PublicKey{
		N: modulus,
		E: int(exponent.Int64()),
	}

	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, []byte(password))
	if err != nil {
		return "", fmt.Errorf("RSA encryption failed: %w", err)
	}

	return hex.EncodeToString(encrypted), nil
}

func (c *Client) setCommonHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Referer", fmt.Sprintf("%s/jwglxt/xtgn/login_slogin.html", c.BaseURL))
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func getInt(m map[string]interface{}, key string) int {
	if v, ok := m[key]; ok {
		switch val := v.(type) {
		case float64:
			return int(val)
		case string:
			var n int
			re := regexp.MustCompile(`\d+`)
			if match := re.FindString(val); match != "" {
				fmt.Sscanf(match, "%d", &n)
			}
			return n
		}
	}
	return 0
}
