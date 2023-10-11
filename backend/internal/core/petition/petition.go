package petition

import (
	"context"
	"html/template"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"

	"mado/pkg/validator"
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

// Todo save to db ?
func (s Service) GeneratePetitionPDF(data *PetitionData) (*PetitionData, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	currentTime := time.Now()
	// Format the date as "21.09.2023"
	data.CreationDate = currentTime.Format("02.01.2006")

	var err error
	// data.SheetNumber, err = s.petitionRepository.GetNextID(ctx)
	// if err != nil {
	// 	s.logger.Error("Failed to get id for next doc: ", zap.Error(err))
	// 	return nil, err
	// }

	// Create a new template It will be used to format the HTML content of the petition.
	tmpl, err := template.New("petition").Parse(TemplateHTML)
	if err != nil {
		s.logger.Error("Failed to parse template: ", zap.Error(err))
		return nil, err
	}

	// Create a buffer to hold the generated HTML content
	var htmlContentBuffer strings.Builder
	if err := tmpl.Execute(&htmlContentBuffer, data); err != nil {
		s.logger.Error("Failed to execute data to buffer: ", zap.Error(err))
		return nil, err
	}

	// he HTML content in the htmlContentBuffer is parsed using the goquery library. This library allows you to manipulate and traverse HTML documents in a structured way
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContentBuffer.String()))
	if err != nil {
		s.logger.Error("Failed to pass htmlContent to goquery: ", zap.Error(err))
		return nil, err
	}

	// Save the modified HTML to a temporary file
	tempHTMLFileName := "temp.html"
	tempHTMLFile, err := os.Create(tempHTMLFileName)
	if err != nil {
		s.logger.Error("Failed to create template file: ", zap.Error(err))
		return nil, err
	}
	defer tempHTMLFile.Close()

	// Use goquery to write the modified HTML to the file
	htmlContent, err := doc.Html()
	if err != nil {
		s.logger.Error("Failed to get html content: ", zap.Error(err))
		return nil, err
	}
	_, err = tempHTMLFile.WriteString(htmlContent)
	if err != nil {
		s.logger.Error("Failed to write the modified HTML to the file: ", zap.Error(err))
		return nil, err
	}

	// Generate PDF using wkhtmltopdf
	//todo change namingOfFile
	// pdfFileName := "output.pdf"
	cmd := exec.Command("wkhtmltopdf", tempHTMLFileName, data.FileName)
	err = cmd.Run()
	if err != nil {
		s.logger.Error("Failed to convert html to pdf: ", zap.Error(err))
		return nil, err
	}
	s.logger.Debug("PDF generated: %s\n", zap.String("filename", data.FileName))

	// Prepare the PDF file data and metadata
	data.PdfData, err = os.ReadFile(data.FileName)
	if err != nil {
		s.logger.Error("Failed to read PDF file: ", zap.Error(err))
		return nil, err
	}
	// return nil, nil
	return s.petitionRepository.Save(ctx, data)
}

func (s Service) GetPetitionPdfByID(id string) (*PetitionData, error) {
	docId, err := validator.IdValidator(id)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.petitionRepository.GetPetitionPdfByID(ctx, docId)
}
