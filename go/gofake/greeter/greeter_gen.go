package greeter

import "gofake/greeter/animal"

type Persons []*Person

func (ss Persons) IDs() []int64 {
	res := make([]int64, len(ss))
	for i, s := range ss {
		res[i] = s.ID
	}
	return res
}

func (ss Persons) Profiles() []Profile {
	res := make([]Profile, len(ss))
	for i, s := range ss {
		res[i] = s.Profile
	}
	return res
}

func (ss Persons) Addresss() []*Address {
	res := make([]*Address, len(ss))
	for i, s := range ss {
		res[i] = s.Address
	}
	return res
}

func (ss Persons) Hobbiess() [][]Hobby {
	res := make([][]Hobby, len(ss))
	for i, s := range ss {
		res[i] = s.Hobbies
	}
	return res
}

func (ss Persons) Petss() [][]*Pet {
	res := make([][]*Pet, len(ss))
	for i, s := range ss {
		res[i] = s.Pets
	}
	return res
}

func (ss Persons) Jobs() []*struct {
	Name string
} {
	res := make([]*struct {
		Name string
	}, len(ss))
	for i, s := range ss {
		res[i] = s.Job
	}
	return res
}

func (ss Persons) Foods() []struct {
	Name string
} {
	res := make([]struct {
		Name string
	}, len(ss))
	for i, s := range ss {
		res[i] = s.Food
	}
	return res
}

func (ss Persons) valAnimals() []animal.Animal {
	res := make([]animal.Animal, len(ss))
	for i, s := range ss {
		res[i] = s.valAnimal
	}
	return res
}

func (ss Persons) ptrAnimals() []*animal.Animal {
	res := make([]*animal.Animal, len(ss))
	for i, s := range ss {
		res[i] = s.ptrAnimal
	}
	return res
}
