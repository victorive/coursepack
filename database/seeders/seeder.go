package seeders

import "gorm.io/gorm"

type Seeder interface {
	Name() string
	Run(db *gorm.DB) error
}

var registry []Seeder

func register(s Seeder) {
	registry = append(registry, s)
}

func All() []Seeder {
	return registry
}

func Find(name string) Seeder {
	for _, s := range registry {
		if s.Name() == name {
			return s
		}
	}
	return nil
}
