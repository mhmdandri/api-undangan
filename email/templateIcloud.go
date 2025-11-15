package email

import (
	"bytes"
	"text/template"
)

const ReservationIcloudTemplate = `
<!DOCTYPE html>
<html lang="id">
  <head>
    <meta charset="UTF-8" />
    <title>Konfirmasi Reservasi Pernikahan</title>
  </head>
  <body style="margin:0; padding:0; font-family: Arial, Helvetica, sans-serif; background-color:#f4f4f4;">
    <table width="100%" cellpadding="0" cellspacing="0" border="0" style="padding:24px 0;">
      <tr>
        <td align="center">
          <table width="520" cellpadding="0" cellspacing="0" border="0" style="background-color:#ffffff; border-radius:6px; overflow:hidden; border:1px solid #e0e0e0;">
            <!-- Header -->
            <tr>
              <td align="left" style="padding:18px 20px; border-bottom:1px solid #eeeeee;">
                <h1 style="margin:0 0 4px; font-size:18px; color:#333333;">
                  Konfirmasi Reservasi Pernikahan
                </h1>
                <p style="margin:0; font-size:12px; color:#777777;">
                  {{.BrideName}} &amp; {{.GroomName}}
                </p>
              </td>
            </tr>

            <!-- Preheader text (help deliverability) -->
            <tr>
              <td style="padding:10px 20px 0 20px; font-size:11px; color:#999999;">
                Terima kasih telah mengonfirmasi kehadiran Anda. Berikut ringkasan detail acara.
              </td>
            </tr>

            <!-- Body -->
            <tr>
              <td style="padding:10px 20px 20px 20px; font-size:13px; color:#333333; line-height:1.6;">
                <p style="margin:0 0 10px;">
                  Halo <strong>{{.Name}}</strong>,
                </p>
                <p style="margin:0 0 14px;">
                  Terima kasih atas konfirmasi kehadiran Anda untuk acara pernikahan kami.
                  Mohon simpan informasi berikut dan catat tanggal acaranya.
                </p>

                <table width="100%" cellpadding="0" cellspacing="0" border="0" style="font-size:13px; color:#333333;">
                  <tr>
                    <td style="padding:4px 0; width:30%; color:#777777;">Tanggal</td>
                    <td style="padding:4px 0;">: {{.EventDate}}</td>
                  </tr>
                  <tr>
                    <td style="padding:4px 0; color:#777777;">Waktu</td>
                    <td style="padding:4px 0;">: {{.EventTime}}</td>
                  </tr>
                  <tr>
                    <td style="padding:4px 0; color:#777777;">Lokasi</td>
                    <td style="padding:4px 0;">: {{.VenueName}}</td>
                  </tr>
                  <tr>
                    <td style="padding:4px 0; color:#777777;">Alamat</td>
                    <td style="padding:4px 0;">: {{.VenueAddress}}</td>
                  </tr>
                  <tr>
                    <td style="padding:4px 0; color:#777777;">Kode Reservasi</td>
                    <td style="padding:4px 0;">
                      : <strong>{{.ReservationCode}}</strong>
                    </td>
                  </tr>
                </table>

                <p style="margin:14px 0 0;">
                  Mohon tunjukkan kode reservasi ini saat hadir di lokasi.
                </p>
              </td>
            </tr>

            <!-- Button -->
            <tr>
              <td align="left" style="padding:0 20px 18px 20px;">
                <a
                  href="{{.ReservationDetailURL}}"
                  style="
                    display:inline-block;
                    padding:10px 20px;
                    background-color:#d9534f;
                    color:#ffffff;
                    text-decoration:none;
                    border-radius:4px;
                    font-size:13px;
                  "
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  Lihat detail reservasi
                </a>
              </td>
            </tr>

            <!-- Footer -->
            <tr>
              <td
                align="left"
                style="
                  padding:12px 20px 16px 20px;
                  font-size:11px;
                  color:#888888;
                  border-top:1px solid #eeeeee;
                  background-color:#fafafa;
                "
              >
                &copy; mohaproject {{.Year}} · Wedding of {{.BrideName}} &amp; {{.GroomName}}<br />
                Anda menerima email ini karena melakukan reservasi kehadiran acara pernikahan kami.
                Jika merasa tidak pernah melakukan reservasi, abaikan saja email ini.
              </td>
            </tr>
          </table>
        </td>
      </tr>
    </table>
  </body>
</html>
`

func BuildWeddingReservationEmailIcloud(data WeddingEmailData) (string, error) {
	tpl, err := template.New("wedding_email").Parse(ReservationIcloudTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}