/*
Напишите программу, которая считывает ввод с stdin,
создаёт структуру student и записывает указатель на структуру в
хранилище map[studentName] *Student.

type Student struct {
	name string
	age int
	grade int
}

Программа должна получать строки в бесконечном цикле,
создать структуру Student через функцию newStudent,
далее сохранить указатель на эту структуру в map,
а после получения EOF (ctrl + d) вывести на экран имена всех
студентов из хранилища. Также необходимо реализовать методы put, get.

Input  ---- go run main.go
Строки ---- Вася 24 1
			Cемен 32 2
EOF
Output ---- Студенты из хранилища:
			Вася 24 1
			Семен 32 2

Зачёт:
	- при получении одной строки (например, «имяСтудента 24 1»)
	программа создаёт студента и сохраняет его, далее ожидает
	следующую строку или сигнал EOF (Сtrl + Z);
	-при получении сигнала EOF программа должна вывести имена
	всех студентов из map.
*/

package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Storage interface {
	Put(string) bool
	Get(string) string
	Students()
}

type Student struct {
	name  string
	age   int
	grade int
}

/*
//---------------------- начало заглушки -------------------

type StubStorage struct{}

func (fs *StubStorage) Put(string) bool {
	return true
}

func (fs *StubStorage) Get(string) string {
	return "default"
}

func (fs *StubStorage) Students() []string {
	return []string{"DefoltName", "DefaultNAme", "DefaultName"}
}
//---------------------- конец заглушки -------------------
*/

type App struct {
	repository Storage
}

func (a *App) Run() {
	for {
		a.printStudents()
		if student, ok := a.inputNextStudent(); ok {
			a.storeStudent(student)
		} else {
			break
		}
	}
}

func (a *App) printStudents() {
	fmt.Println("Список введенных студентов с оценками:")
	a.repository.Students()
}

func (a *App) inputNextStudent() (string, bool) {
	for {
		fmt.Print("Введите данные студента или Ctrl + Z для завершения: ")
		var inputName string
		var inputAge, inputGrade int
		_, err := fmt.Scanln(&inputName, &inputAge, &inputGrade)

		inputInfo := inputName + " " + strconv.Itoa(inputAge) + " " + strconv.Itoa(inputGrade)

		if err == io.EOF {
			var getStudent, name string
			fmt.Println("Если хотите найти студента в списке, введите 'find', если нет, еще раз введите Ctrl + Z")
			fmt.Scan(&getStudent)
			if getStudent == "find" {
				fmt.Print("Введите имя студента: ")
				fmt.Scan(&name)
				fmt.Println(a.repository.Get(name))
				return "end", false
			} else {
				a.printStudents()
				return "end", false
			}
		} else {
			return inputInfo, true
		}
	}
}

func (a *App) storeStudent(student string) {
	msg := "Студент уже присутствует в коллекции\n"
	if ok := a.repository.Put(student); ok {
		msg = "Студент успешно добавлен\n"
	}
	fmt.Printf(msg, student)
}

//-----------------реализация репозитория----------------

type MemStorage struct {
	students map[string]*Student
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		students: make(map[string]*Student),
	}
}

func (ms *MemStorage) Put(inputInfo string) bool {
	inputInfoSplited := strings.Split(inputInfo, " ")
	var StudentInfo Student
	StudentInfo.name = inputInfoSplited[0]
	StudentInfo.age, _ = strconv.Atoi(inputInfoSplited[1])
	StudentInfo.grade, _ = strconv.Atoi(inputInfoSplited[2])

	if ms.contains(StudentInfo.name) {
		return false
	}
	ms.students[StudentInfo.name] = &StudentInfo
	return true
}

func (ms *MemStorage) Get(name string) string {
	getStudent, ok := ms.students[name]
	if ok == true {
		return getStudent.name + " " + strconv.Itoa(getStudent.age) + " " + strconv.Itoa(getStudent.grade)
	} else {
		return "Такого студента нет в коллекции"
	}
}

func (ms *MemStorage) Students() {
	for key, v := range ms.students {
		fmt.Printf("%s: %v %d %d\n", key, (v.name), (v.age), (v.grade))
	}
}

func (ms *MemStorage) contains(stName string) bool {
	for key, _ := range ms.students {
		if ms.students[key].name == stName {
			return true
		}
	}
	return false
}

func main() {
	//repository := &StubStorage{}
	repository := NewMemStorage()
	app := &App{repository}
	app.Run()
}
