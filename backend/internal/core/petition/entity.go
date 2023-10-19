package petition

import (
	"io"
)

// Define a struct to hold the data
type PetitionData struct {
	FileName          string        `json:"file_name"`
	SheetNumber       *int          `json:"sheet_number,omitempty"`  //get it from potsgreSQL
	CreationDate      string        `json:"creation_date,omitempty"` //getted from back
	Location          string        `json:"location"`
	ResponsiblePerson string        `json:"responsible_person"`
	Questions         []Question    `json:"questions"`
	OwnerName         string        `json:"owner_name"`
	OwnerAddress      string        `json:"owner_address"`
	PdfData           []byte        `json:"pdf_data,omitempty"` //we did't get it from frontend we give it in backend
	File              io.ReadCloser `json:"file"`
}

// Define a struct to represent each question
type Question struct {
	Number   int    `json:"number"`
	Text     string `json:"text"`
	Decision string `json:"description"` // For, Abstain, Against За Воздержусь, Против
}

const TemplateHTML = `<!DOCTYPE html>
<html>

<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title></title>
</head>

<body>
    <p style='color: #000000;line-height: 100%;text-align: left;margin-bottom: 0in;background: transparent;font-family: "Times New Roman", serif;font-size:11px;'><span size="2" style="font-size:15px;">&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; Приложение 3</span></p>
    <p style='color: rgb(0, 0, 0);text-align: center;margin: 0.04in 0.2in 0in 2.93in;background: transparent;font-family: "Times New Roman", serif;font-size:11px;line-height: 115%;text-indent: 0in;'><span size="2" style="font-size:15px;">к Правилам принятия решений по управлению&nbsp;объектом&nbsp;кондоминиума&nbsp;и&nbsp;содержанию&nbsp;общего&nbsp;имущества&nbsp;объекта&nbsp;кондоминиума</span></p>
    <p style='color: rgb(0, 0, 0);line-height: 100%;text-align: center;margin-bottom: 0in;background: transparent;font-family: "Times New Roman", serif;font-size:11px;margin-left: 0.05in;margin-top: 0.15in;'><span size="2" style="font-size:15px;"><strong>Лист</strong><strong>&nbsp;</strong><strong>№</strong><strong>&nbsp;</strong><u>&nbsp; {{.SheetNumber}}</u></span></p>
    <p style="color: rgb(0, 0, 0);text-align: center;margin: 0.05in 0.1in 0in 0.05in;background: transparent;line-height: 120%;"><span size="4" style="font-size:19px;"><strong>голосования&nbsp;при&nbsp;проведении&nbsp;письменного&nbsp;опроса&nbsp;собственников&nbsp;квартир,&nbsp;нежилых&nbsp;помещений</strong></span></p>
    <p style='color: #000000;line-height: 100%;text-align: left;margin-bottom: 0in;background: transparent;font-family: "Times New Roman", serif;font-size:11px;margin-top: 0.01in;'><br></p>
    <p style='color: #000000;line-height: 100%;text-align: left;margin-bottom: 0in;background: transparent;font-family: "Times New Roman", serif;font-size:11px;margin-left: 3.42in;margin-top: 0.06in;'><span size="2" style="font-size:15px;">&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &laquo;__&raquo; {{.CreationDate}} года</span></p>
    <p style='color: #000000;line-height: 100%;text-align: left;margin-bottom: 0in;background: transparent;font-family: "Times New Roman", serif;font-size:11px;margin-left: 0.57in;margin-top: 0.17in;'><br></p>
    <p style='color: #000000;line-height: 100%;text-align: left;margin-bottom: 0in;background: transparent;font-family: "Times New Roman", serif;font-size:11px;margin-left: 0.57in;margin-top: 0.17in;'><span size="2" style="font-size:15px;">время <u>__</u></span></p>
    <p style='color: #000000;line-height: 100%;text-align: left;margin-bottom: 0in;background: transparent;font-family: "Times New Roman", serif;font-size:11px;margin-left: 0.57in;margin-top: 0.12in;'><span size="2" style="font-size:15px;">Местонахождение&nbsp;многоквартирного&nbsp;жилого&nbsp;дома:</span></p>
    <p style='color: #000000;line-height: 100%;text-align: left;margin-bottom: 0in;background: transparent;font-family: "Times New Roman", serif;font-size:11px;margin-top: 0in;'>&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; {{.Location}}</p>
    <p style='color: #000000;line-height: 100%;text-align: left;margin-bottom: 0in;background: transparent;font-family: "Times New Roman", serif;font-size:11px;margin-left: 0.57in;margin-top: 0.1in;'><span size="2" style="font-size:15px;">Ответственные&nbsp;лица:</span></p>
    <p style='color: #000000;line-height: 100%;text-align: left;margin-bottom: 0in;background: transparent;font-family: "Times New Roman", serif;font-size:11px;margin-top: 0in;'>&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; {{.ResponsiblePerson}} </p>
    <p style='color: #000000;line-height: 100%;text-align: left;margin-bottom: 0in;background: transparent;font-family: "Times New Roman", serif;font-size:11px;'><br></p>
    <p style='color: #000000;line-height: 100%;text-align: left;margin-bottom: 0in;background: transparent;font-family: "Times New Roman", serif;font-size:11px;margin-top: 0in;'><br></p>
    <p style='color: #000000;line-height: 100%;text-align: left;margin-bottom: 0in;background: transparent;font-family: "Times New Roman", serif;font-size:11px;'><br></p>
    <p style='color: #000000;line-height: 100%;text-align: left;margin-bottom: 0in;background: transparent;font-family: "Times New Roman", serif;font-size:11px;margin-top: 0in;'><br></p>
    <p style='color: #000000;line-height: 100%;text-align: left;margin-bottom: 0in;background: transparent;font-family: "Times New Roman", serif;font-size:11px;margin-left: 0.57in;margin-top: 0.1in;'><span size="2" style="font-size:15px;"><em>(назначаемые</em><em>&nbsp;</em><em>из</em><em>&nbsp;</em><em>числа</em><em>&nbsp;</em><em>собственников</em><em>&nbsp;</em><em>квартир,</em><em>&nbsp;</em><em>нежилого</em><em>&nbsp;</em><em>помещения)</em></span></p>
    <p style='color: #000000;line-height: 100%;text-align: left;margin-bottom: 0in;background: transparent;font-family: "Times New Roman", serif;font-size:11px;margin-top: 0in;'><br></p>
    <table style="width: 5.0e+2pt;" cellpadding="1">
        <tbody>
            <tr>
                <td rowspan="2" style="border-top: 1.00pt solid #000000;border-bottom: 1px solid #000000;border-left: 1px solid #000000;border-right: 1px solid #000000;padding-top: 0in;padding-bottom: 0in;padding-left: 0in;padding-right: 0in;">
                    <p style="color: #000000;line-height: 100%;text-align: center;margin-bottom: 0in;background: transparent;margin-left: 0in;margin-top: 0.03in;">№</p>
                </td>
                <td rowspan="2" style="border: 1px solid #000000;">
                    <p style="color: #000000;line-height: 100%;text-align: left;margin-bottom: 0in;background: transparent;margin-left: 0.96in;margin-right: 0.65in;text-indent: -0.3in;margin-top: 0.03in;"><span size="2" style="font-size:15px;">Вопросы, внесенные для обсуждения:</span></p>
                </td>
                <td colspan="3" style="border-top: 1.00pt solid #000000;border-bottom: 1.00pt solid #000000;border-left: 1px solid #000000;border-right: 1px solid #000000;padding-top: 0in;padding-bottom: 0in;padding-left: 0in;padding-right: 0in;">
                    <p style="color: #000000;line-height: 100%;text-align: center;margin-bottom: 0in;background: transparent;margin-left: 1.7in;margin-right: 1.7in;margin-top: 0.03in;"><span size="2" style="font-size:15px;">Голосую</span></p>
                </td>
            </tr>
            <tr>
                <td style="border-top: 1.00pt solid #000000;border-bottom: 1px solid #000000;border-left: 1px solid #000000;border-right: 1px solid #000000;padding-top: 0in;padding-bottom: 0in;padding-left: 0in;padding-right: 0in;">
                    <p style="color: rgb(0, 0, 0);line-height: 100%;text-align: center;margin: 0.03in 0.31in 0in;background: transparent;"><span size="2" style="font-size:15px;">&laquo;За&raquo;</span></p>
                    <p style="color: rgb(0, 0, 0);line-height: 100%;text-align: center;margin: 0.07in 0.31in 0in;background: transparent;"><span size="2" style="font-size:15px;">(подпись)</span></p>
                </td>
                <td style="border-top: 1.00pt solid #000000;border-bottom: 1px solid #000000;border-left: 1px solid #000000;border-right: 1px solid #000000;padding-top: 0in;padding-bottom: 0in;padding-left: 0in;padding-right: 0in;">
                    <p style="color: #000000;line-height: 100%;text-align: left;margin-bottom: 0in;background: transparent;margin-left: 0.33in;margin-top: 0.03in;"><span size="2" style="font-size:15px;">&laquo;Против&raquo;</span></p>
                    <p style="color: #000000;line-height: 100%;text-align: left;margin-bottom: 0in;background: transparent;margin-left: 0.33in;margin-top: 0.07in;"><span size="2" style="font-size:15px;">(подпись)</span></p>
                </td>
                <td style="border-top: 1.00pt solid #000000;border-bottom: 1px solid #000000;border-left: 1px solid #000000;border-right: 1px solid #000000;padding-top: 0in;padding-bottom: 0in;padding-left: 0in;padding-right: 0in;">
                    <p style="color: rgb(0, 0, 0);line-height: 100%;text-align: center;margin: 0.03in 0.41in 0in 0.42in;background: transparent;"><span size="2" style="font-size:15px;">&laquo;Воздержусь&raquo;</span></p>
                    <p style="color: rgb(0, 0, 0);line-height: 100%;text-align: center;margin: 0.07in 0.41in 0in 0.42in;background: transparent;"><span size="2" style="font-size:15px;">(подпись)</span></p>
                </td>
            </tr>
			{{range .Questions}}
            <tr>
                <td style="border-top: 1px solid #000000;border-bottom: 1px solid #000000;border-left: 1px solid #000000;border-right: 1px solid #000000;padding-top: 0in;padding-bottom: 0in;padding-left: 0in;padding-right: 0in;">
                    <p style="color: #000000;line-height: 100%;text-align: left;margin-bottom: 0in;background: transparent;margin-left: 0.08in;margin-top: 0.04in;"><span size="2" style="font-size:15px;">{{.Number}}</span></p>
                </td>
                <td style="border: 1px solid #000000;">
                    <p style="color: #000000;line-height: 100%;text-align: left;margin-bottom: 0in;background: transparent;">{{.Text}}</p>
                </td>
                <td style="border: 1px solid #000000;">
                    <p style="color: #000000;line-height: 100%;text-align: left;margin-bottom: 0in;background: transparent;">{{if eq .Decision "За"}}X{{end}}</p>
                </td>
                <td style="border: 1px solid #000000;">
                    <p style="color: #000000;line-height: 100%;text-align: left;margin-bottom: 0in;background: transparent;">{{if eq .Decision "Против"}}X{{end}}</p>
                </td>
                <td style="border-top: 1px solid #000000;border-bottom: 1px solid #000000;border-left: 1px solid #000000;border-right: 1px solid #000000;padding-top: 0in;padding-bottom: 0in;padding-left: 0in;padding-right: 0in;">
                    <p style="color: #000000;line-height: 100%;text-align: left;margin-bottom: 0in;background: transparent;">{{if eq .Decision "Воздержусь"}}X{{end}}</p>
                </td>
            </tr>
			{{end}}
        </tbody>
    </table>
    <p style='color: #000000;line-height: 100%;text-align: left;margin-bottom: 0in;background: transparent;font-family: "Times New Roman", serif;font-size:11px;margin-left: 0.08in;margin-top: 0.07in;'>Ф.И.О.&nbsp;собственника&nbsp;квартиры,&nbsp;нежилого&nbsp;помещения<u>&nbsp;{{.OwnerName}}</u></p>
    <p style='color: #000000;line-height: 115%;text-align: left;margin-bottom: 0in;background: transparent;font-family: "Times New Roman", serif;font-size:11px;margin-left: 0.08in;margin-right: 2.59in;margin-top: 0.02in;'><u>&nbsp;\</u> Адрес собственника квартиры, нежилого помещения<u>&nbsp;{{.OwnerAddress}}</u></p>
    <p style='color: #000000;line-height: 100%;text-align: left;margin-bottom: 0in;background: transparent;font-family: "Times New Roman", serif;font-size:11px;'><br></p>
    <p style='color: #000000;line-height: 115%;text-align: left;margin-bottom: 0in;background: transparent;font-family: "Times New Roman", serif;font-size:11px;margin-left: 0.73in;margin-right: 2.58in;text-indent: -0.65in;margin-top: 0in;'>Подпись (собственника&nbsp;квартиры,&nbsp;нежилого&nbsp;помещения)__</p>
    <p style='color: #000000;line-height: 115%;text-align: left;margin-bottom: 0in;background: transparent;font-family: "Times New Roman", serif;font-size:11px;margin-left: 0.91in;margin-right: 2.58in;text-indent: -0.83in;margin-top: 0in;'>Подпись (ответственные лица)__</p>
    <p style='color: #000000;line-height: 115%;text-align: left;margin-bottom: 0in;background: transparent;font-family: "Times New Roman", serif;font-size:11px;margin-left: 0.97in;margin-right: 2.58in;text-indent: -0.89in;margin-top: 0in;'>Подпись (члена&nbsp;совета&nbsp;дома)__</p>
    <p style='color: #000000;line-height: 115%;text-align: left;margin-bottom: 0in;background: transparent;font-family: "Times New Roman", serif;font-size:11px;margin-left: 0.97in;margin-right: 2.58in;text-indent: -0.89in;margin-top: 0in;'>Подпись (члена&nbsp;совета&nbsp;дома)__</p>
    <p style='color: #000000;line-height: 115%;text-align: left;margin-bottom: 0in;background: transparent;font-family: "Times New Roman", serif;font-size:11px;margin-left: 0.97in;margin-right: 2.58in;text-indent: -0.89in;margin-top: 0in;'>Подпись (члена совета дома)__</p>
</body>

</html>`

/*
const TemplateHTMLOptimized = `<!DOCTYPE html>
<html lang="ru,en">

<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title></title>
    <style>
        body {
            font-family: "Times New Roman", serif;
            font-size: 11px;
            color: #000;
            line-height: 100%;
            margin: 0;
            background: transparent;
        }

        p {
            text-align: left;
            margin: 0;
        }

        .center {
            text-align: center;
            font-size: 15px;
        }

        .strong {
            font-weight: bold;
        }

        .indent {
            margin-left: 0.57in;
        }

        .table {
            border-collapse: collapse;
            width: 500pt;
            margin: 0 auto;
        }

        th, td {
            border: 1px solid #000;
            padding: 0.03in;
            text-align: left;
        }

        .checkbox {
            text-align: center;
        }
    </style>
</head>

<body>
    <p class="center strong">Приложение 3</p>
    <p class="center">к Правилам принятия решений по управлению объектом кондоминиума и содержанию общего имущества объекта кондоминиума</p>
    <p class="center strong">Лист № <u>{{.SheetNumber}}</u></p>
    <p class="center strong">голосования при проведении письменного опроса собственников квартир, нежилых помещений</p>
    <br>
    <p class="indent strong">&laquo;__&raquo; {{.CreationDate}} года</p>
    <br>
    <p class="indent strong">время <u>__</u></p>
    <p>Местонахождение многоквартирного жилого дома: {{.Location}}</p>
    <p>Ответственные лица: {{.ResponsiblePerson}}</p>
    <br><br><br><br>
    <p class="indent strong"><em>(назначаемые из числа собственников квартир, нежилого помещения)</em></p>
    <br>

    <table class="table" cellpadding="1">
        <thead>
            <tr>
                <th style="width: 60pt;">№</th>
                <th style="width: 200pt;">Вопросы, внесенные для обсуждения:</th>
                <th style="width: 100pt;">&laquo;За&raquo;</th>
                <th style="width: 100pt;">&laquo;Против&raquo;</th>
                <th style="width: 100pt;">&laquo;Воздержусь&raquo;</th>
            </tr>
        </thead>
        <tbody>
            {{range .Questions}}
            <tr>
                <td class="checkbox">{{if eq .Decision "За"}}X{{end}}</td>
                <td>{{.Text}}</td>
                <td class="checkbox">{{if eq .Decision "Против"}}X{{end}}</td>
                <td class="checkbox">{{if eq .Decision "Воздержусь"}}X{{end}}</td>
            </tr>
            {{end}}
        </tbody>
    </table>

    <p class="indent">Ф.И.О. собственника квартиры, нежилого помещения <u>{{.OwnerName}}</u></p>
    <p class="indent"><u> Адрес собственника квартиры, нежилого помещения <u>{{.OwnerAddress}}</u></u></p>
    <br>
    <p class="indent">Подпись (собственника квартиры, нежилого помещения)__</p>
    <p class="indent">Подпись (ответственные лица)__</p>
    <p class="indent">Подпись (члена совета дома)__</p>
    <p class="indent">Подпись (члена совета дома)__</p>
    <p class="indent">Подпись (члена совета дома)__</p>
</body>

</html>`
*/
