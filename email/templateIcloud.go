package email

// BuildWeddingReservationEmailIcloud renders the invitation email for iCloud / Apple Mail recipients.
// It uses the same bulletproof dark luxury template as BuildWeddingReservationEmail to ensure
// iCloud recipients receive the complete, full-featured invitation experience rather than plain text.
func BuildWeddingReservationEmailIcloud(data WeddingEmailData) (string, error) {
	return BuildWeddingReservationEmail(data)
}