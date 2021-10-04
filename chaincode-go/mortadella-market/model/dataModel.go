package model

//MortadellaIndexString the partial key used for build the composite key for a MortadellaType
const MortadellaIndexString string = "MortadellaType"

//CompanyIndexString the partial key used for build composite key for a company
const CompanyIndexString string = "Company"

// The StateInterface defines the basic method for a state
type StateInterface interface {
	GetIndex() string
	GetKey() []string
}

//State is the basic data structure used in this chaincode
type State struct {
	ID string `json:"id"`
}

//Company this class represent a Comapny
type Company struct {
	State
	Name            string `json:"name"`
	Identity        string `json:"identity"`
	MortadellaUnits int64  `json:"mortadellaUnits"`
}

//ComapanyDTO a DTO for the Company Class
type ComapanyDTO struct {
	ID              string `json:"ID"`
	Name            string `json:"name"`
	MortadellaUnits int64  `json:"mortadellaUnits"`
}

//FromCompany create a DTO from a company
func (dto *ComapanyDTO) FromCompany(company Company) {
	dto.ID = company.ID
	dto.Name = company.Name
	dto.MortadellaUnits = company.MortadellaUnits
}

//GetIndex git index for the state
func (c Company) GetIndex() string {
	return CompanyIndexString
}

//GetCompanyKey get a key for a Company
func GetCompanyKey(CompanyID string) []string {
	return []string{CompanyID}
}

//GetKey get a key for a state
func (c Company) GetKey() []string {
	return GetCompanyKey(c.ID)
}
