package integrations

import (
	"context"
	"fmt"
	"time"

	routeros "github.com/go-routeros/routeros"
)

type MikroTikIntegration struct {
	client *routeros.Client
	config config.MikroTikConfig
}

func NewMikroTikIntegration(cfg config.MikroTikConfig) (*MikroTikIntegration, error) {
	address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	client, err := routeros.Dial(address, cfg.Username, cfg.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MikroTik: %v", err)
	}
	
	return &MikroTikIntegration{
		client: client,
		config: cfg,
	}, nil
}

func (m *MikroTikIntegration) AddHotspotUser(username, password string, profile string, timeout time.Duration) error {
	cmd := fmt.Sprintf("/ip/hotspot/user/add name=%s password=%s profile=%s limit-uptime=%s", 
		username, password, profile, timeout.String())
		
	_, err := m.client.RunArgs(strings.Split(cmd, " "))
	return err
}

func (m *MikroTikIntegration) RemoveHotspotUser(username string) error {
	// Find user ID first
	cmd := fmt.Sprintf("/ip/hotspot/user/print where name=%s", username)
	resp, err := m.client.RunArgs(strings.Split(cmd, " "))
	if err != nil {
		return err
	}
	
	if len(resp.Re) == 0 {
		return fmt.Errorf("user not found")
	}
	
	// Extract user ID
	userID := resp.Re[0].Map[".id"]
	
	// Remove user
	cmd = fmt.Sprintf("/ip/hotspot/user/remove .id=%s", userID)
	_, err = m.client.RunArgs(strings.Split(cmd, " "))
	return err
}

func (m *MikroTikIntegration) GetActiveUsers() ([]map[string]string, error) {
	cmd := "/ip/hotspot/active/print"
	resp, err := m.client.RunArgs(strings.Split(cmd, " "))
	if err != nil {
		return nil, err
	}
	
	users := make([]map[string]string, len(resp.Re))
	for i, item := range resp.Re {
		users[i] = item.Map
	}
	
	return users, nil
}

func (m *MikroTikIntegration) GetUserBytes(username string) (int64, int64, error) {
	cmd := fmt.Sprintf("/ip/hotspot/user/print stats where name=%s", username)
	resp, err := m.client.RunArgs(strings.Split(cmd, " "))
	if err != nil {
		return 0, 0, err
	}
	
	if len(resp.Re) == 0 {
		return 0, 0, fmt.Errorf("user not found")
	}
	
	user := resp.Re[0].Map
	bytesIn, _ := strconv.ParseInt(user["bytes-in"], 10, 64)
	bytesOut, _ := strconv.ParseInt(user["bytes-out"], 10, 64)
	
	return bytesIn, bytesOut, nil
}

func (m *MikroTikIntegration) Close() {
	m.client.Close()
}