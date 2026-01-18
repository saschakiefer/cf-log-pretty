/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package parser

import (
	"strings"
	"testing"
)

func TestParseLine_JSONWithStackTrace(t *testing.T) {
	input := `   2023-04-30T08:39:15.71+0200 [APP/PROC/WEB/0] OUT { "written_at":"2023-04-30T06:39:15.716Z","level":"ERROR","logger":"com.sap.shell.services.navigation.client.ChangeRequestMessageListener","msg":"MQ: processChange >> java.lang.IllegalArgumentException: SITEREFERENCES: is not supported. Dropping change for: sitereferences.default MODIFIED","stacktrace":["java.lang.IllegalArgumentException: SITEREFERENCES: is not supported.","\tat com.example.Class.method(Class.java:123)"] }`

	msg, ok := ParseLine(input)
	if !ok {
		t.Fatal("Expected log line to be parsed")
	}
	if msg.Level != "ERROR" {
		t.Errorf("Expected level ERROR, got %s", msg.Level)
	}
	if !strings.Contains(msg.Message, "MQ: processChange") {
		t.Errorf("Unexpected message content: %s", msg.Message)
	}
	if len(msg.StackTrace) < 2 {
		t.Errorf("Expected stacktrace entries, got %v", msg.StackTrace)
	}
}

func TestParseLine_JSONWithoutStackTrace(t *testing.T) {
	input := `   2023-04-30T08:39:16.76+0200 [APP/PROC/WEB/0] OUT { "written_at":"2023-04-30T06:39:16.766Z","level":"INFO","logger":"com.sap.shell.services.navigation.client.ChangeRequestMessageListener","msg":" MQ: Incoming message on topic: cdm/site/entities/deleted" }`

	msg, ok := ParseLine(input)
	if !ok {
		t.Fatal("Expected log line to be parsed")
	}
	if msg.Level != "INFO" {
		t.Errorf("Expected level INFO, got %s", msg.Level)
	}
	if msg.Logger != "com.sap.shell.services.navigation.client.ChangeRequestMessageListener" {
		t.Errorf("Unexpected logger: %s", msg.Logger)
	}
	if !strings.Contains(msg.Message, "Incoming message") {
		t.Errorf("Unexpected message: %s", msg.Message)
	}
	if len(msg.StackTrace) != 0 {
		t.Errorf("Expected no stacktrace, got %v", msg.StackTrace)
	}
}

func TestParseLine_PlainServiceMessage(t *testing.T) {
	input := `   2023-04-30T08:43:34.02+0200 [APP/PROC/WEB/2] OUT [SERVICE: InfluxDB] Couldn't write to server: 404 Not Found`

	msg, ok := ParseLine(input)
	if !ok {
		t.Fatal("Expected plain log to be parsed")
	}
	if !strings.Contains(msg.Message, "InfluxDB") {
		t.Errorf("Unexpected message content: %s", msg.Message)
	}
	if msg.Level != "-----" {
		t.Errorf("Expected no level (-> -----), got %s", msg.Level)
	}
}

func TestParseLine_RouterAccessLog(t *testing.T) {
	input := `   2024-01-20T09:37:58.99+0100 [RTR/6] OUT portal-service.cfapps.eu12.hana.ondemand.com - [2024-01-20T08:37:58.990006873Z] "GET /navigation HTTP/1.1" 200`

	msg, ok := ParseLine(input)
	if !ok {
		t.Fatal("Expected RTR log to be parsed")
	}
	if !strings.Contains(msg.Message, `"GET /navigation`) {
		t.Errorf("Unexpected message: %s", msg.Message)
	}
}

func TestParseLine_Unstructured(t *testing.T) {
	input := "test"

	msg, ok := ParseLine(input)
	if !ok {
		t.Fatal("Expected simple line to be parsed")
	}
	if msg.Message != "test" {
		t.Errorf("Expected message to be 'test', got %s", msg.Message)
	}
}
