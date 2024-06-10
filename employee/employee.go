package employee

type Employee struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Work   string `json:"work"`
	Age    int    `json:"age"`
	Salary int    `json:"salary"`
}

func NewEmployee(id int, name, work string, age, salary int) *Employee {
	return &Employee{
		Id:     id,
		Name:   name,
		Work:   work,
		Age:    age,
		Salary: salary,
	}
}
