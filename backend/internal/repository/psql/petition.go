package psql

import (
	"context"

	"go.uber.org/zap"

	"mado/internal/core/petition"
	"mado/pkg/database/postgres"
	"mado/pkg/errs"
)

type PetitionRepository struct {
	db     *postgres.Postgres
	logger *zap.Logger
}

// PetitionRepository creates a new UserRepository.
func NewPetitionRepository(db *postgres.Postgres, logger *zap.Logger) PetitionRepository {
	return PetitionRepository{
		db:     db,
		logger: logger,
	}
}

// todo create table for this and think about numering documents IDEA: just before create new file get id for the next row in postgre
func (p PetitionRepository) Save(ctx context.Context, dto *petition.PetitionData) (*petition.PetitionData, error) {
	// Insert the PDF data into the database
	_, err := p.db.Pool.Exec(ctx, `
	  INSERT INTO petition_pdf_files (file_name, sheet_number, creation_date, location, responsible_person, owner_name, owner_address, pdf_data)
	  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		dto.FileName, dto.SheetNumber, dto.CreationDate, dto.Location, dto.ResponsiblePerson, dto.OwnerName, dto.OwnerAddress, dto.PdfData)
	if err != nil {
		p.logger.Error("Failed to insert PDF data into the database: ", zap.Error(err))
		return nil, err
	}
	return dto, nil
}

// We use the nextval function with the sequence name petitions_id_seq to retrieve the next available ID.
func (p PetitionRepository) GetNextID(ctx context.Context) (*int, error) {
	query := "SELECT nextval('petition_pdf_files_id_seq')" //name of table petition_pdf_files + _id_seq

	var id int
	if err := p.db.Pool.QueryRow(ctx, query).Scan(&id); err != nil {
		return nil, err
	}

	return &id, nil
}

// CREATE TABLE petition_pdf_files (
//     id serial PRIMARY KEY,
//     file_name varchar(255) NOT NULL,
//     creation_date timestamp,
//     location varchar(255),
//     responsible_person varchar(255),
//     owner_name varchar(255),
//     owner_address varchar(255),
//     pdf_data bytea NOT NULL
// );

func (p PetitionRepository) GetPetitionPdfByID(ctx context.Context, pdfID *int) (*petition.PetitionData, error) {
	// Retrieve the PDF data from the database based on the ID
	var pdfData []byte
	var filename string
	err := p.db.Pool.QueryRow(ctx, "SELECT pdf_data, file_name FROM petition_pdf_files WHERE id = $1", pdfID).Scan(&pdfData, &filename)
	if err != nil {
		p.logger.Error(errs.ErrPdfFileNotFound.Error(), zap.Error(err))
		return nil, errs.ErrPdfFileNotFound
	}

	return &petition.PetitionData{
		PdfData:  pdfData,
		FileName: filename,
	}, nil
}
