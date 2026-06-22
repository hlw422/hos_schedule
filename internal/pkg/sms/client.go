package sms

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"hos_schedule/internal/config"
)

type Client struct {
	cfg *config.SMSConfig
}

func NewClient(cfg *config.SMSConfig) *Client {
	return &Client{cfg: cfg}
}

func (c *Client) SendSMS(phone, templateID string, params map[string]string) error {
	if c.cfg == nil || c.cfg.SecretID == "" {
		log.Printf("SMS not configured, skip sending to %s, template: %s, params: %v", phone, templateID, params)
		return nil
	}

	switch c.cfg.Provider {
	case "tencent":
		return c.sendTencentSMS(phone, templateID, params)
	case "aliyun":
		return c.sendAliyunSMS(phone, templateID, params)
	default:
		log.Printf("Unknown SMS provider: %s, skip sending", c.cfg.Provider)
		return nil
	}
}

func (c *Client) sendTencentSMS(phone, templateID string, params map[string]string) error {
	tplID := templateID
	if c.cfg.TemplateID != "" {
		tplID = c.cfg.TemplateID
	}

	var tplParams []string
	for _, v := range params {
		tplParams = append(tplParams, v)
	}

	reqBody := map[string]interface{}{
		"PhoneNumberSet": []string{phone},
		"SmsSdkAppId":    c.cfg.AppID,
		"SignName":       c.cfg.SignName,
		"TemplateId":     tplID,
		"TemplateParamSet": tplParams,
	}

	body, _ := json.Marshal(reqBody)

	host := "sms.tencentcloudapi.com"
	service := "sms"
	action := "SendSms"
	version := "2021-01-11"
	region := "ap-guangzhou"
	timestamp := time.Now().Unix()

	// Build canonical request
	httpMethod := "POST"
	canonicalURI := "/"
	canonicalQueryString := ""
	contentType := "application/json; charset=utf-8"
	canonicalHeaders := fmt.Sprintf("content-type:%s\nhost:%s\nx-tc-action:%s\n",
		strings.ToLower(contentType), host, strings.ToLower(action))
	signedHeaders := "content-type;host;x-tc-action"
	hashedPayload := sha256Hex(string(body))
	canonicalRequest := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
		httpMethod, canonicalURI, canonicalQueryString,
		canonicalHeaders, signedHeaders, hashedPayload)

	// Build string to sign
	date := time.Unix(timestamp, 0).UTC().Format("2006-01-02")
	credentialScope := fmt.Sprintf("%s/%s/tc3_request", date, service)
	hashedCanonicalRequest := sha256Hex(canonicalRequest)
	stringToSign := fmt.Sprintf("TC3-HMAC-SHA256\n%d\n%s\n%s",
		timestamp, credentialScope, hashedCanonicalRequest)

	// Calculate signature
	secretDate := hmacSHA256("TC3"+c.cfg.SecretKey, date)
	secretService := hmacSHA256(secretDate, service)
	secretSigning := hmacSHA256(secretService, "tc3_request")
	signature := hex.EncodeToString([]byte(hmacSHA256(secretSigning, stringToSign)))

	// Build authorization header
	authorization := fmt.Sprintf("TC3-HMAC-SHA256 Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		c.cfg.SecretID, credentialScope, signedHeaders, signature)

	url := fmt.Sprintf("https://%s", host)
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Host", host)
	req.Header.Set("X-TC-Action", action)
	req.Header.Set("X-TC-Version", version)
	req.Header.Set("X-TC-Timestamp", fmt.Sprintf("%d", timestamp))
	req.Header.Set("X-TC-Region", region)
	req.Header.Set("Authorization", authorization)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	log.Printf("Tencent SMS response: %s", string(respBody))

	var result struct {
		Response struct {
			Error struct {
				Code    string `json:"Code"`
				Message string `json:"Message"`
			} `json:"Error"`
			SendStatusSet []struct {
				Code         string `json:"Code"`
				PhoneNumber  string `json:"PhoneNumber"`
			} `json:"SendStatusSet"`
		} `json:"Response"`
	}
	json.Unmarshal(respBody, &result)

	if result.Response.Error.Code != "" {
		return fmt.Errorf("tencent SMS error: %s - %s", result.Response.Error.Code, result.Response.Error.Message)
	}

	return nil
}

func (c *Client) sendAliyunSMS(phone, templateID string, params map[string]string) error {
	// Aliyun SMS implementation placeholder
	log.Printf("Aliyun SMS not implemented yet, skip sending to %s", phone)
	return nil
}

func sha256Hex(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func hmacSHA256(key, data string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return string(h.Sum(nil))
}

func (c *Client) SendAppointmentSuccess(phone, doctorName, date, timePeriod string) error {
	period := "上午"
	if timePeriod == "AFTERNOON" {
		period = "下午"
	}
	return c.SendSMS(phone, "appointment_success", map[string]string{
		"1": doctorName,
		"2": date,
		"3": period,
	})
}

func (c *Client) SendAppointmentReminder(phone, doctorName, date, timePeriod string) error {
	period := "上午"
	if timePeriod == "AFTERNOON" {
		period = "下午"
	}
	return c.SendSMS(phone, "appointment_reminder", map[string]string{
		"1": doctorName,
		"2": date,
		"3": period,
	})
}

func (c *Client) SendAppointmentCancelled(phone, doctorName, date string) error {
	return c.SendSMS(phone, "appointment_cancelled", map[string]string{
		"1": doctorName,
		"2": date,
	})
}

func sortKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
