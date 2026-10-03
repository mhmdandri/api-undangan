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
		"https://andricica.mohaproject.tech?code=88219",
		"Q.S. AR-RUM : 21",
		"KODE RESERVASI RESMI",
		"Terkonfirmasi Hadir",
		"2 Orang",
		"BUKA UNDANGAN DIGITAL",
	}

	for _, sub := range requiredSubstrings {
		if !strings.Contains(html, sub) {
			t.Errorf("Expected HTML email to contain %q, but it was missing", sub)
		}
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
		"KODE RESERVASI RESMI",
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
		"https://andricica.mohaproject.tech?code=88219",
		"mohaproject.tech",
	}

	for _, sub := range requiredSubstrings {
		if !strings.Contains(text, sub) {
			t.Errorf("Expected PlainText to contain %q, but it was missing", sub)
		}
	}
}
