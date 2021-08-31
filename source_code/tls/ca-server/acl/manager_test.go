package acl

import (
	"os"
	"testing"
)

func TestACLManager(t *testing.T)  {
	os.Setenv("MY_NODE_NAME", "local-node")
	directory := "D:\\workspace\\gocode\\gomodule\\data"
	config := TunnelACLConfig{
		TLSCAFile:       	directory + "\\rootCA.crt",
		TLSCertFile:        directory + "\\server.crt",
		TLSPrivateKeyFile:  directory + "\\server.key",
		//Token:             "76228ccd639dfad66251cbef3879e70dcdaee29699487d19f3339962a5ec6439.eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzAzOTk0NzV9.jvmQFIzGMFXM9Dx5cPPre2xVa5phT4jkUpwYkF4xHeM",
		Token: "f28610da4b1e7bb3617d8154d02688beb68ccb78099701629079509b08aa7cf4.eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzA0NjgxNDl9.sDDc2UcKSdM3H1NvjXYF3dba0srehSQul_ZNWh8wyJc",
		//HTTPServer:        "https://119.8.58.38:10002",
		HTTPServer:        "https://127.0.0.1:10002",
	}
	manager := NewACLManager(config)
	manager.Start()
}
