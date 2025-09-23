package model

import "github.com/google/uuid"

// RADEK: není to uplně nutnost, ale já mám rádši, když má model konstruktor, protože, potom se ti lépe hledá a upravuje definice takového modelu
// to platí pro všechny modely
type Answer struct {
	ID         uuid.UUID
	QuestionID uuid.UUID
	Text       string
	IsCorrect  bool
}
