package sms

import (
	"log"
)

type Client struct {
	// TODO: Add SMS provider config (e.g., Alibaba Cloud SMS, Tencent Cloud SMS)
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) SendSMS(phone, templateID string, params map[string]string) error {
	// TODO: Integrate with actual SMS provider
	// For now, just log the message
	log.Printf("SMS sent to %s, template: %s, params: %v", phone, templateID, params)
	return nil
}

func (c *Client) SendAppointmentSuccess(phone, doctorName, date, timePeriod string) error {
	return c.SendSMS(phone, "appointment_success", map[string]string{
		"doctor":      doctorName,
		"date":        date,
		"time_period": timePeriod,
	})
}

func (c *Client) SendAppointmentReminder(phone, doctorName, date, timePeriod string) error {
	return c.SendSMS(phone, "appointment_reminder", map[string]string{
		"doctor":      doctorName,
		"date":        date,
		"time_period": timePeriod,
	})
}

func (c *Client) SendAppointmentCancelled(phone, doctorName, date string) error {
	return c.SendSMS(phone, "appointment_cancelled", map[string]string{
		"doctor": doctorName,
		"date":   date,
	})
}
