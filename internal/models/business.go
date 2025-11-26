package models

import (
	"meu_job/utils/validator"
	"slices"
	"unicode"
)

type Business struct {
	ID    int64
	Name  string
	CNPJ  string
	Email string
	Phone string
	User  *User
	BaseModel
}

type BusinessDTO struct {
	ID    *int64   `json:"business_id"`
	Name  *string  `json:"name"`
	CNPJ  *string  `json:"cnpj"`
	Email *string  `json:"email"`
	Phone *string  `json:"phone"`
	User  *UserDTO `json:"user"`
}

func (b Business) ToDTO() *BusinessDTO {
	return &BusinessDTO{
		ID:    &b.ID,
		Name:  &b.Name,
		CNPJ:  &b.CNPJ,
		Email: &b.Email,
		Phone: &b.Phone,
		User:  b.User.ToDTO(),
	}
}

func (b BusinessDTO) ToModel() *Business {
	var model = &Business{}
	if b.ID != nil {
		model.ID = *b.ID
	}
	if b.Name != nil {
		model.Name = *b.Name
	}

	if b.CNPJ != nil {
		model.CNPJ = *b.CNPJ
	}

	if b.Phone != nil {
		model.Phone = *b.Phone
	}

	if b.Email != nil {
		model.Email = *b.Email
	}

	if b.User != nil {
		model.User = b.User.ToModel()
	}
	return model
}

func (m *Business) ValidateBusiness(v *validator.Validator) {
	v.Check(m.Name != "", "name", "must be provided")
	v.Check(len(m.Name) <= 500, "name", "must not be more than 500 bytes long")
	v.Check(m.Phone != "", "phone", "must be provided")
	ValidateEmail(v, m.Email)
	ValidateCNPJ(v, m.CNPJ)
}

func ValidateCNPJ(v *validator.Validator, cnpj string) {
	v.Check(cnpj != "", "cnpj", "must be provided")
	v.Check(isValidCNPJ(cnpj), "cnpj", "must be a valid cnpj")
}

func normalizeCNPJ(cnpj string) string {
	var out []rune
	for _, r := range cnpj {
		if unicode.IsDigit(r) {
			out = append(out, r)
		}
	}

	return string(out)
}

func isValidCNPJ(cnpj string) bool {
	cnpj = normalizeCNPJ(cnpj)

	if len(cnpj) != 14 {
		return false
	}

	invalids := []string{
		"00000000000000", "11111111111111", "22222222222222",
		"33333333333333", "44444444444444", "55555555555555",
		"66666666666666", "77777777777777", "88888888888888",
		"99999999999999",
	}

	if slices.Contains(invalids, cnpj) {
		return false
	}

	var nums []int
	for _, c := range cnpj {
		nums = append(nums, int(c-'0'))
	}

	w1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	w2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}

	sum := 0
	for i := range 12 {
		sum += nums[i] * w1[i]
	}
	r1 := sum % 11
	d1 := 0
	if r1 >= 2 {
		d1 = 11 - r1
	}

	if nums[12] != d1 {
		return false
	}

	sum = 0
	for i := range 13 {
		sum += nums[i] * w2[i]
	}
	r2 := sum % 11
	d2 := 0
	if r2 >= 2 {
		d2 = 11 - r2
	}

	return nums[13] == d2
}
