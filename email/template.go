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

// ReservationTemplate is a rock-solid, cross-client email template.
// Built specifically to prevent layout breakage across Gmail Web/App, Apple Mail, iCloud Mail, and Outlook:
// 1. Fixed-width outer table (580px) with align="center" to prevent Gmail 100% viewport stretching.
// 2. Pure table-based layout with explicit column widths (no awkward wrapping of labels and colons).
// 3. Solid HEX color codes instead of rgba() which is stripped by Gmail and older mobile clients.
// 4. Double-layer bgcolor attributes on all body, wrapper, and card cells to maintain dark aesthetics in all email readers.
// 5. Anti-spam compliant hidden preheader with non-breaking zero-width spaces.
// 6. Inline typography with standard web-safe fallbacks.
const ReservationTemplate = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" lang="id">
<head>
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <meta http-equiv="X-UA-Compatible" content="IE=edge" />
  <meta name="format-detection" content="telephone=no, date=no, address=no, email=no" />
  <title>Konfirmasi Reservasi - The Wedding of Andri &amp; Cica</title>
  <style type="text/css">
    body, table, td, a { -webkit-text-size-adjust: 100%; -ms-text-size-adjust: 100%; }
    table, td { mso-table-lspace: 0pt; mso-table-rspace: 0pt; }
    img { -ms-interpolation-mode: bicubic; border: 0; outline: none; text-decoration: none; }
    body { margin: 0 !important; padding: 0 !important; width: 100% !important; background-color: #0b0c0e !important; }
    a[x-apple-data-detectors] {
      color: inherit !important;
      text-decoration: none !important;
      font-size: inherit !important;
      font-family: inherit !important;
      font-weight: inherit !important;
      line-height: inherit !important;
    }
    @media only screen and (max-width: 620px) {
      .responsive-table {
        width: 100% !important;
        max-width: 100% !important;
      }
      .mobile-padding {
        padding-left: 16px !important;
        padding-right: 16px !important;
      }
      .mobile-code {
        font-size: 24px !important;
        letter-spacing: 4px !important;
      }
      .mobile-title {
        font-size: 26px !important;
      }
      .detail-label {
        width: 105px !important;
      }
    }
  </style>
</head>
<body bgcolor="#0b0c0e" style="margin: 0; padding: 0; background-color: #0b0c0e; font-family: Arial, Helvetica, sans-serif; -webkit-font-smoothing: antialiased;">

  <!-- PREHEADER (Hidden preview snippet without leaking into email body) -->
  <div style="display: none; max-height: 0px; overflow: hidden; font-size: 1px; line-height: 1px; max-width: 0px; opacity: 0; mso-hide: all;">
    Konfirmasi kehadiran pernikahan Andri &amp; Cica untuk {{.Name}}. Kode Reservasi Anda: {{.ReservationCode}}. Terima kasih atas konfirmasi kehadiran Anda.
    &nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;
  </div>

  <!-- OUTER FULL-WIDTH TABLE -->
  <table role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%" bgcolor="#0b0c0e" style="background-color: #0b0c0e; margin: 0; padding: 0;">
    <tr>
      <td align="center" valign="top" bgcolor="#0b0c0e" style="padding: 24px 10px 36px 10px;">

        <!--[if (gte mso 9)|(IE)]>
        <table align="center" border="0" cellspacing="0" cellpadding="0" width="580">
        <tr>
        <td align="center" valign="top" width="580">
        <![endif]-->

        <!-- MAIN CARD TABLE (Strictly 580px width) -->
        <table role="presentation" border="0" cellpadding="0" cellspacing="0" width="580" align="center" class="responsive-table" bgcolor="#14171f" style="width: 580px; max-width: 580px; background-color: #14171f; border-radius: 12px; border: 1px solid #242934; overflow: hidden; margin: 0 auto;">
          
          <!-- TOP GOLD ACCENT BAR -->
          <tr>
            <td bgcolor="#c5a059" height="3" style="background-color: #c5a059; height: 3px; font-size: 0px; line-height: 0px;">&nbsp;</td>
          </tr>

          <!-- HEADER SECTION -->
          <tr>
            <td align="center" class="mobile-padding" style="padding: 34px 28px 18px 28px; text-align: center;">
              <table role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%">
                <tr>
                  <td align="center" style="font-family: Arial, Helvetica, sans-serif; font-size: 11px; font-weight: bold; letter-spacing: 3px; text-transform: uppercase; color: #c5a059; padding-bottom: 8px;">
                    THE WEDDING OF
                  </td>
                </tr>
                <tr>
                  <td align="center" class="mobile-title" style="font-family: Georgia, 'Times New Roman', serif; font-size: 32px; font-weight: normal; color: #ffffff; letter-spacing: 1px; line-height: 1.2; padding-bottom: 8px;">
                    {{.CoupleNames}}
                  </td>
                </tr>
                <tr>
                  <td align="center" style="font-family: Arial, Helvetica, sans-serif; font-size: 12px; font-weight: bold; letter-spacing: 2px; text-transform: uppercase; color: #94a3b8; padding-bottom: 18px;">
                    {{.EventDay}}, {{.EventDate}}
                  </td>
                </tr>
                <tr>
                  <td align="center">
                    <table role="presentation" border="0" cellpadding="0" cellspacing="0" width="70" align="center">
                      <tr>
                        <td bgcolor="#363c4a" height="1" style="background-color: #363c4a; height: 1px; font-size: 0px; line-height: 0px;">&nbsp;</td>
                      </tr>
                    </table>
                  </td>
                </tr>
              </table>
            </td>
          </tr>

          <!-- HOLY VERSE / QUOTE -->
          <tr>
            <td class="mobile-padding" style="padding: 0 28px 22px 28px;">
              <table role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%" bgcolor="#191d26" style="background-color: #191d26; border: 1px solid #242934; border-radius: 8px;">
                <tr>
                  <td align="center" style="padding: 16px 20px; text-align: center;">
                    <p style="margin: 0 0 6px 0; font-family: Arial, Helvetica, sans-serif; font-size: 11px; font-weight: bold; letter-spacing: 2px; text-transform: uppercase; color: #c5a059;">
                      Q.S. AR-RUM : 21
                    </p>
                    <p style="margin: 0; font-family: Georgia, 'Times New Roman', serif; font-size: 12px; line-height: 1.7; color: #a0aec0; font-style: italic;">
                      &ldquo;Dan di antara tanda-tanda (kebesaran)-Nya ialah Dia menciptakan pasangan-pasangan untukmu dari jenismu sendiri, agar kamu cenderung dan merasa tenteram kepadanya, dan Dia menjadikan di antaramu rasa kasih dan sayang.&rdquo;
                    </p>
                  </td>
                </tr>
              </table>
            </td>
          </tr>

          <!-- GREETING & CONFIRMATION NOTE -->
          <tr>
            <td class="mobile-padding" style="padding: 0 28px 22px 28px; font-family: Arial, Helvetica, sans-serif; font-size: 14px; line-height: 1.6; color: #e2e8f0; text-align: left;">
              <p style="margin: 0 0 10px 0; font-size: 14px; color: #e2e8f0;">
                Kepada Yth. Bapak/Ibu/Saudara/i <strong style="color: #ffffff;">{{.Name}}</strong>,
              </p>
              <p style="margin: 0; font-size: 13px; line-height: 1.6; color: #cbd5e1;">
                Terima kasih atas konfirmasi kehadiran Anda. Merupakan suatu kehormatan dan kebahagiaan bagi kami sekeluarga apabila Anda berkenan hadir untuk memberikan doa restu secara langsung pada hari pernikahan kami.
              </p>
            </td>
          </tr>

          <!-- RESERVATION PASS / CODE BOX -->
          <tr>
            <td class="mobile-padding" style="padding: 0 28px 22px 28px;">
              <table role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%" bgcolor="#1c202a" style="background-color: #1c202a; border: 1px solid #3d3522; border-radius: 10px;">
                <tr>
                  <td align="center" style="padding: 22px 18px; text-align: center;">
                    <p style="margin: 0 0 8px 0; font-family: Arial, Helvetica, sans-serif; font-size: 11px; font-weight: bold; letter-spacing: 2px; text-transform: uppercase; color: #c5a059;">
                      KODE CHECK-IN MASUK
                    </p>
                    
                    <!-- PIN Voucher Box -->
                    <table role="presentation" border="0" cellpadding="0" cellspacing="0" align="center" style="margin: 0 auto 8px auto;">
                      <tr>
                        <td align="center" bgcolor="#0b0c0e" style="background-color: #0b0c0e; border: 1px dashed #c5a059; border-radius: 8px; padding: 10px 26px;">
                          <span class="mobile-code" style="font-family: 'Courier New', Courier, monospace; font-size: 28px; font-weight: bold; letter-spacing: 6px; color: #f6d289;">
                            {{.ReservationCode}}
                          </span>
                        </td>
                      </tr>
                    </table>

                    <p style="margin: 0 0 14px 0; font-family: Arial, Helvetica, sans-serif; font-size: 11px; color: #94a3b8; line-height: 1.4;">
                      Masukkan 5 digit kode di atas saat Anda melakukan scan barcode di lokasi acara
                    </p>

                    <!-- Divider -->
                    <table role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%">
                      <tr>
                        <td bgcolor="#282f3d" height="1" style="background-color: #282f3d; height: 1px; font-size: 0px; line-height: 0px;">&nbsp;</td>
                      </tr>
                    </table>

                    <!-- Status & Quota Details -->
                    <table role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%" style="margin-top: 12px; font-family: Arial, Helvetica, sans-serif;">
                      <tr>
                        <td width="50%" align="center" style="padding: 4px 10px; border-right: 1px solid #282f3d;">
                          <span style="color: #8e9ba8; display: block; font-size: 10px; text-transform: uppercase; letter-spacing: 1px; margin-bottom: 2px;">
                            Status RSVP
                          </span>
                          <span style="color: #4ade80; font-weight: bold; font-size: 12px;">
                            &#10003; Terdaftar Hadir
                          </span>
                        </td>
                        <td width="50%" align="center" style="padding: 4px 10px;">
                          <span style="color: #8e9ba8; display: block; font-size: 10px; text-transform: uppercase; letter-spacing: 1px; margin-bottom: 2px;">
                            Jumlah Kuota
                          </span>
                          <span style="color: #ffffff; font-weight: bold; font-size: 12px;">
                            {{.Guests}} Orang
                          </span>
                        </td>
                      </tr>
                    </table>

                  </td>
                </tr>
              </table>
            </td>
          </tr>

          <!-- EVENT DETAILS SECTION (BULLETPROOF TABLE WITH FIXED COLUMNS) -->
          <tr>
            <td class="mobile-padding" style="padding: 0 28px 24px 28px;">
              <table role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%" bgcolor="#191d26" style="background-color: #191d26; border: 1px solid #242934; border-radius: 8px;">
                <!-- Table Header -->
                <tr>
                  <td colspan="3" bgcolor="#1c202a" style="background-color: #1c202a; padding: 12px 18px; border-bottom: 1px solid #242934; border-top-left-radius: 8px; border-top-right-radius: 8px;">
                    <span style="font-family: Arial, Helvetica, sans-serif; font-size: 11px; font-weight: bold; letter-spacing: 2px; text-transform: uppercase; color: #c5a059;">
                      DETAIL ACARA PERNIKAHAN
                    </span>
                  </td>
                </tr>
                <!-- Rows container -->
                <tr>
                  <td colspan="3" style="padding: 12px 18px 16px 18px;">
                    <table role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%">
                      <!-- Hari & Tanggal -->
                      <tr>
                        <td class="detail-label" width="125" valign="top" style="padding: 6px 0; font-family: Arial, Helvetica, sans-serif; font-size: 13px; color: #94a3b8;">
                          Hari, Tanggal
                        </td>
                        <td width="15" valign="top" style="padding: 6px 0; font-family: Arial, Helvetica, sans-serif; font-size: 13px; color: #94a3b8;">:</td>
                        <td valign="top" style="padding: 6px 0; font-family: Arial, Helvetica, sans-serif; font-size: 13px; font-weight: bold; color: #ffffff;">
                          {{.EventDay}}, {{.EventDate}}
                        </td>
                      </tr>
                      <!-- Akad Nikah -->
                      <tr>
                        <td class="detail-label" width="125" valign="top" style="padding: 6px 0; font-family: Arial, Helvetica, sans-serif; font-size: 13px; color: #94a3b8;">
                          Akad Nikah
                        </td>
                        <td width="15" valign="top" style="padding: 6px 0; font-family: Arial, Helvetica, sans-serif; font-size: 13px; color: #94a3b8;">:</td>
                        <td valign="top" style="padding: 6px 0; font-family: Arial, Helvetica, sans-serif; font-size: 13px; color: #ffffff;">
                          {{.AkadTime}}
                        </td>
                      </tr>
                      <!-- Resepsi -->
                      <tr>
                        <td class="detail-label" width="125" valign="top" style="padding: 6px 0; font-family: Arial, Helvetica, sans-serif; font-size: 13px; color: #94a3b8;">
                          Resepsi
                        </td>
                        <td width="15" valign="top" style="padding: 6px 0; font-family: Arial, Helvetica, sans-serif; font-size: 13px; color: #94a3b8;">:</td>
                        <td valign="top" style="padding: 6px 0; font-family: Arial, Helvetica, sans-serif; font-size: 13px; color: #ffffff;">
                          {{.ResepsiTime}}
                        </td>
                      </tr>
                      <!-- Lokasi Acara -->
                      <tr>
                        <td class="detail-label" width="125" valign="top" style="padding: 6px 0; font-family: Arial, Helvetica, sans-serif; font-size: 13px; color: #94a3b8;">
                          Lokasi Acara
                        </td>
                        <td width="15" valign="top" style="padding: 6px 0; font-family: Arial, Helvetica, sans-serif; font-size: 13px; color: #94a3b8;">:</td>
                        <td valign="top" style="padding: 6px 0; font-family: Arial, Helvetica, sans-serif; font-size: 13px; font-weight: bold; color: #ffffff;">
                          {{.VenueName}}
                        </td>
                      </tr>
                      <!-- Alamat Lengkap -->
                      <tr>
                        <td class="detail-label" width="125" valign="top" style="padding: 6px 0; font-family: Arial, Helvetica, sans-serif; font-size: 13px; color: #94a3b8;">
                          Alamat Lengkap
                        </td>
                        <td width="15" valign="top" style="padding: 6px 0; font-family: Arial, Helvetica, sans-serif; font-size: 13px; color: #94a3b8;">:</td>
                        <td valign="top" style="padding: 6px 0; font-family: Arial, Helvetica, sans-serif; font-size: 13px; line-height: 1.5; color: #cbd5e1;">
                          {{.VenueAddress}}
                        </td>
                      </tr>
                    </table>
                  </td>
                </tr>
              </table>
            </td>
          </tr>

          <!-- INSTRUCTION NOTE -->
          <tr>
            <td class="mobile-padding" align="center" style="padding: 0 28px 22px 28px; text-align: center;">
              <p style="margin: 0; font-family: Arial, Helvetica, sans-serif; font-size: 12px; color: #8e9ba8; line-height: 1.5;">
                Harap simpan email ini atau catat kode reservasi Anda. Tunjukkan kode ini kepada penerima tamu saat tiba di lokasi acara.
              </p>
            </td>
          </tr>

          <!-- CTA BUTTON -->
          <tr>
            <td class="mobile-padding" align="center" style="padding: 0 28px 32px 28px; text-align: center;">
              <table role="presentation" border="0" cellpadding="0" cellspacing="0" align="center" style="margin: 0 auto;">
                <tr>
                  <td align="center" bgcolor="#c5a059" style="border-radius: 24px; background-color: #c5a059;">
                    <a href="{{.ReservationDetailURL}}" target="_blank" rel="noopener noreferrer" style="display: inline-block; padding: 13px 30px; font-family: Arial, Helvetica, sans-serif; font-size: 12px; font-weight: bold; color: #0b0c0e; text-decoration: none; text-transform: uppercase; letter-spacing: 2px; border-radius: 24px;">
                      BUKA UNDANGAN DIGITAL
                    </a>
                  </td>
                </tr>
              </table>
              <p style="margin: 14px 0 0 0; font-family: Arial, Helvetica, sans-serif; font-size: 11px; color: #64748b;">
                Atau lihat petunjuk arah di 
                <a href="{{.MapsURL}}" target="_blank" rel="noopener noreferrer" style="color: #c5a059; text-decoration: underline;">
                  Google Maps
                </a>
              </p>
            </td>
          </tr>

          <!-- FOOTER -->
          <tr>
            <td class="mobile-padding" bgcolor="#0e1015" style="background-color: #0e1015; border-top: 1px solid #1f232d; padding: 22px 28px; text-align: center; font-family: Arial, Helvetica, sans-serif; color: #64748b; font-size: 11px; line-height: 1.6;">
              <p style="margin: 0 0 6px 0; color: #8e9ba8; font-weight: bold;">
                The Wedding of {{.BrideName}} &amp; {{.GroomName}}
              </p>
              <p style="margin: 0 0 8px 0; color: #64748b;">
                Email konfirmasi otomatis dari sistem RSVP website resmi 
                <a href="{{.ReservationDetailURL}}" target="_blank" style="color: #8e9ba8; text-decoration: underline;">andricica.mohaproject.tech</a>.
              </p>
              <p style="margin: 0; font-size: 10px; color: #475569;">
                &copy; {{.Year}} mohaproject.tech &middot; All rights reserved.
              </p>
            </td>
          </tr>

        </table>
        <!-- END MAIN CARD TABLE -->

        <!--[if (gte mso 9)|(IE)]>
        </td>
        </tr>
        </table>
        <![endif]-->

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
	sb.WriteString("TIKET CHECK-IN MASUK\n")
	sb.WriteString("------------------------------------------------------------\n")
	sb.WriteString(fmt.Sprintf("Kode Check-in   : %s\n", data.ReservationCode))
	sb.WriteString("(Masukkan 5 digit kode ini saat scan barcode di lokasi acara)\n\n")
	sb.WriteString(fmt.Sprintf("Status RSVP     : %s\n", data.ReservationStatus))
	sb.WriteString(fmt.Sprintf("Jumlah Kuota    : %d Orang\n\n", data.Guests))

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