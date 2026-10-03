package email

import (
	"bytes"
	"html/template"
)

// ReservationIcloudTemplate is tuned specifically for Apple Mail and iCloud Mail clients.
// It uses clean Apple WebKit styles, system font fallbacks, solid table styling,
// and full compliance with Apple's Mail Privacy Protection and spam heuristics.
const ReservationIcloudTemplate = `<!DOCTYPE html>
<html lang="id">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <meta name="format-detection" content="telephone=no, date=no, address=no, email=no, url=no">
  <meta name="x-apple-disable-message-reformatting">
  <meta name="color-scheme" content="dark light">
  <meta name="supported-color-schemes" content="dark light">
  <title>Konfirmasi Reservasi - The Wedding of Andri &amp; Cica</title>
  <style>
    body, table, td, a { -webkit-text-size-adjust: 100%; -ms-text-size-adjust: 100%; }
    table, td { mso-table-lspace: 0pt; mso-table-rspace: 0pt; }
    img { -ms-interpolation-mode: bicubic; border: 0; outline: none; }
    body { margin: 0; padding: 0; width: 100% !important; background-color: #0b0c0e; font-family: -apple-system, BlinkMacSystemFont, 'SF Pro Text', 'Helvetica Neue', Helvetica, Arial, sans-serif; -webkit-font-smoothing: antialiased; }
    a[x-apple-data-detectors] {
      color: inherit !important;
      text-decoration: none !important;
      font-size: inherit !important;
      font-family: inherit !important;
      font-weight: inherit !important;
      line-height: inherit !important;
    }
    @media only screen and (max-width: 600px) {
      .email-container { width: 100% !important; max-width: 100% !important; }
      .content-padding { padding-left: 20px !important; padding-right: 20px !important; }
      .header-title { font-size: 28px !important; }
      .code-display { font-size: 24px !important; letter-spacing: 4px !important; }
    }
  </style>
</head>
<body style="margin:0; padding:0; background-color: #0b0c0e; color: #e2e8f0;">

  <!-- Apple Mail Preheader -->
  <div style="display:none; font-size:1px; color:#0b0c0e; line-height:1px; max-height:0px; max-width:0px; opacity:0; overflow:hidden;">
    Konfirmasi kehadiran pernikahan Andri &amp; Cica untuk {{.Name}}. Kode reservasi: {{.ReservationCode}}.
    &nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;
  </div>

  <table role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%" style="background-color: #0b0c0e; table-layout: fixed;">
    <tr>
      <td align="center" style="padding: 24px 12px 32px 12px;">

        <table role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%" class="email-container" style="max-width: 560px; background-color: #13161c; border-radius: 12px; border: 1px solid #242934; overflow: hidden;">
          
          <!-- Top Accent -->
          <tr>
            <td style="background-color: #c5a059; height: 3px; font-size: 0px; line-height: 0px;">&nbsp;</td>
          </tr>

          <!-- Header -->
          <tr>
            <td align="center" class="content-padding" style="padding: 32px 28px 20px 28px; text-align: center;">
              <p style="margin: 0 0 8px 0; font-size: 11px; letter-spacing: 4px; text-transform: uppercase; color: #c5a059; font-weight: 600;">
                THE WEDDING OF
              </p>
              <h1 class="header-title" style="margin: 0 0 6px 0; font-family: 'Playfair Display', Georgia, 'Times New Roman', serif; font-size: 32px; font-weight: 400; color: #ffffff; letter-spacing: 1px; line-height: 1.2;">
                {{.CoupleNames}}
              </h1>
              <p style="margin: 0; font-size: 12px; letter-spacing: 2px; text-transform: uppercase; color: #94a3b8;">
                {{.EventDay}}, {{.EventDate}}
              </p>
            </td>
          </tr>

          <!-- Greeting -->
          <tr>
            <td class="content-padding" style="padding: 0 28px 20px 28px; font-size: 14px; line-height: 1.6; color: #e2e8f0;">
              <p style="margin: 0 0 10px 0;">
                Kepada Yth. <strong style="color: #ffffff;">{{.Name}}</strong>,
              </p>
              <p style="margin: 0; color: #cbd5e1;">
                Terima kasih telah mengonfirmasi kehadiran Anda pada acara pernikahan kami. Kami sangat berbahagia dan menantikan kehadiran serta doa restu Anda.
              </p>
            </td>
          </tr>

          <!-- Reservation Code Card -->
          <tr>
            <td class="content-padding" style="padding: 0 28px 22px 28px;">
              <table role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%" style="background-color: #191e28; border: 1px solid #383120; border-radius: 10px; text-align: center;">
                <tr>
                  <td style="padding: 20px;">
                    <p style="margin: 0 0 6px 0; font-size: 11px; letter-spacing: 2px; text-transform: uppercase; color: #c5a059; font-weight: 600;">
                      KODE RESERVASI RESMI
                    </p>
                    <div class="code-display" style="font-family: 'Courier New', Courier, monospace; font-size: 28px; font-weight: bold; letter-spacing: 6px; color: #f6d289; margin: 4px 0 10px 0;">
                      {{.ReservationCode}}
                    </div>
                    <table role="presentation" border="0" cellpadding="0" cellspacing="0" align="center" style="margin: 0 auto 8px auto;">
                      <tr>
                        <td style="background-color: rgba(34, 197, 94, 0.15); border: 1px solid rgba(34, 197, 94, 0.4); border-radius: 9999px; padding: 4px 14px; font-size: 11px; font-weight: 600; color: #4ade80; text-transform: uppercase;">
                          &#10003; {{.ReservationStatus}}
                        </td>
                      </tr>
                    </table>
                    <p style="margin: 0; font-size: 12px; color: #94a3b8;">
                      Jumlah Tamu: <strong style="color: #ffffff;">{{.Guests}} Orang</strong>
                    </p>
                  </td>
                </tr>
              </table>
            </td>
          </tr>

          <!-- Event Detail Table -->
          <tr>
            <td class="content-padding" style="padding: 0 28px 24px 28px;">
              <table role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%" style="background-color: #171b23; border: 1px solid #242934; border-radius: 8px; font-size: 13px;">
                <tr>
                  <td colspan="2" style="padding: 14px 18px 10px 18px; border-bottom: 1px solid #242934;">
                    <p style="margin: 0; font-size: 12px; letter-spacing: 2px; text-transform: uppercase; color: #c5a059; font-weight: 600;">
                      INFORMASI ACARA
                    </p>
                  </td>
                </tr>
                <tr>
                  <td width="35%" valign="top" style="padding: 10px 18px 6px 18px; color: #8e9ba8;">Waktu Akad</td>
                  <td width="65%" valign="top" style="padding: 10px 18px 6px 0; color: #ffffff;">{{.AkadTime}}</td>
                </tr>
                <tr>
                  <td width="35%" valign="top" style="padding: 6px 18px; color: #8e9ba8;">Waktu Resepsi</td>
                  <td width="65%" valign="top" style="padding: 6px 18px 6px 0; color: #ffffff;">{{.ResepsiTime}}</td>
                </tr>
                <tr>
                  <td width="35%" valign="top" style="padding: 6px 18px; color: #8e9ba8;">Lokasi</td>
                  <td width="65%" valign="top" style="padding: 6px 18px 6px 0; color: #ffffff; font-weight: 500;">{{.VenueName}}</td>
                </tr>
                <tr>
                  <td width="35%" valign="top" style="padding: 6px 18px 14px 18px; color: #8e9ba8;">Alamat</td>
                  <td width="65%" valign="top" style="padding: 6px 18px 14px 0; color: #cbd5e1; line-height: 1.5;">{{.VenueAddress}}</td>
                </tr>
              </table>
            </td>
          </tr>

          <!-- CTA Button -->
          <tr>
            <td align="center" class="content-padding" style="padding: 0 28px 30px 28px;">
              <table role="presentation" border="0" cellpadding="0" cellspacing="0">
                <tr>
                  <td align="center" style="border-radius: 9999px; background-color: #c5a059;">
                    <a href="{{.ReservationDetailURL}}" target="_blank" rel="noopener noreferrer" style="display: inline-block; padding: 12px 28px; font-size: 12px; font-weight: 600; color: #0b0c0e; text-decoration: none; text-transform: uppercase; letter-spacing: 2px; border-radius: 9999px;">
                      BUKA UNDANGAN
                    </a>
                  </td>
                </tr>
              </table>
              <p style="margin: 12px 0 0 0; font-size: 11px; color: #64748b;">
                <a href="{{.MapsURL}}" target="_blank" rel="noopener noreferrer" style="color: #c5a059; text-decoration: underline;">
                  Petunjuk Arah Google Maps
                </a>
              </p>
            </td>
          </tr>

          <!-- Footer -->
          <tr>
            <td class="content-padding" style="background-color: #0e1015; border-top: 1px solid #1f232d; padding: 20px 28px; text-align: center; color: #64748b; font-size: 11px; line-height: 1.5;">
              <p style="margin: 0 0 4px 0; color: #8e9ba8;">
                Wedding of {{.BrideName}} &amp; {{.GroomName}}
              </p>
              <p style="margin: 0 0 8px 0;">
                Email konfirmasi otomatis dari <a href="{{.ReservationDetailURL}}" target="_blank" style="color: #8e9ba8; text-decoration: none;">andricica.mohaproject.tech</a>.
              </p>
              <p style="margin: 0; font-size: 10px; color: #475569;">
                &copy; {{.Year}} mohaproject.tech
              </p>
            </td>
          </tr>

        </table>

      </td>
    </tr>
  </table>

</body>
</html>
`

// BuildWeddingReservationEmailIcloud renders the iCloud/Apple Mail tailored HTML email
func BuildWeddingReservationEmailIcloud(data WeddingEmailData) (string, error) {
	tpl, err := template.New("wedding_email_icloud").Parse(ReservationIcloudTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}