package wechat

import (
	"bytes"
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"math/rand"
	"net/http"
	"sort"
	"strings"
)

type PayClient struct {
	appID     string
	mchID     string
	apiKey    string
	notifyURL string
}

func NewPayClient(appID, mchID, apiKey, notifyURL string) *PayClient {
	return &PayClient{
		appID:     appID,
		mchID:     mchID,
		apiKey:    apiKey,
		notifyURL: notifyURL,
	}
}

type UnifiedOrderRequest struct {
	XMLName        xml.Name `xml:"xml"`
	AppID          string   `xml:"appid"`
	MchID          string   `xml:"mch_id"`
	NonceStr       string   `xml:"nonce_str"`
	Sign           string   `xml:"sign"`
	Body           string   `xml:"body"`
	OutTradeNo     string   `xml:"out_trade_no"`
	TotalFee       int      `xml:"total_fee"`
	SpbillCreateIP string   `xml:"spbill_create_ip"`
	NotifyURL      string   `xml:"notify_url"`
	TradeType      string   `xml:"trade_type"`
	OpenID         string   `xml:"openid,omitempty"`
}

type UnifiedOrderResponse struct {
	XMLName    xml.Name `xml:"xml"`
	ReturnCode string   `xml:"return_code"`
	ReturnMsg  string   `xml:"return_msg"`
	ResultCode string   `xml:"result_code"`
	PrepayID   string   `xml:"prepay_id"`
	NonceStr   string   `xml:"nonce_str"`
	Sign       string   `xml:"sign"`
}

type PayNotification struct {
	XMLName       xml.Name `xml:"xml"`
	ReturnCode    string   `xml:"return_code"`
	ResultCode    string   `xml:"result_code"`
	OutTradeNo    string   `xml:"out_trade_no"`
	TransactionID string   `xml:"transaction_id"`
	TotalFee      int      `xml:"total_fee"`
	TimeEnd       string   `xml:"time_end"`
}

func (c *PayClient) CreatePrepayOrder(outTradeNo, body, openID string, totalFee int, clientIP string) (string, error) {
	req := &UnifiedOrderRequest{
		AppID:          c.appID,
		MchID:          c.mchID,
		NonceStr:       generateNonceStr(),
		Body:           body,
		OutTradeNo:     outTradeNo,
		TotalFee:       totalFee,
		SpbillCreateIP: clientIP,
		NotifyURL:      c.notifyURL,
		TradeType:      "JSAPI",
		OpenID:         openID,
	}
	req.Sign = c.sign(req)

	xmlData, err := xml.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(
		"https://api.mch.weixin.qq.com/pay/unifiedorder",
		"application/xml",
		bytes.NewBuffer(xmlData),
	)
	if err != nil {
		return "", fmt.Errorf("failed to call wechat pay api: %w", err)
	}
	defer resp.Body.Close()

	var orderResp UnifiedOrderResponse
	if err := xml.NewDecoder(resp.Body).Decode(&orderResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if orderResp.ReturnCode != "SUCCESS" {
		return "", fmt.Errorf("wechat pay error: %s", orderResp.ReturnMsg)
	}
	if orderResp.ResultCode != "SUCCESS" {
		return "", fmt.Errorf("wechat pay order failed: %s", orderResp.ReturnMsg)
	}

	return orderResp.PrepayID, nil
}

func (c *PayClient) VerifyNotification(data []byte) (*PayNotification, error) {
	var notification PayNotification
	if err := xml.Unmarshal(data, &notification); err != nil {
		return nil, fmt.Errorf("failed to unmarshal notification: %w", err)
	}

	if notification.ReturnCode != "SUCCESS" || notification.ResultCode != "SUCCESS" {
		return nil, fmt.Errorf("payment notification failed")
	}

	return &notification, nil
}

func (c *PayClient) sign(params interface{}) string {
	m := structToMap(params)
	m["key"] = c.apiKey

	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buf strings.Builder
	for i, k := range keys {
		if i > 0 {
			buf.WriteString("&")
		}
		buf.WriteString(k)
		buf.WriteString("=")
		buf.WriteString(m[k])
	}

	hash := md5.Sum([]byte(buf.String()))
	return fmt.Sprintf("%X", hash)
}

func structToMap(params interface{}) map[string]string {
	m := make(map[string]string)
	switch v := params.(type) {
	case *UnifiedOrderRequest:
		if v.AppID != "" {
			m["appid"] = v.AppID
		}
		if v.MchID != "" {
			m["mch_id"] = v.MchID
		}
		if v.NonceStr != "" {
			m["nonce_str"] = v.NonceStr
		}
		if v.Body != "" {
			m["body"] = v.Body
		}
		if v.OutTradeNo != "" {
			m["out_trade_no"] = v.OutTradeNo
		}
		m["total_fee"] = fmt.Sprintf("%d", v.TotalFee)
		if v.SpbillCreateIP != "" {
			m["spbill_create_ip"] = v.SpbillCreateIP
		}
		if v.NotifyURL != "" {
			m["notify_url"] = v.NotifyURL
		}
		if v.TradeType != "" {
			m["trade_type"] = v.TradeType
		}
		if v.OpenID != "" {
			m["openid"] = v.OpenID
		}
	}
	return m
}

func generateNonceStr() string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 32)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}
