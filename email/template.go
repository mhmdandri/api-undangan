package email

import (
	"bytes"
	"text/template"
)

type WeddingEmailData struct {
	Name           string
	Email          string
	Guests         int
	ReservationStatus   string
	EventDate           string
	EventTime           string
	VenueName           string
	VenueAddress        string
	ReservationCode     string
	ReservationDetailURL string
	BrideName           string
	GroomName           string
	Year                int
}
const ReservationTemplate = `<!DOCTYPE html>
<html lang="id">
  <head>
    <meta charset="UTF-8" />
    <title>Wedding Reservation</title>
  </head>
  <body
    style="margin: 0; padding: 0; font-family: Arial, Helvetica, sans-serif"
  >
    <table
      width="100%"
      cellpadding="0"
      cellspacing="0"
      border="0"
      style="padding: 24px 0"
    >
      <tr>
        <td align="center">
          <table
            width="480"
            cellpadding="0"
            cellspacing="0"
            border="0"
            style="
              background-color: #141414;
              border-radius: 6px;
              overflow: hidden;
            "
          >
            <!-- Header -->
            <tr>
              <td
                align="center"
                style="
                  background-color: #e50914;
                  padding: 18px 12px;
                  color: #ffffff;
                "
              >
                <h2 style="margin: 0; font-size: 20px">Konfirmasi Kehadiran</h2>
                <span style="font-size: 12px">
                  {{.BrideName}} &amp; {{.GroomName}}
                </span>
              </td>
            </tr>

            <!-- Body -->
            <tr>
              <td style="padding: 20px; color: #ffffff">
                <p style="margin: 0 0 10px; font-size: 13px">
                  Halo <strong>{{.Name}}</strong>,
                </p>
                <p
                  style="
                    margin: 0 0 14px;
                    font-size: 13px;
                    color: #b3b3b3;
                    line-height: 1.6;
                  "
                >
                  Terimasih atas konfirmasi kehadiran Anda untuk acara
                  pernikahan kami. Jangan lupa untuk catat tanggal nya ya.
                </p>

                <table
                  width="100%"
                  cellpadding="0"
                  cellspacing="0"
                  border="0"
                  style="font-size: 12px; color: #ffffff"
                >
                  <tr>
                    <td style="padding: 4px 0; width: 35%; color: #b3b3b3">
                      Tanggal
                    </td>
                    <td style="padding: 4px 0">: {{.EventDate}}</td>
                  </tr>
                  <tr>
                    <td style="padding: 4px 0; color: #b3b3b3">Waktu</td>
                    <td style="padding: 4px 0">: {{.EventTime}}</td>
                  </tr>
                  <tr>
                    <td style="padding: 4px 0; color: #b3b3b3">Lokasi</td>
                    <td style="padding: 4px 0">: {{.VenueName}}</td>
                  </tr>
                  <tr>
                    <td style="padding: 4px 0; color: #b3b3b3">Alamat</td>
                    <td style="padding: 4px 0">: {{.VenueAddress}}</td>
                  </tr>
                  <tr>
                    <td style="padding: 4px 0; color: #b3b3b3">Kode</td>
                    <td style="padding: 4px 0">
                      : <strong>{{.ReservationCode}}</strong>
                    </td>
                  </tr>
                </table>

                <p
                  style="
                    margin: 14px 0 0;
                    font-size: 12px;
                    color: #b3b3b3;
                    line-height: 1.5;
                  "
                >
                  Mohon simpan kode ini dan tunjukkan saat hadir di lokasi.
                </p>
              </td>
            </tr>

            <!-- Button -->
            <tr>
              <td align="center" style="padding: 0 20px 18px 20px">
                <a
                  href="{{.ReservationDetailURL}}"
                  style="
                    display: inline-block;
                    padding: 10px 22px;
                    background-color: #e50914;
                    color: #ffffff;
                    text-decoration: none;
                    border-radius: 4px;
                    font-size: 13px;
                    font-weight: bold;
                  "
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  Lihat Detail
                </a>
              </td>
            </tr>

            <!-- Footer -->
            <tr>
              <td
                align="center"
                style="
                  padding: 10px 20px 14px 20px;
                  font-size: 10px;
                  color: #808080;
                  background-color: #141414;
                "
              >
                &copy; mohaproject {{.Year}} <br /> Wedding of {{.BrideName}} &amp;
                {{.GroomName}}
              </td>
            </tr>
          </table>
        </td>
      </tr>
    </table>
  </body>
</html>
`


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