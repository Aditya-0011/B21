package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"proxy-server/totp"
	"proxy-server/utils"
	"time"
)

type downloadRequest struct {
	URL string `json:"url"`
	OTP string `json:"otp"`
}

type logRequest struct {
	OTP string `json:"otp"`
}

type Controller struct {
	Log     *utils.Logger
	LogFile string
	Secret string
}

func (c *Controller) GetData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		c.Log.Entry(fmt.Sprintf("[INVALID METHOD] %s request from IP %s", r.Method, r.RemoteAddr))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req downloadRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	allowed, err := (&totp.Auth{OTP: req.OTP, Secret: c.Secret}).Validate()
	if err != nil || !allowed {
		c.Log.Entry(fmt.Sprintf("[AUTH FAIL] Download attempt for %s from IP %s", req.URL, r.RemoteAddr))
		http.Error(w, "Invalid OTP", http.StatusForbidden)
		return
	}

	c.Log.Entry(fmt.Sprintf("[PROXY] Starting download: %s (IP: %s)", req.URL, r.RemoteAddr))

	resp, err := http.Get(req.URL)
	if err != nil {
		c.Log.Entry(fmt.Sprintf("[ERROR] Upstream failure: %v", err))
		http.Error(w, "Failed to reach target", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.Header().Set("Content-Length", resp.Header.Get("Content-Length"))
	w.Header().Set("Content-Disposition", "attachment; filename=downloaded_file")

	n, err := io.Copy(w, resp.Body)
	if err != nil {
		c.Log.Entry(fmt.Sprintf("[ERROR] Stream interrupted after %d bytes: %v", n, err))
	} else {
		c.Log.Entry(fmt.Sprintf("[SUCCESS] Transferred %d bytes", n))
	}
}

func (c *Controller) GetLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		c.Log.Entry(fmt.Sprintf("[INVALID METHOD] %s request from IP %s", r.Method, r.RemoteAddr))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req logRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	allowed, err := (&totp.Auth{OTP: req.OTP, Secret: c.Secret}).Validate()
	if err != nil || !allowed {
		if err.Error() == "Secret not found" {
			c.Log.Entry("[ERROR] Secret not configured")
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		c.Log.Entry(fmt.Sprintf("[AUTH FAIL] Log access attempt from IP %s", r.RemoteAddr))
		http.Error(w, "Invalid OTP", http.StatusForbidden)
		return
	}

	tempFilename := fmt.Sprintf("logs_snapshot_%d.tmp", time.Now().Unix())
	utils.LogMutex.Lock()

	sourceFile, err := os.Open(c.LogFile)
	if err != nil {
		utils.LogMutex.Unlock()
		http.Error(w, "Log file missing", http.StatusInternalServerError)
		return
	}

	destFile, err := os.Create(tempFilename)
	if err != nil {
		sourceFile.Close()
		utils.LogMutex.Unlock()
		http.Error(w, "Server Disk Error", http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(destFile, sourceFile)

	sourceFile.Close()
	destFile.Close()
	utils.LogMutex.Unlock()

	if err != nil {
		http.Error(w, "Error snapshotting logs", http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFilename)

	c.Log.Entry(fmt.Sprintf("[ADMIN] Logs exported to IP %s", r.RemoteAddr))

	http.ServeFile(w, r, tempFilename)
}
