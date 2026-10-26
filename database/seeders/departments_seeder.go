package seeders

import (
	"log/slog"

	"github.com/google/uuid"
	"github.com/victorive/coursepack/internal/domain"
	"gorm.io/gorm"
)

type departmentsSeeder struct{}

func (s departmentsSeeder) Name() string {
	return "departments"
}

func (s departmentsSeeder) Run(db *gorm.DB) error {
	logger := slog.Default()

	var unnUniv, uniabujaUniv domain.University
	if err := db.Where("acronym = ?", "UNN").First(&unnUniv).Error; err != nil {
		logger.Error("failed to find UNN university", "error", err)
		return err
	}
	if err := db.Where("acronym = ?", "UNIABUJA").First(&uniabujaUniv).Error; err != nil {
		logger.Error("failed to find UNIABUJA university", "error", err)
		return err
	}

	unnDepts := []domain.Department{
		{UniversityID: unnUniv.ID, Name: "Agriculture", Code: "AG"},
		{UniversityID: unnUniv.ID, Name: "Agronomy", Code: "AGR"},
		{UniversityID: unnUniv.ID, Name: "Animal Science", Code: "ANS"},
		{UniversityID: unnUniv.ID, Name: "Agricultural Economics", Code: "AEC"},
		{UniversityID: unnUniv.ID, Name: "Agricultural Extension", Code: "AEX"},
		{UniversityID: unnUniv.ID, Name: "Agricultural & Bioresources Engineering", Code: "ABE"},
		{UniversityID: unnUniv.ID, Name: "Crop Science", Code: "CSC"},
		{UniversityID: unnUniv.ID, Name: "Food Science & Technology", Code: "FST"},
		{UniversityID: unnUniv.ID, Name: "Soil Science", Code: "SSC"},
		{UniversityID: unnUniv.ID, Name: "Archaeology & Tourism", Code: "ART"},
		{UniversityID: unnUniv.ID, Name: "English & Literary Studies", Code: "ELS"},
		{UniversityID: unnUniv.ID, Name: "Fine & Applied Arts", Code: "FAA"},
		{UniversityID: unnUniv.ID, Name: "French", Code: "FRE"},
		{UniversityID: unnUniv.ID, Name: "History & International Studies", Code: "HIS"},
		{UniversityID: unnUniv.ID, Name: "Linguistics & Nigerian Languages", Code: "LNL"},
		{UniversityID: unnUniv.ID, Name: "Mass Communication", Code: "MAC"},
		{UniversityID: unnUniv.ID, Name: "Music", Code: "MUS"},
		{UniversityID: unnUniv.ID, Name: "Philosophy", Code: "PHV"},
		{UniversityID: unnUniv.ID, Name: "Religious and Cultural Studies", Code: "RCS"},
		{UniversityID: unnUniv.ID, Name: "Theatre & Festival Studies", Code: "TFS"},
		{UniversityID: unnUniv.ID, Name: "Theatre Arts", Code: "THA"},
		{UniversityID: unnUniv.ID, Name: "Biochemistry", Code: "BCH"},
		{UniversityID: unnUniv.ID, Name: "Biology", Code: "BIO"},
		{UniversityID: unnUniv.ID, Name: "Botany", Code: "BOT"},
		{UniversityID: unnUniv.ID, Name: "Microbiology", Code: "MCB"},
		{UniversityID: unnUniv.ID, Name: "Plant Science & Biotechnology", Code: "PSB"},
		{UniversityID: unnUniv.ID, Name: "Zoology & Environmental Biology", Code: "ZEB"},
		{UniversityID: unnUniv.ID, Name: "Zoology", Code: "ZOO"},
		{UniversityID: unnUniv.ID, Name: "Accountancy", Code: "ACC"},
		{UniversityID: unnUniv.ID, Name: "Banking & Finance", Code: "BNF"},
		{UniversityID: unnUniv.ID, Name: "Business Administration", Code: "BUS"},
		{UniversityID: unnUniv.ID, Name: "Management", Code: "MAN"},
		{UniversityID: unnUniv.ID, Name: "Marketing", Code: "MKT"},
		{UniversityID: unnUniv.ID, Name: "Adult Education", Code: "ADU"},
		{UniversityID: unnUniv.ID, Name: "Arts Education", Code: "AED"},
		{UniversityID: unnUniv.ID, Name: "Commercial Education", Code: "CED"},
		{UniversityID: unnUniv.ID, Name: "Education", Code: "EDU"},
		{UniversityID: unnUniv.ID, Name: "Educational Foundations", Code: "EDF"},
		{UniversityID: unnUniv.ID, Name: "Health & Physical Education", Code: "HPE"},
		{UniversityID: unnUniv.ID, Name: "Home Science, Nutrition & Dietetics", Code: "HND"},
		{UniversityID: unnUniv.ID, Name: "Library & Information Science", Code: "LSE"},
		{UniversityID: unnUniv.ID, Name: "Science Education", Code: "SED"},
		{UniversityID: unnUniv.ID, Name: "Social Science Education", Code: "SSE"},
		{UniversityID: unnUniv.ID, Name: "Vocational Teacher Education", Code: "VTE"},
		{UniversityID: unnUniv.ID, Name: "Civil Engineering", Code: "CE"},
		{UniversityID: unnUniv.ID, Name: "Electrical Engineering", Code: "ELG"},
		{UniversityID: unnUniv.ID, Name: "Electronic Engineering", Code: "EEG"},
		{UniversityID: unnUniv.ID, Name: "Materials & Metallurgical Engineering", Code: "MME"},
		{UniversityID: unnUniv.ID, Name: "Mechanical Engineering", Code: "MEG"},
		{UniversityID: unnUniv.ID, Name: "Architecture", Code: "ARC"},
		{UniversityID: unnUniv.ID, Name: "Estate Management", Code: "ESM"},
		{UniversityID: unnUniv.ID, Name: "Estate Management", Code: "ESMA"},
		{UniversityID: unnUniv.ID, Name: "Surveying & Geodesy", Code: "SVY"},
		{UniversityID: unnUniv.ID, Name: "Urban & Regional Planning", Code: "URP"},
		{UniversityID: unnUniv.ID, Name: "Law", Code: "LAW"},
		{UniversityID: unnUniv.ID, Name: "International Law & Jurisprudence", Code: "ILJ"},
		{UniversityID: unnUniv.ID, Name: "Political Science", Code: "POL"},
		{UniversityID: unnUniv.ID, Name: "Property Law", Code: "PRL"},
		{UniversityID: unnUniv.ID, Name: "Public & Private Law", Code: "PPL"},
		{UniversityID: unnUniv.ID, Name: "Anaesthesia", Code: "ANA"},
		{UniversityID: unnUniv.ID, Name: "Anatomy", Code: "ANT"},
		{UniversityID: unnUniv.ID, Name: "Chemical Pathology", Code: "CPT"},
		{UniversityID: unnUniv.ID, Name: "Child Dental Health", Code: "CDH"},
		{UniversityID: unnUniv.ID, Name: "Community Medicine", Code: "CMA"},
		{UniversityID: unnUniv.ID, Name: "Dermatology", Code: "DER"},
		{UniversityID: unnUniv.ID, Name: "Haematology & Immunology", Code: "HIM"},
		{UniversityID: unnUniv.ID, Name: "Medical Biochemistry", Code: "MBH"},
		{UniversityID: unnUniv.ID, Name: "Medical Laboratory Technology", Code: "MLT"},
		{UniversityID: unnUniv.ID, Name: "Medical Microbiology", Code: "MCM"},
		{UniversityID: unnUniv.ID, Name: "Medical Rehabilitation", Code: "MRH"},
		{UniversityID: unnUniv.ID, Name: "Morbid Anatomy", Code: "MPA"},
		{UniversityID: unnUniv.ID, Name: "Nursing Sciences", Code: "NUR"},
		{UniversityID: unnUniv.ID, Name: "Obstetrics & Gynecology", Code: "OGY"},
		{UniversityID: unnUniv.ID, Name: "Ophthalmology", Code: "OPT"},
		{UniversityID: unnUniv.ID, Name: "Oral Maxillofacial Surgery", Code: "OMS"},
		{UniversityID: unnUniv.ID, Name: "Otolaryngology", Code: "OTL"},
		{UniversityID: unnUniv.ID, Name: "Paediatric Surgery", Code: "PSG"},
		{UniversityID: unnUniv.ID, Name: "Paediatrics", Code: "PED"},
		{UniversityID: unnUniv.ID, Name: "Pharmacology & Therapeutics", Code: "PHR"},
		{UniversityID: unnUniv.ID, Name: "Physiological Medicine", Code: "PMS"},
		{UniversityID: unnUniv.ID, Name: "Preventive Dentistry", Code: "PDE"},
		{UniversityID: unnUniv.ID, Name: "Radiation Medicine", Code: "RAD"},
		{UniversityID: unnUniv.ID, Name: "Restorative Dentistry", Code: "RDE"},
		{UniversityID: unnUniv.ID, Name: "Surgery", Code: "SUR"},
		{UniversityID: unnUniv.ID, Name: "Clinical Pharmacy", Code: "CLP"},
		{UniversityID: unnUniv.ID, Name: "Pharmaceutical Chemistry", Code: "PCH"},
		{UniversityID: unnUniv.ID, Name: "Pharmaceutical Technology", Code: "PHE"},
		{UniversityID: unnUniv.ID, Name: "Pharmaceutical and Medicinal Chemistry", Code: "PMC"},
		{UniversityID: unnUniv.ID, Name: "Pharmaceutics", Code: "PHS"},
		{UniversityID: unnUniv.ID, Name: "Pharmacognosy", Code: "PHG"},
		{UniversityID: unnUniv.ID, Name: "Pharmacognosy and Environmental Medicines", Code: "PEM"},
		{UniversityID: unnUniv.ID, Name: "Pharmacology and Toxicology", Code: "PHT"},
		{UniversityID: unnUniv.ID, Name: "Computer Science", Code: "COS"},
		{UniversityID: unnUniv.ID, Name: "Computer Science", Code: "CS"},
		{UniversityID: unnUniv.ID, Name: "Computer Applications", Code: "CSP"},
		{UniversityID: unnUniv.ID, Name: "Geology", Code: "GEOL"},
		{UniversityID: unnUniv.ID, Name: "Geology", Code: "GLG"},
		{UniversityID: unnUniv.ID, Name: "Mathematics", Code: "MTH"},
		{UniversityID: unnUniv.ID, Name: "Physics & Astronomy", Code: "PHY"},
		{UniversityID: unnUniv.ID, Name: "Pure and Industrial Chemistry", Code: "CHM"},
		{UniversityID: unnUniv.ID, Name: "Statistics", Code: "STA"},
		{UniversityID: unnUniv.ID, Name: "Statistics", Code: "STAT"},
		{UniversityID: unnUniv.ID, Name: "Economics", Code: "ECO"},
		{UniversityID: unnUniv.ID, Name: "Geography", Code: "GMK"},
		{UniversityID: unnUniv.ID, Name: "Psychology", Code: "PSY"},
		{UniversityID: unnUniv.ID, Name: "Public Administration", Code: "PAD"},
		{UniversityID: unnUniv.ID, Name: "Social Work", Code: "SWK"},
		{UniversityID: unnUniv.ID, Name: "Sociology & Anthropology", Code: "SOC"},
		{UniversityID: unnUniv.ID, Name: "Animal Health & Production", Code: "AHP"},
		{UniversityID: unnUniv.ID, Name: "Veterinary Anatomy", Code: "VAN"},
		{UniversityID: unnUniv.ID, Name: "Veterinary Medicine", Code: "VET"},
		{UniversityID: unnUniv.ID, Name: "Veterinary Obstetrics and Reproductive Diseases", Code: "VOR"},
		{UniversityID: unnUniv.ID, Name: "Veterinary Parasitology & Entomology", Code: "VPE"},
		{UniversityID: unnUniv.ID, Name: "Veterinary Pathology and Microbiology", Code: "VPM"},
		{UniversityID: unnUniv.ID, Name: "Veterinary Physiology and Pharmacology", Code: "VPP"},
		{UniversityID: unnUniv.ID, Name: "Veterinary Public Health & Preventive Medicine", Code: "VPH"},
		{UniversityID: unnUniv.ID, Name: "Veterinary Surgery", Code: "VSG"},
		{UniversityID: unnUniv.ID, Name: "General Studies", Code: "GST"},
		{UniversityID: unnUniv.ID, Name: "General Studies Programme", Code: "GSP"},
		{UniversityID: unnUniv.ID, Name: "Science", Code: "SCI"},
	}

	uniabujaDepts := []domain.Department{
		{UniversityID: uniabujaUniv.ID, Name: "Agricultural Economics", Code: "AGEC"},
		{UniversityID: uniabujaUniv.ID, Name: "Agricultural Economics and Extension", Code: "AGEE"},
		{UniversityID: uniabujaUniv.ID, Name: "Agricultural Extension and Rural Sociology", Code: "AGERS"},
		{UniversityID: uniabujaUniv.ID, Name: "Agronomy", Code: "AGRO"},
		{UniversityID: uniabujaUniv.ID, Name: "Animal Science", Code: "ANSC"},
		{UniversityID: uniabujaUniv.ID, Name: "Crop Protection", Code: "CRPR"},
		{UniversityID: uniabujaUniv.ID, Name: "Crop Science", Code: "CRSC"},
		{UniversityID: uniabujaUniv.ID, Name: "Diary Science", Code: "DIRY"},
		{UniversityID: uniabujaUniv.ID, Name: "Fisheries Aquaculture and Wildlife", Code: "FAQW"},
		{UniversityID: uniabujaUniv.ID, Name: "Food Science and Technology", Code: "FSAT"},
		{UniversityID: uniabujaUniv.ID, Name: "Forestry and Bioresources", Code: "FOR"},
		{UniversityID: uniabujaUniv.ID, Name: "Horticulture and Landscaping", Code: "HORT"},
		{UniversityID: uniabujaUniv.ID, Name: "Soil Science", Code: "SOIL"},
		{UniversityID: uniabujaUniv.ID, Name: "Arabic", Code: "ARAB"},
		{UniversityID: uniabujaUniv.ID, Name: "Christian Religious studies", Code: "CRST"},
		{UniversityID: uniabujaUniv.ID, Name: "English", Code: "ENGL"},
		{UniversityID: uniabujaUniv.ID, Name: "History and Diplomatic Studies", Code: "HIST"},
		{UniversityID: uniabujaUniv.ID, Name: "Islamic Studies", Code: "ISSL"},
		{UniversityID: uniabujaUniv.ID, Name: "Linguistics and African Language", Code: "LING"},
		{UniversityID: uniabujaUniv.ID, Name: "Philosophy", Code: "PHIL"},
		{UniversityID: uniabujaUniv.ID, Name: "Theatre Arts", Code: "THEA"},
		{UniversityID: uniabujaUniv.ID, Name: "Accounting", Code: "ACCN"},
		{UniversityID: uniabujaUniv.ID, Name: "Banking and Finance", Code: "BNKF"},
		{UniversityID: uniabujaUniv.ID, Name: "Business Administration", Code: "BUSN"},
		{UniversityID: uniabujaUniv.ID, Name: "Entrepreneurship Studies", Code: "ENTR"},
		{UniversityID: uniabujaUniv.ID, Name: "Public Administration", Code: "PUAD"},
		{UniversityID: uniabujaUniv.ID, Name: "Tourism and Hospitality Management", Code: "THOM"},
		{UniversityID: uniabujaUniv.ID, Name: "Computer Science", Code: "COSC"},
		{UniversityID: uniabujaUniv.ID, Name: "Mathematics", Code: "MATH"},
		{UniversityID: uniabujaUniv.ID, Name: "Microbiology", Code: "MICR"},
		{UniversityID: uniabujaUniv.ID, Name: "Physics", Code: "PHYS"},
		{UniversityID: uniabujaUniv.ID, Name: "Statistics", Code: "STAT"},
		{UniversityID: uniabujaUniv.ID, Name: "Zoology", Code: "ZOOL"},
		{UniversityID: uniabujaUniv.ID, Name: "Biochemistry", Code: "BIOC"},
		{UniversityID: uniabujaUniv.ID, Name: "Biology", Code: "BIOL"},
		{UniversityID: uniabujaUniv.ID, Name: "Botany", Code: "BOTA"},
		{UniversityID: uniabujaUniv.ID, Name: "Chemistry", Code: "CHEM"},
		{UniversityID: uniabujaUniv.ID, Name: "Geology and Mining", Code: "GEOL"},
		{UniversityID: uniabujaUniv.ID, Name: "Aeronautic and Astronautical Engineering", Code: "AAEG"},
		{UniversityID: uniabujaUniv.ID, Name: "Agricultural Engineering", Code: "AGEG"},
		{UniversityID: uniabujaUniv.ID, Name: "Chemical Engineering", Code: "CHEG"},
		{UniversityID: uniabujaUniv.ID, Name: "Civil Engineering", Code: "CIEG"},
		{UniversityID: uniabujaUniv.ID, Name: "Electrical/Electronics Engineering", Code: "ELEG"},
		{UniversityID: uniabujaUniv.ID, Name: "Mechanical Engineering", Code: "MEEG"},
		{UniversityID: uniabujaUniv.ID, Name: "Nuclear Engineering", Code: "NUEG"},
		{UniversityID: uniabujaUniv.ID, Name: "Railway Engineering", Code: "RAEG"},
	}

	allDepts := append(unnDepts, uniabujaDepts...)
	for i := range allDepts {
		var existing domain.Department
		err := db.Where(domain.Department{UniversityID: allDepts[i].UniversityID, Code: allDepts[i].Code}).First(&existing).Error
		if err == nil {
			continue
		}

		allDepts[i].UUID = uuid.NewString()
		if err := db.Create(&allDepts[i]).Error; err != nil {
			logger.Error("failed to seed department", "error", err, "name", allDepts[i].Name, "code", allDepts[i].Code)
			return err
		}
	}

	logger.Info("departments seeded successfully", "count", len(allDepts))
	return nil
}
