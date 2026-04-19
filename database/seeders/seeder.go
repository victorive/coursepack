package seeders

import "gorm.io/gorm"

type Seeder interface {
	Name() string
	Run(db *gorm.DB) error
}

var orderedRegistry = []Seeder{
	universitiesSeeder{},
	departmentsSeeder{},
	permissionsSeeder{},
	levelsSeeder{},
	rolesSeeder{},
	rolePermissionsSeeder{},
	superAdminSeeder{},
}

func All() []Seeder {
	return orderedRegistry
}

func Find(name string) Seeder {
	for _, s := range orderedRegistry {
		if s.Name() == name {
			return s
		}
	}
	return nil
}
