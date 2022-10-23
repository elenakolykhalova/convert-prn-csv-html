package model

type User struct {
	Name     string
	Address  string
	Postcode string
	Mobile   string
	Limit    string
	Birthday string
}

func GetAllUsers() (users []User, err error) {
	users = []User{
		{"Oliver, El","Via Archimede, 103-91","2343aa","000 1119381","6000000","01/01/1999"},
		{"Harry","Leonardo da Vinci 1","4532 AA","010 1118986","343434","31/12/1965"},
	}
	return
}
