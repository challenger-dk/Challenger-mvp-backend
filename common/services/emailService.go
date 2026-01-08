package services

import (
	"context"
	"fmt"
	"log/slog"
	"server/common/config"
	"sync"

	"github.com/mrz1836/postmark"
)

var (
	postmarkClient *postmark.Client
	postmarkOnce   sync.Once
)

// getPostmarkClient returns a singleton Postmark client instance
func getPostmarkClient() *postmark.Client {
	postmarkOnce.Do(func() {
		postmarkClient = postmark.NewClient(config.AppConfig.PostmarkAPIKey, "")
	})
	return postmarkClient
}

func SendPasswordResetEmail(email, resetCode string) error {
	client := getPostmarkClient()

	htmlBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
		</head>
		<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px;">
			<div style="background-color: #f4f4f4; padding: 30px; border-radius: 5px;">
				<h2 style="color: #333; margin-top: 0; text-align: center;">Password Reset Request</h2>
				<p style="text-align: center;">You requested to reset your password. Enter this code in the app:</p>
				<div style="text-align: center; margin: 30px 0;">
					<div style="background-color: #fff; border: 2px solid #007bff; border-radius: 8px; padding: 20px; display: inline-block;">
						<div style="font-size: 36px; font-weight: bold; letter-spacing: 8px; color: #007bff; font-family: 'Courier New', monospace;">%s</div>
					</div>
				</div>
				<p style="text-align: center; color: #666; font-size: 14px;">This code will expire in 1 hour.</p>
				<p style="text-align: center; color: #999; font-size: 12px; margin-top: 30px;">If you didn't request this, please ignore this email.</p>
			</div>
		</body>
		</html>
	`, resetCode)

	textBody := fmt.Sprintf(`
Password Reset Request

You requested to reset your password. Enter this code in the app:

%s

This code will expire in 1 hour. If you didn't request this, please ignore this email.
	`, resetCode)

	emailMessage := postmark.Email{
		From:          config.AppConfig.PostmarkFromEmail,
		To:            email,
		Subject:       "Reset Your Password",
		HTMLBody:      htmlBody,
		TextBody:      textBody,
		MessageStream: "outbound",
		Tag:           "password-reset",
	}

	ctx := context.Background()
	response, err := client.SendEmail(ctx, emailMessage)
	if err != nil {
		slog.Error("Failed to send password reset email", "error", err, "to", email)
		return fmt.Errorf("failed to send email: %w", err)
	}

	slog.Info("Password reset email sent successfully", "to", email, "message_id", response.MessageID)
	return nil
}
