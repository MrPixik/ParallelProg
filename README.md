# Параллельное программирование

A URL shortening service written in Go that allows users to submit URLs and receive shortened versions in return.

## Project Structure

```bash
ParallelProg/
├── lab1/                       #1 laboratory work
│   └── main.go          
├── lab2/                       #2 laboratory work
│   ├── proc/               
│   │   ├── process.go      
│   │   └── process.exe     
│   └── main.go                      
├── go.mod                          # Go module file
└── .gitignore                      # Git ignore file
```
## Установка

Чтобы установить этот проект локально, выполните следующие шаги:

1. Клонируйте репозиторий:

``` bash
git clone https://github.com/MrPixik/ParallelProg.git
cd ParallelProg
```
2. Установите зависимости:

```bash
go mod tidy
```

## 1 Лабораторная работа
### Формулировка задания:
Измерить производительность вашего персонального компьютера на операциях следующих типов:
- сложение,
- вычитание,
- умножение,
- деление,
- возведение в степень,
  а также на вычислениях функций математической библиотеки – exp, log, sin.
### Запуск программы:
Синтаксис запуска основной программы:
```bash
go run lab1/main.go
```
### Ожидаемый вывод:
``` bash
Time: <время_выполнения> '+' perf.: <кол-во GFlops>
Time: <время_выполнения> '-' perf.: <кол-во GFlops>
...
```
## 2 Лабораторная работа
### Формулировка задания:
Написать программу, решающую задачу численного интегрирования
с помощью тяжелых и легких процессов одновременно.
### Сборка:
В корневой директории проекта выполнить команду:
``` bash
go build -o lab2/proc/process.exe lab2/proc/process.go
```
### Запуск программы:
Синтаксис запуска основной программы:
```bash
go run main.go <кол-во_процессов> <кол-во_потоков_в_каждом_процессе>
```
### Ожидаемый вывод:
``` bash
Time: <время_выполнения> Integral: <значение_интеграла>
```
