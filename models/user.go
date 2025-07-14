package models

type User struct {
	id     int64
	age    int
	height float64
}

func (u *User) SetAge(age int) {
	u.age = age
}

func (u *User) SetHeight(height float64) {
	u.height = height
}

func (u *User) GetId() int64 {
	return u.id
}

func (u *User) GetAge() int {
	return u.age
}

func (u *User) GetHeight() float64 {
	return u.height
}

func NewUser(id int64, age int, height float64) User {
	return User{id: id, age: age, height: height}
}
