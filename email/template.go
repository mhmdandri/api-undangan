package email

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
)

// WeddingEmailData holds all template variables for the invitation email
type WeddingEmailData struct {
	Name                 string
	Email                string
	Guests               int
	ReservationStatus    string
	EventDay             string
	EventDate            string
	EventTime            string
	AkadTime             string
	ResepsiTime          string
	VenueName            string
	VenueAddress         string
	ReservationCode      string
	ReservationDetailURL string
	MapsURL              string
	BrideName            string
	GroomName            string
	CoupleNames          string
	Year                 int
}

// DefaultWeddingEmailData creates a ready-to-use dataset populated with the official wedding details
func DefaultWeddingEmailData(name, recipientEmail, code string, guests int) WeddingEmailData {
	if guests <= 0 {
		guests = 1
	}
	if guests > 3 {
		guests = 3
	}
	trimmedName := strings.TrimSpace(name)
	if trimmedName == "" {
		trimmedName = "Tamu Undangan"
	}

	detailURL := "https://andricica.mohaproject.tech"
	if code != "" {
		detailURL = fmt.Sprintf("https://andricica.mohaproject.tech?code=%s", code)
	}

	return WeddingEmailData{
		Name:                 trimmedName,
		Email:                recipientEmail,
		Guests:               guests,
		ReservationStatus:    "Terkonfirmasi Hadir",
		EventDay:             "Sabtu",
		EventDate:            "21 November 2026",
		EventTime:            "09.00 WIB – Selesai",
		AkadTime:             "09.00 WIB",
		ResepsiTime:          "11.00 WIB – Selesai",
		VenueName:            "Turi Jaya Gang IV",
		VenueAddress:         "Jl. Turi Jaya Gang IV No 1, Sagara Makmur, Kec. Tarumajaya, Kab. Bekasi, Jawa Barat",
		ReservationCode:      code,
		ReservationDetailURL: detailURL,
		MapsURL:              "https://maps.google.com/?q=Jl.+Turi+Jaya+Gang+IV+No+1+Sagara+Makmur+Kec.+Tarumajaya+Kab.+Bekasi",
		BrideName:            "Cica Purwanti, S.Pd. Gr.",
		GroomName:            "Muhamad Andriyansyah, S.Kom",
		CoupleNames:          "Andri & Cica",
		Year:                 2026,
	}
}

// Main Wedding Reservation HTML Template
// Designed with modern dark luxury aesthetics matching andricica.mohaproject.tech
// Adheres strictly to Gmail and general email client anti-spam requirements:
// - Multipart compatibility with matching plain text
// - Hidden preheader with zero-width non-breaking spacers
// - Table-based responsive layout with inline CSS
// - WCAG AA compliant high contrast colors
// - Legitimate sender identification and unsubscribe/reasoning footnote
const ReservationTemplate = `<!DOCTYPE html>
<html lang="id" xmlns="http://www.w3.org/1999/xhtml" xmlns:v="urn:schemas-microsoft-com:vml" xmlns:o="urn:schemas-microsoft-com:office:office">
<head>
  <meta charset="UTF-8">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta name="format-detection" content="telephone=no, date=no, address=no, email=no, url=no">
  <meta name="x-apple-disable-message-reformatting">
  <meta name="color-scheme" content="dark light">
  <meta name="supported-color-schemes" content="dark light">
  <title>Konfirmasi Reservasi - The Wedding of Andri &amp; Cica</title>
  <!--[if mso]>
  <noscript>
    <xml>
      <o:OfficeDocumentSettings>
        <o:PixelsPerInch>96</o:PixelsPerInch>
      </o:OfficeDocumentSettings>
    </xml>
  </noscript>
  <![endif]-->
  <style>
    /* Client specific reset */
    body, table, td, a { -webkit-text-size-adjust: 100%; -ms-text-size-adjust: 100%; }
    table, td { mso-table-lspace: 0pt; mso-table-rspace: 0pt; }
    img { -ms-interpolation-mode: bicubic; border: 0; outline: none; text-decoration: none; }
    table { border-collapse: collapse !important; }
    body { margin: 0 !important; padding: 0 !important; width: 100% !important; height: 100% !important; background-color: #0b0c0e; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; }
    /* iOS blue link fix */
    a[x-apple-data-detectors] {
      color: inherit !important;
      text-decoration: none !important;
      font-size: inherit !important;
      font-family: inherit !important;
      font-weight: inherit !important;
      line-height: inherit !important;
    }
    /* Mobile responsive */
    @media only screen and (max-width: 600px) {
      .email-container { width: 100% !important; max-width: 100% !important; margin: auto !important; }
      .content-padding { padding-left: 20px !important; padding-right: 20px !important; }
      .header-title { font-size: 28px !important; }
      .code-display { font-size: 24px !important; letter-spacing: 4px !important; }
      .mobile-stack { display: block !important; width: 100% !important; }
    }
  </style>
</head>
<body style="margin: 0; padding: 0; background-color: #0b0c0e; color: #e2e8f0; -webkit-font-smoothing: antialiased;">

  <!-- PREHEADER (Hidden snippet for inbox preview without leaking into content) -->
  <div style="display: none; max-height: 0px; overflow: hidden; font-size: 1px; line-height: 1px; max-width: 0px; opacity: 0; mso-hide: all;">
    Konfirmasi kehadiran pernikahan Andri &amp; Cica untuk {{.Name}}. Kode Reservasi Anda: {{.ReservationCode}}. Terima kasih atas konfirmasi kehadiran Anda.
    &nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;
  </div>

  <!-- WRAPPER TABLE -->
  <table role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%" style="background-color: #0b0c0e; table-layout: fixed;">
    <tr>
      <td align="center" style="padding: 28px 12px 36px 12px;">

        <!-- MAIN CARD (Max width 580px) -->
        <table role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%" class="email-container" style="max-width: 580px; background-color: #13161c; border-radius: 12px; border: 1px solid #242934; overflow: hidden; box-shadow: 0 10px 30px rgba(0, 0, 0, 0.5);">
          
          <!-- TOP GOLD ACCENT BAR -->
          <tr>
            <td style="background-color: #c5a059; height: 3px; font-size: 0px; line-height: 0px;">&nbsp;</td>
          </tr>

          <!-- HEADER SECTION -->
          <tr>
            <td align="center" class="content-padding" style="padding: 36px 32px 20px 32px; text-align: center;">
              <p style="margin: 0 0 10px 0; font-size: 11px; letter-spacing: 4px; text-transform: uppercase; color: #c5a059; font-weight: 600;">
                THE WEDDING OF
              </p>
              <h1 class="header-title" style="margin: 0 0 8px 0; font-family: 'Playfair Display', Georgia, 'Times New Roman', serif; font-size: 34px; font-weight: 400; color: #ffffff; letter-spacing: 1px; line-height: 1.2;">
                {{.CoupleNames}}
              </h1>
              <p style="margin: 0; font-size: 12px; letter-spacing: 3px; text-transform: uppercase; color: #94a3b8; font-weight: 500;">
                {{.EventDay}}, {{.EventDate}}
              </p>

              <!-- SUBTLE DECORATIVE LINE -->
              <table role="presentation" border="0" cellpadding="0" cellspacing="0" width="80" style="margin: 20px auto 0 auto;">
                <tr>
                  <td style="border-bottom: 1px solid #363c4a; height: 1px; font-size: 0px; line-height: 0px;">&nbsp;</td>
                </tr>
              </table>
            </td>
          </tr>

          <!-- HOLY VERSE / QUOTE -->
          <tr>
            <td class="content-padding" style="padding: 0 32px 24px 32px; text-align: center;">
              <table role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%" style="background-color: #171b23; border: 1px solid #242934; border-radius: 8px;">
                <tr>
                  <td style="padding: 16px 20px; text-align: center;">
                    <p style="margin: 0 0 6px 0; font-size: 11px; letter-spacing: 2px; text-transform: uppercase; color: #c5a059; font-weight: 600;">
                      Q.S. AR-RUM : 21
                    </p>
                    <p style="margin: 0; font-size: 12px; line-height: 1.7; color: #a0aec0; font-style: italic;">
                      &ldquo;Dan di antara tanda-tanda (kebesaran)-Nya ialah Dia menciptakan pasangan-pasangan untukmu dari jenismu sendiri, agar kamu cenderung dan merasa tenteram kepadanya, dan Dia menjadikan di antaramu rasa kasih dan sayang.&rdquo;
                    </p>
                  </td>
                </tr>
              </table>
            </td>
          </tr>

          <!-- GREETING & CONFIRMATION NOTE -->
          <tr>
            <td class="content-padding" style="padding: 0 32px 24px 32px; text-align: left; color: #e2e8f0; font-size: 14px; line-height: 1.65;">
              <p style="margin: 0 0 12px 0;">
                Kepada Yth. Bapak/Ibu/Saudara/i <strong style="color: #ffffff;">{{.Name}}</strong>,
              </p>
              <p style="margin: 0; color: #cbd5e1;">
                Terima kasih atas konfirmasi kehadiran Anda. Merupakan suatu kehormatan dan kebahagiaan bagi kami sekeluarga apabila Anda berkenan hadir untuk memberikan doa restu secara langsung pada hari bahagia kami.
              </p>
            </td>
          </tr>

          <!-- RESERVATION PASS / CODE BOX -->
          <tr>
            <td class="content-padding" style="padding: 0 32px 24px 32px;">
              <table role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%" style="background-color: #191e28; border: 1px solid #383120; border-radius: 10px; text-align: center;">
                <tr>
                  <td style="padding: 22px 20px;">
                    <p style="margin: 0 0 8px 0; font-size: 11px; letter-spacing: 3px; text-transform: uppercase; color: #c5a059; font-weight: 600;">
                      KODE RESERVASI RESMI
                    </p>
                    <div class="code-display" style="font-family: 'Courier New', Courier, monospace; font-size: 30px; font-weight: bold; letter-spacing: 6px; color: #f6d289; margin: 4px 0 12px 0;">
                      {{.ReservationCode}}
                    </div>
                    <table role="presentation" border="0" cellpadding="0" cellspacing="0" align="center" style="margin: 0 auto 10px auto;">
                      <tr>
                        <td style="background-color: rgba(34, 197, 94, 0.15); border: 1px solid rgba(34, 197, 94, 0.4); border-radius: 9999px; padding: 4px 14px; font-size: 11px; font-weight: 600; color: #4ade80; text-transform: uppercase; letter-spacing: 1px;">
                          &#10003; {{.ReservationStatus}}
                        </td>
                      </tr>
                    </table>
                    <p style="margin: 6px 0 0 0; font-size: 12px; color: #94a3b8;">
                      Jumlah Reservasi: <strong style="color: #ffffff;">{{.Guests}} Orang</strong>
                    </p>
                  </td>
                </tr>
              </table>
            </td>
          </tr>

          <!-- EVENT DETAILS SECTION -->
          <tr>
            <td class="content-padding" style="padding: 0 32px 28px 32px;">
              <table role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%" style="background-color: #171b23; border: 1px solid #242934; border-radius: 10px; font-size: 13px;">
                <!-- Header Title -->
                <tr>
                  <td colspan="2" style="padding: 16px 20px 12px 20px; border-bottom: 1px solid #242934;">
                    <p style="margin: 0; font-size: 12px; letter-spacing: 2px; text-transform: uppercase; color: #c5a059; font-weight: 600;">
                      WAKTU &amp; TEMPAT ACARA
                    </p>
                  </td>
                </tr>
                <!-- Tanggal -->
                <tr>
                  <td width="35%" valign="top" style="padding: 12px 20px 8px 20px; color: #8e9ba8;">Hari &amp; Tanggal</td>
                  <td width="65%" valign="top" style="padding: 12px 20px 8px 0; color: #ffffff; font-weight: 500;">
                    {{.EventDay}}, {{.EventDate}}
                  </td>
                </tr>
                <!-- Akad Nikah -->
                <tr>
                  <td width="35%" valign="top" style="padding: 8px 20px; color: #8e9ba8;">Akad Nikah</td>
                  <td width="65%" valign="top" style="padding: 8px 20px 8px 0; color: #ffffff;">
                    {{.AkadTime}}
                  </td>
                </tr>
                <!-- Resepsi -->
                <tr>
                  <td width="35%" valign="top" style="padding: 8px 20px; color: #8e9ba8;">Resepsi Pernikahan</td>
                  <td width="65%" valign="top" style="padding: 8px 20px 8px 0; color: #ffffff;">
                    {{.ResepsiTime}}
                  </td>
                </tr>
                <!-- Lokasi -->
                <tr>
                  <td width="35%" valign="top" style="padding: 8px 20px; color: #8e9ba8;">Lokasi Acara</td>
                  <td width="65%" valign="top" style="padding: 8px 20px 8px 0; color: #ffffff; font-weight: 500;">
                    {{.VenueName}}
                  </td>
                </tr>
                <!-- Alamat -->
                <tr>
                  <td width="35%" valign="top" style="padding: 8px 20px 16px 20px; color: #8e9ba8;">Alamat</td>
                  <td width="65%" valign="top" style="padding: 8px 20px 16px 0; color: #cbd5e1; line-height: 1.5;">
                    {{.VenueAddress}}
                  </td>
                </tr>
              </table>
            </td>
          </tr>

          <!-- NOTICE NOTE -->
          <tr>
            <td class="content-padding" style="padding: 0 32px 24px 32px; text-align: center;">
              <p style="margin: 0; font-size: 12px; color: #8e9ba8; line-height: 1.6;">
                Harap simpan email ini atau catat kode reservasi Anda. Tunjukkan kode ini kepada penerima tamu saat tiba di lokasi acara.
              </p>
            </td>
          </tr>

          <!-- CTA BUTTON -->
          <tr>
            <td align="center" class="content-padding" style="padding: 0 32px 36px 32px;">
              <table role="presentation" border="0" cellpadding="0" cellspacing="0" style="margin: 0 auto;">
                <tr>
                  <td align="center" style="border-radius: 9999px; background-color: #c5a059;">
                    <a href="{{.ReservationDetailURL}}" target="_blank" rel="noopener noreferrer" style="display: inline-block; padding: 14px 32px; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; font-size: 12px; font-weight: 600; color: #0b0c0e; text-decoration: none; text-transform: uppercase; letter-spacing: 2px; border-radius: 9999px;">
                      BUKA UNDANGAN DIGITAL
                    </a>
                  </td>
                </tr>
              </table>
              <p style="margin: 14px 0 0 0; font-size: 11px; color: #64748b;">
                Atau lihat petunjuk arah di 
                <a href="{{.MapsURL}}" target="_blank" rel="noopener noreferrer" style="color: #c5a059; text-decoration: underline;">
                  Google Maps
                </a>
              </p>
            </td>
          </tr>

          <!-- FOOTER -->
          <tr>
            <td class="content-padding" style="background-color: #0e1015; border-top: 1px solid #1f232d; padding: 24px 32px; text-align: center; color: #64748b; font-size: 11px; line-height: 1.6;">
              <p style="margin: 0 0 6px 0; color: #8e9ba8; font-weight: 500;">
                The Wedding of {{.BrideName}} &amp; {{.GroomName}}
              </p>
              <p style="margin: 0 0 10px 0;">
                Anda menerima email ini karena melakukan konfirmasi kehadiran di website resmi 
                <a href="{{.ReservationDetailURL}}" target="_blank" style="color: #8e9ba8; text-decoration: none;">andricica.mohaproject.tech</a>.
              </p>
              <p style="margin: 0; font-size: 10px; color: #475569;">
                &copy; {{.Year}} mohaproject.tech &middot; All rights reserved.
              </p>
            </td>
          </tr>

        </table>
        <!-- END MAIN CARD -->

      </td>
    </tr>
  </table>

</body>
</html>
`

// BuildWeddingReservationEmail renders the primary HTML email
func BuildWeddingReservationEmail(data WeddingEmailData) (string, error) {
	tpl, err := template.New("wedding_email").Parse(ReservationTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// BuildWeddingReservationPlainText generates the exact text equivalent of the reservation email.
// This is critical for anti-spam filters (Gmail, iCloud) to satisfy MIME multipart/alternative rules.
func BuildWeddingReservationPlainText(data WeddingEmailData) string {
	var sb strings.Builder
	sb.WriteString("============================================================\n")
	sb.WriteString(fmt.Sprintf("THE WEDDING OF %s\n", strings.ToUpper(data.CoupleNames)))
	sb.WriteString(fmt.Sprintf("%s, %s\n", data.EventDay, data.EventDate))
	sb.WriteString("============================================================\n\n")

	sb.WriteString(fmt.Sprintf("Kepada Yth. Bapak/Ibu/Saudara/i %s,\n\n", data.Name))
	sb.WriteString("Terima kasih atas konfirmasi kehadiran Anda untuk merayakan momen bahagia pernikahan kami.\n\n")

	sb.WriteString("------------------------------------------------------------\n")
	sb.WriteString("DETAIL RESERVASI ANDA\n")
	sb.WriteString("------------------------------------------------------------\n")
	sb.WriteString(fmt.Sprintf("Kode Reservasi : %s\n", data.ReservationCode))
	sb.WriteString(fmt.Sprintf("Status         : %s\n", data.ReservationStatus))
	sb.WriteString(fmt.Sprintf("Jumlah Tamu    : %d Orang\n\n", data.Guests))

	sb.WriteString("------------------------------------------------------------\n")
	sb.WriteString("WAKTU & TEMPAT ACARA\n")
	sb.WriteString("------------------------------------------------------------\n")
	sb.WriteString(fmt.Sprintf("Hari & Tanggal : %s, %s\n", data.EventDay, data.EventDate))
	sb.WriteString(fmt.Sprintf("Akad Nikah     : %s\n", data.AkadTime))
	sb.WriteString(fmt.Sprintf("Resepsi        : %s\n", data.ResepsiTime))
	sb.WriteString(fmt.Sprintf("Lokasi Acara   : %s\n", data.VenueName))
	sb.WriteString(fmt.Sprintf("Alamat Lengkap : %s\n\n", data.VenueAddress))

	sb.WriteString("Buka Undangan Digital:\n")
	sb.WriteString(fmt.Sprintf("%s\n\n", data.ReservationDetailURL))

	sb.WriteString("Petunjuk Arah Google Maps:\n")
	sb.WriteString(fmt.Sprintf("%s\n\n", data.MapsURL))

	sb.WriteString("Mohon simpan kode reservasi ini dan tunjukkan kepada penerima tamu saat tiba di lokasi.\n\n")
	sb.WriteString(fmt.Sprintf("Salam hangat,\n%s & %s\n", data.GroomName, data.BrideName))
	sb.WriteString("mohaproject.tech\n\n")
	sb.WriteString("------------------------------------------------------------\n")
	sb.WriteString(fmt.Sprintf("Anda menerima email ini karena telah melakukan RSVP di andricica.mohaproject.tech.\n© %d mohaproject.tech\n", data.Year))
	sb.WriteString("============================================================\n")

	return sb.String()
}