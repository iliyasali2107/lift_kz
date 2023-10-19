package petition

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"

	"mado/pkg/validator"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type Repository interface {
	Save(ctx context.Context, dto *PetitionData) (*PetitionData, error)
	GetNextID(ctx context.Context) (*int, error)
	GetPetitionPdfByID(ctx context.Context, pdfID *int) (*PetitionData, error)
}

// Service is a user service interface.
type Service struct {
	petitionRepository Repository
	logger             *zap.Logger
}

// NewService creates a new user service.
func NewService(petitionRepository Repository, logger *zap.Logger) Service {
	return Service{
		petitionRepository: petitionRepository,
		logger:             logger,
	}
}

func (s Service) GeneratePDF(ctx context.Context, t *template.Template, pageData interface{}, outFilePath string, templatePath string) error {
	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	buf := &bytes.Buffer{}
	err := t.Execute(buf, pageData)
	if err != nil {
		return fmt.Errorf("executing template: %w", err)
	}

	html := buf.String()

	htmlAssetsDir := "file://" + templatePath
	fmt.Println(htmlAssetsDir)

	var wg sync.WaitGroup
	wg.Add(1)
	err = chromedp.Run(ctx,
		chromedp.Navigate(htmlAssetsDir),

		// Add EventLoadEventFired listener.
		chromedp.ActionFunc(func(ctx context.Context) error {
			lctx, cancel := context.WithCancel(ctx)
			chromedp.ListenTarget(lctx, func(ev interface{}) {
				if _, ok := ev.(*page.EventLoadEventFired); ok {
					// It's a good habit to remove the event
					// listener if we don't need it anymore.
					wg.Done()
					cancel()
				}
			})
			return nil
		}),

		// Set the page content.
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				fmt.Println(1)
				return err
			}

			err = page.SetDocumentContent(frameTree.Frame.ID, html).Do(ctx)
			if err != nil {
				fmt.Println(2)
				return err
			}

			return nil
		}),

		// Wait until the page is loaded (including its resources).
		chromedp.ActionFunc(func(ctx context.Context) error {
			wg.Wait()
			return nil
		}),

		// Print to PDF.
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().WithPrintBackground(true).
				WithPaperWidth(5.8).WithPaperHeight(8.3).Do(ctx)
			if err != nil {
				fmt.Println(3)
				return err
			}

			err = os.WriteFile(outFilePath, buf, 0644)
			if err != nil {
				fmt.Println(4)
				return err
			}
			return nil
		}),
	)
	if err != nil {
		return fmt.Errorf("chromedp: %w", err)
	}

	return nil
}

// Todo save to db ?
// func (s Service) GeneratePetitionPDF(data *PetitionData) (*PetitionData, error) {

// 	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	// defer cancel()

// 	currentTime := time.Now()
// 	// Format the date as "21.09.2023"
// 	data.CreationDate = currentTime.Format("02.01.2006")

// 	var err error
// 	// data.SheetNumber, err = s.petitionRepository.GetNextID(ctx)
// 	// if err != nil {
// 	// 	s.logger.Error("Failed to get id for next doc: ", zap.Error(err))
// 	// 	return nil, err
// 	// }

// 	// Create a new template It will be used to format the HTML content of the petition.
// 	tmpl, err := template.New("petition").Parse(TemplateHTML)
// 	if err != nil {
// 		s.logger.Error("Failed to parse template: ", zap.Error(err))
// 		return nil, err
// 	}

// 	// Create a buffer to hold the generated HTML content
// 	var htmlContentBuffer strings.Builder
// 	if err := tmpl.Execute(&htmlContentBuffer, data); err != nil {
// 		s.logger.Error("Failed to execute data to buffer: ", zap.Error(err))
// 		return nil, err
// 	}

// 	// he HTML content in the htmlContentBuffer is parsed using the goquery library. This library allows you to manipulate and traverse HTML documents in a structured way
// 	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContentBuffer.String()))
// 	if err != nil {
// 		s.logger.Error("Failed to pass htmlContent to goquery: ", zap.Error(err))
// 		return nil, err
// 	}

// 	// Save the modified HTML to a temporary file
// 	tempHTMLFileName := "temp.html"
// 	tempHTMLFile, err := os.Create(tempHTMLFileName)
// 	if err != nil {
// 		s.logger.Error("Failed to create template file: ", zap.Error(err))
// 		return nil, err
// 	}
// 	defer tempHTMLFile.Close()

// 	// Use goquery to write the modified HTML to the file
// 	htmlContent, err := doc.Html()
// 	if err != nil {
// 		s.logger.Error("Failed to get html content: ", zap.Error(err))
// 		return nil, err
// 	}
// 	_, err = tempHTMLFile.WriteString(htmlContent)
// 	if err != nil {
// 		s.logger.Error("Failed to write the modified HTML to the file: ", zap.Error(err))
// 		return nil, err
// 	}

// 	// // Generate PDF using wkhtmltopdf
// 	// //todo change namingOfFile
// 	// // pdfFileName := "output.pdf"
// 	// cmd := exec.Command("wkhtmltopdf", tempHTMLFileName, data.FileName)
// 	// err = cmd.Run()
// 	// if err != nil {
// 	// 	s.logger.Error("Failed to convert html to pdf: ", zap.Error(err))
// 	// 	return nil, err
// 	// }
// 	// s.logger.Debug("PDF generated: %s\n", zap.String("filename", data.FileName))

// 	var requestBody bytes.Buffer
// 	writer := multipart.NewWriter(&requestBody)
// 	htmlPart, err := writer.CreateFormFile("files", "index.html")
// 	if err != nil {
// 		fmt.Println(err)
// 		return nil, err
// 	}

// 	_, err = io.Copy(htmlPart, tempHTMLFile)
// 	if err != nil {
// 		fmt.Println(err)
// 		return nil, err
// 	}

// 	err = writer.Close()
// 	if err != nil {
// 		fmt.Println(err)
// 		return nil, err
// 	}

// 	url := "http://gotenberg:3000/forms/chromium/convert/html"
// 	req, err := http.NewRequest("POST", url, &requestBody)
// 	if err != nil {
// 		fmt.Println(err)
// 		return nil, err
// 	}

// 	req.Header.Set("Content-Type", writer.FormDataContentType())

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println("error cliend.Do(): ", err)
// 		return nil, err
// 	}

// 	defer resp.Body.Close()

// 	var p []byte

// 	if resp.StatusCode == http.StatusOK {
// 		fmt.Println(resp.Status)
// 		resp.Body.Read(p)
// 		fmt.Println(string(p))
// 		return nil, fmt.Errorf("status is not OK")
// 	}

// 	// Save the resulting PDF to a file.
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	orig, _ := os.Getwd()
// 	os.Chdir("../../")
// 	fmt.Println(os.Getwd())
// 	pdfFile, err := os.Create("internal/core/petition/outputs/my.pdf")
// 	if err != nil {
// 		fmt.Println("Error creating PDF file:", err)
// 		return nil, err
// 	}
// 	defer pdfFile.Close()
// 	os.Chdir(orig)
// 	_, err = io.Copy(pdfFile, resp.Body)
// 	if err != nil {
// 		fmt.Println("Error saving PDF:", err)
// 		return nil, err
// 	}

// 	s.logger.Debug(fmt.Sprintf("file name: %s", data.FileName))

// 	data.File = resp.Body

// 	return data, nil
// 	// // Prepare the PDF file data and metadata
// 	// data.PdfData, err = os.ReadFile(data.FileName)
// 	// if err != nil {
// 	// 	s.logger.Error("Failed to read PDF file: ", zap.Error(err))
// 	// 	return nil, err
// 	// }
// 	// // return nil, nil
// 	// return s.petitionRepository.Save(ctx, data)

// }

func (s Service) GetPetitionPdfByID(id string) (*PetitionData, error) {
	docId, err := validator.IdValidator(id)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.petitionRepository.GetPetitionPdfByID(ctx, docId)
}

func (s Service) GeneratePetitionPDF(data *PetitionData) (*string, error) {
	os.Chdir("../../")
	curr, _ := os.Getwd()
	// /app
	curr = curr + "/internal/core/petition"
	fmt.Println(os.ReadDir(curr))
	url := "http://localhost:3000/forms/chromium/convert/html"
	filePath := "internal/core/petition/index.html"
	outputPath := "internal/core/petition/my2.pdf"

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("OPEN FILE", err)
		return nil, err
	}
	defer file.Close()

	// Create a new buffer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("files", "index.html")
	if err != nil {
		fmt.Println("Create From File: ", err)
		return nil, err
	}

	// Copy the file content to the part
	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println("Copy: ", err)
		return nil, err
	}
	writer.Close()

	// Create a new HTTP request
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		fmt.Println("New Request: ", err)
		return nil, err
	}
	request.Header.Add("Content-Type", writer.FormDataContentType())

	// Send the request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Client.Do(): ", err)
		return nil, err
	}
	defer response.Body.Close()

	// Create the output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Println("os.Create: ", err)
		return nil, err
	}
	defer outputFile.Close()

	// Copy the response body to the output file
	_, err = io.Copy(outputFile, response.Body)
	if err != nil {
		fmt.Println("Copy 2: ", err)
		return nil, err
	}

	return &outputPath, nil

}
