//go:build ignore

package main

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/mdp/qrterminal/v3"
	"github.com/pquerna/otp/totp"
)

const (
	issuer      = "Server"
	accountName = "John Doe"
)

const (
	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorCyan    = "\033[36m"
	bold         = "\033[1m"
	colorFGWhite = "\033[97m"
)

func main() {
	imgPath, key, err := generate(issuer, accountName)
	if err != nil {
		fmt.Println(colorReset+bold+"\nError generating TOTP key:"+colorReset+colorRed, err)
		fmt.Println(colorReset+bold+"\nKey (if any):"+colorReset+bold+colorYellow, key)
		fmt.Println(colorReset+bold+"\nQr Code path (if any):"+colorReset+colorYellow, imgPath)
		fmt.Println(colorReset)
		return
	}
	fmt.Println(colorReset+bold+"\nTOTP Key:"+colorReset+colorGreen, key)
	fmt.Println(colorReset+bold+"\nQr Code Path:"+colorReset+colorCyan, imgPath)
	fmt.Println(colorReset)
}

func generate(issuer, accountName string) (string, string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: accountName,
	})
	if err != nil {
		return "", "", fmt.Errorf("Failed to generate TOTP key: %w", err)
	}

	qrCode, err := qr.Encode(key.URL(), qr.M, qr.Auto)
	if err != nil {
		return "", "", fmt.Errorf("Failed to generate QR code: %w", err)
	}

	qrCode, err = barcode.Scale(qrCode, 200, 200)
	if err != nil {
		return "", "", fmt.Errorf("Failed to scale QR code: %w", err)
	}

	pngFile, err := os.Create("code.png")
	if err != nil {
		return "", "", fmt.Errorf("Failed to create temporary file for QR code: %w", err)
	}
	defer pngFile.Close()

	err = png.Encode(pngFile, qrCode)
	if err != nil {
		return "", "", fmt.Errorf("Failed to encode QR code to PNG: %w", err)
	}

	imgPath, err := filepath.Abs(pngFile.Name())
	if err != nil {
		return "", "", fmt.Errorf("Failed to get absolute path of QR code image: %w", err)
	}

	fmt.Print(colorFGWhite)
	qrterminal.GenerateHalfBlock(key.URL(), qrterminal.L, os.Stdout)
	fmt.Print(colorReset)
	
	return imgPath, key.Secret(), nil
}
