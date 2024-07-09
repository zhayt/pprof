# Profiling and Optimization in Go
В этой репе постараюсь достаточно просто обяснить что такое профилирование и как его пременять на практике стобы оптемелировать ваши программы.
## What is profiling and optimization?
Если ваша программа работает недостаточно быстро, сильно нагружает CPU или потребляет много памяти, процесс выявления проблем называется профилирование, а их исправление — оптимизация.

## Profiling in Go
В Go существует встроенный профайлер и утилита для визуализации результатов профилирования. Основные типы профилирования:
- Инструментальное профилирование (source or binary)
- Семплирующее профилирование

В Go используется семплирующий профайлер. Это значит, что с какой-то периодичностью профайлер прерывает работу программы, берет стек-трейс и записывает его.

###### Подробнее о классификации профилировщиков можно прочитать по этой [ссылке](https://www.delphitools.info/samplingprofiler/).

## How to get profiling
Существует три способа получения профилирования (по крайней мере, которые я знаю):
1. С помощью benchmark тестов. Мы можем запустить бенчмарк тест с флагами `-cpuprofile` и `-memprofile`.
```shell
go test -bench=. -cpuprofile cpu.out -memprofile mem.out
```
2. Импортируя библиотеку `_ "net/http/pprof"`. Эта библиотека содержит `init` функцию, которая инициализирует хэндлеры для профилирования.
```go
package somepackage

import _ "net/http/pprof"
```
3. Можно использовать функции `runtime.StartCPUProfile` или `runtime.WriteHeapProfile` в коде.

## Practicing
Рассмотрим пример простой программы на Go. Задача программы — найти захэшированный пароль, зная его максимальную длину и возможные символы.
```go
package main

import (
	"crypto/md5"
	"encoding/hex"
)

// Максимальная длина пароля
const maxPasswordLength = 5

// Возможные символы
var chars = []string{"a", "b", "c", "d", "e", "f", "g"}

func BruteForcePassword(hash string) string {
	return bruteForceRecursively(hash, "")
}

func bruteForceRecursively(hash string, passwd string) string {
	if compareHash(hash, getMD5Hash(passwd)) {
		return passwd
	}
	for i := 0; i < len(chars); i++ {
		if len(passwd) == maxPasswordLength {
			return ""
		}
		if str := bruteForceRecursively(hash, passwd+chars[i]); str != "" {
			return str
		}
	}
	return ""
}

// Хэширующая функция
func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func compareHash(a, b string) bool {
	return a == b
}
```
Функция `bruteForceRecursively` отвечает за генерацию паролей. Это рекурсивная функция, которая строит дерево всех возможных вариантов паролей и сравнивает хэш с паролем, который необходимо найти. Если хэш совпадает, функция возвращает этот пароль.

Напишем простой unit тест для проверки правильности работы программы.
###### Для написания тестов используется стандартная библиотека `testing`. Написание тестов — полезный навык, который помогает быстро проверить правильность работы программы и найти баги перед тем, как оно сломает прод (если QA профукал).
```go
package main

import (
	"fmt"
	"testing"
)

func TestBruteForcePassword(t *testing.T) {
	var table = []struct {
		input string
	}{
		{input: "a"},
		{input: "ba"},
		{input: "cf"},
	}

	for _, tab := range table {
		t.Run(fmt.Sprintf("input_%s", tab.input), func(t *testing.T) {
			if got := BruteForcePassword(getMD5Hash(tab.input)); got != tab.input {
				t.Errorf("Error want %s, got %s", tab.input, got)
			}
		})
	}
}
```
Команда для запуска теста:
```shell
go test -run=TestBruteForcePassword
```
Напишем benchmark тест для получения профиля программы.
###### Benchmark тесты помогают проверить производительность всей программы, или каждой отдельной функции.
```go
package main

import (
	"fmt"
	"testing"
)

func BenchmarkBruteForcePassword(b *testing.B) {
	var table = []struct {
		input string
	}{
		{input: "a"},
		{input: "ba"},
		{input: "cf"},
	}
	for _, tab := range table {
		b.Run(fmt.Sprintf("input_%s", tab.input), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				BruteForcePassword(getMD5Hash(tab.input))
			}
		})
	}
}
```
Запуск бенчмарк теста:
```shell
go test -bench=BenchmarkBruteForcePassword -cpuprofile cpu_1.out -memprofile mem_1.out
```
Используя `-bench`, указываем, какой тест запустить, а флаги `-cpuprofile` и `-memprofile` используются для получения профиля по CPU и памяти, сохраняя их в файлы `cpu_1.out` и `mem_1.out`.

Теперь у нас в корневом каталоге появились два новых файла `cpu_1.out` и `mem_1.out`. Мы можем воспользоваться инструментом `pprof` для визуализации этих данных.

### pprof tool
Инструмент `pprof` помогает визуализировать результаты профилирования в командной строке или в браузере.

Запуск профиля в командной строке
```shell
go tool pprof cpu_1.out
```
```shell
File: pprof.test
Type: cpu
Time: Jul 5, 2024 at 5:00pm (+05)
Duration: 24.71s, Total samples = 26.38s (106.76%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) 
```

Если ввести команду `help`, можно увидеть множество команд для работы с профилем:
```shell
(pprof) help
  Commands:
    callgrind        Outputs a graph in callgrind format
    comments         Output all profile comments
    disasm           Output assembly listings annotated with samples
    dot              Outputs a graph in DOT format
    eog              Visualize graph through eog
    evince           Visualize graph through evince
    gif              Outputs a graph image in GIF format
    gv               Visualize graph through gv
    kcachegrind      Visualize report in KCachegrind
    list             Output annotated source for functions matching regexp
    pdf              Outputs a graph in PDF format
    peek             Output callers/callees of functions matching regexp
    png              Outputs a graph image in PNG format
    proto            Outputs the profile in compressed protobuf format
    ps               Outputs a graph in PS format
    raw              Outputs a text representation of the raw profile
    svg              Outputs a graph in SVG format
    tags             Outputs all tags in the profile
    text             Outputs top entries in text form
    top              Outputs top entries in text form
    topproto         Outputs top entries in compressed protobuf format
    traces           Outputs all profile samples in text form
    tree             Outputs a text rendering of call graph
    web              Visualize graph through web browser
    weblist          Display annotated source in a web browser
    o/options        List options and their current values
    q/quit/exit/^D   Exit pprof

  Options:
    call_tree        Create a context-sensitive call tree
    compact_labels   Show minimal headers
    divide_by        Ratio to divide all samples before visualization
    drop_negative    Ignore negative differences
    edgefraction     Hide edges below <f>*total
    focus            Restricts to samples going through a node matching regexp
    hide             Skips nodes matching regexp
    ignore           Skips paths going through any nodes matching regexp
    intel_syntax     Show assembly in Intel syntax
    mean             Average sample value over first value (count)
    nodecount        Max number of nodes to show
    nodefraction     Hide nodes below <f>*total
    noinlines        Ignore inlines.
    normalize        Scales profile based on the base profile.
    output           Output filename for file-based outputs
    prune_from       Drops any functions below the matched frame.
    relative_percentages Show percentages relative to focused subgraph
    sample_index     Sample value to report (0-based index or name)
    show             Only show nodes matching regexp
    show_from        Drops functions above the highest matched frame.
    source_path      Search path for source files
    tagfocus         Restricts to samples with tags in range or matched by regexp
    taghide          Skip tags matching this regexp
    tagignore        Discard samples with tags in range or matched by regexp
    tagleaf          Adds pseudo stack frames for labels key/value pairs at the callstack leaf.
    tagroot          Adds pseudo stack frames for labels key/value pairs at the callstack root.
    tagshow          Only consider tags matching this regexp
    trim             Honor nodefraction/edgefraction/node

count even at expense of dropping details
    unit             Measurement units for the data

```

С помощью команды `top` можно посмотреть топ 10 вызовов по использованию времени процессора или памяти, в зависимости от того, какой профиль мы смотрим.
```shell
(pprof) top
Showing nodes accounting for 19930ms, 75.55% of 26380ms total
Dropped 146 nodes (cum <= 131.90ms)
Showing top 10 nodes out of 62
      flat  flat%   sum%        cum   cum%
    8130ms 30.82% 30.82%     8130ms 30.82%  crypto/md5.block
    3730ms 14.14% 44.96%     6760ms 25.63%  runtime.mallocgc
    2020ms  7.66% 52.62%     2020ms  7.66%  encoding/hex.Encode (inline)
    1150ms  4.36% 56.97%     1150ms  4.36%  runtime.nextFreeFast (inline)
    1030ms  3.90% 60.88%     1030ms  3.90%  runtime.memmove
    1010ms  3.83% 64.71%     2120ms  8.04%  runtime.concatstrings
     990ms  3.75% 68.46%     9880ms 37.45%  crypto/md5.(*digest).checkSum
     710ms  2.69% 71.15%     9340ms 35.41%  crypto/md5.(*digest).Write
     580ms  2.20% 73.35%    21950ms 83.21%  github.com/zhayt/pprof.bruteForceRecursively
     580ms  2.20% 75.55%     1950ms  7.39%  runtime.growslice
```
Каждое поле означает следующее: допустим, у нас есть функция `foo()`, которая внутри вызывает ещё четыре функции `f1(), f2(), f3(), f4()`, и ещё сама делает что-то.
```shell
func foo() {
  f1()
  f2()
  f3()
  // do something here
  f4()
}
```
`Flat` — это процессорное время или память, потраченное только на `do something here`, а `Cum` — это всё вместе взятое `f1+f2+f3+f4+do something`.

###### Подробно можно прочитать по этой [ссылке](https://stackoverflow.com/questions/32571396/pprof-and-golang-how-to-interpret-a-results).

Можно заметить, что полная работа функции `bruteForceRecursively` занимает 83.21% процессорного времени.

Давайте воспользуемся командой `list` и посмотрим изнутри функцию `BruteForcePassword`.
```shell
(pprof) list BruteForcePassword
Total: 10 samples
ROUTINE ======================== main.BruteForcePassword in /go/src/bruteForce/main.go
         0      0%      2      20% 	2
         .          .          .	}
         .          .          .
         0      0%      0       0%	func BruteForcePassword(hash string) string {
         .          .          .		return bruteForceRecursively(hash, "")
         .          .          .	}
         .          .          .
         0      0%      2      20%	func bruteForceRecursively(hash string, passwd string) string {
         .          .          .		if compareHash(hash, getMD5Hash(passwd)) {
         .          .          .			return passwd
         .          .          .		}
         .          .          .		for i := 0; i < len(chars); i++ {
         .          .          .			if len(passwd) == maxPasswordLength {
         0      0%      0       0%				return ""
         .          .          .			}
         .          .          .			if str := bruteForceRecursively(hash, passwd+chars[i]); str != "" {
         .          .          .				return str
         .          .          .			}
         .          .          .		}
         .          .          .		return ""
         .          .          .	}
(pprof)
```
Как видно в разделе `Total:`, функция `bruteForceRecursively` была вызвана 10 раз, так как она рекурсивна. Это не очень оптимально по памяти, поскольку каждый вызов функции добавляется в стек вызовов.

Если посмотреть на профиль памяти, можно увидеть, что больше всего потребляет память не функция `bruteForceRecursively`, а кодирование среза в строку `encoding/hex.EncodeToString`.
```shell
go tool pprof mem_1.out
```
```shell
File: pprof.test
Type: alloc_space
Time: Jul 7, 2024 at 7:14pm (+05)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 698.02MB, 100% of 698.02MB total
Dropped 3 nodes (cum <= 3.49MB)
      flat  flat%   sum%        cum   cum%
  520.02MB 74.50% 74.50%   520.02MB 74.50%  encoding/hex.EncodeToString (inline)
  149.50MB 21.42% 95.92%   149.50MB 21.42%  crypto/md5.(*digest).Sum (inline)
   28.50MB  4.08%   100%   571.52MB 81.88%  github.com/zhayt/pprof.bruteForceRecursively
         0     0%   100%   697.02MB 99.86%  github.com/zhayt/pprof.BenchmarkBruteForcePassword.func1
         0     0%   100%   571.52MB 81.88%  github.com/zhayt/pprof.BruteForcePassword (inline)
         0     0%   100%   669.52MB 95.92%  github.com/zhayt/pprof.getMD5Hash
         0     0%   100%   695.52MB 99.64%  testing.(*B).launch
         0     0%   100%   697.02MB 99.86%  testing.(*B).runN
```
Введя команду `web`, можно увидеть профиль в браузере.
```shell
(pprof) web
Gtk-Message: 18:45:47.077: Not loading module "atk-bridge": The functionality is provided by GTK natively. Please try to not load it.
```
Но у меня не работает модуль `gtk`, поэтому я перезапущу профилирование с флагом `http`.
```shell
go tool pprof -http=:8001 mem_1.out
```
![img_2.png](img_2.png)
Тут показывается профиль в виде аллоцированной памяти, а также можно показать количество аллоцированных объектов. Это можно сделать, нажав на `SAMPLE` в хэдере и выбрав параметр `alloc_objects`.

