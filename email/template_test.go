package email

import (
	"strings"
	"testing"
)

func TestBuildWeddingReservationEmail(t *testing.T) {
	data := DefaultWeddingEmailData("Budi Santoso", "budi@example.com", "88219", 2)

	html, err := BuildWeddingReservationEmail(data)
	if err != nil {
		t.Fatalf("BuildWeddingReservationEmail returned error: %v", err)
	}

	requiredSubstrings := []string{
		"Budi Santoso",
		"88219",
		"Andri &amp; Cica",
		"21 November 2026",
		"09.00 WIB",
		"Turi Jaya Gang IV",
		"https://weddingofandricica.me?code=88219",
		"api.qrserver.com/v1/create-qr-code",
		"data=https%3A%2F%2Fweddingofandricica.me%3Fcode%3D88219",
		"Q.S. AR-RUM : 21",
		"KODE CHECK-IN MASUK",
		"Terdaftar Hadir",
		"2 Orang",
		"BUKA UNDANGAN DIGITAL",
	}

	for _, sub := range requiredSubstrings {
		if !strings.Contains(html, sub) {
			t.Errorf("Expected HTML email to contain %q, but it was missing", sub)
		}
	}

	qrIdx := strings.Index(html, "api.qrserver.com/v1/create-qr-code")
	checkinCodeIdx := strings.Index(html, ">88219</span>")
	if qrIdx == -1 || checkinCodeIdx == -1 || qrIdx >= checkinCodeIdx {
		t.Errorf("Expected QR code to appear before reservation code in checkin box, but qrIdx=%d, checkinCodeIdx=%d", qrIdx, checkinCodeIdx)
	}
}

func TestBuildWeddingReservationEmailIcloud(t *testing.T) {
	data := DefaultWeddingEmailData("Sarah Apple", "sarah@icloud.com", "44102", 1)

	html, err := BuildWeddingReservationEmailIcloud(data)
	if err != nil {
		t.Fatalf("BuildWeddingReservationEmailIcloud returned error: %v", err)
	}

	requiredSubstrings := []string{
		"Sarah Apple",
		"44102",
		"Andri &amp; Cica",
		"21 November 2026",
		"Turi Jaya Gang IV",
		"KODE CHECK-IN MASUK",
		"1 Orang",
	}

	for _, sub := range requiredSubstrings {
		if !strings.Contains(html, sub) {
			t.Errorf("Expected iCloud HTML to contain %q, but it was missing", sub)
		}
	}
}

func TestBuildWeddingReservationPlainText(t *testing.T) {
	data := DefaultWeddingEmailData("Budi Santoso", "budi@example.com", "88219", 2)

	text := BuildWeddingReservationPlainText(data)

	requiredSubstrings := []string{
		"THE WEDDING OF ANDRI & CICA",
		"Budi Santoso",
		"88219",
		"21 November 2026",
		"Turi Jaya Gang IV",
		"https://weddingofandricica.me?code=88219",
		"weddingofandricica.me",
	}

	for _, sub := range requiredSubstrings {
		if !strings.Contains(text, sub) {
			t.Errorf("Expected PlainText to contain %q, but it was missing", sub)
		}
	}
}
