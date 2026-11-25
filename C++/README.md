# Лабораторная работа №3 — C++

Структуры данных, сериализация, модульные тесты и бенчмарки.

Исходный проект: https://github.com/VorobevAB/Lab_3 (Максим Воробьев). В этой версии исправлены вызовы методов сериализации в тестах, возвращаемые значения методов списков и скрипт формирования отчётов.

## Запуск (Ubuntu)

```bash
sudo apt install build-essential cppcheck gcovr libgtest-dev libboost-test-dev libbenchmark-dev
g++ -std=c++17 main.cpp -o lab3
./lab3
```

## Линтер, тесты и покрытие

```bash
cppcheck --enable=warning,style --std=c++17 --suppress=missingIncludeSystem struct/ main.cpp tests/ bench/
chmod +x build_reports.sh
./build_reports.sh
```

HTML-отчёты находятся в `build_gtest/`, `build_boost/`, `build_catch/`.

## Бенчмарки

```bash
g++ -std=c++17 -O2 -I. bench/bench.cpp -lbenchmark -pthread -o lab3_bench
./lab3_bench
```

Catch2 2.13.10 (`catch.hpp`) распространяется по Boost Software License 1.0, см. `LICENSE-Catch2.txt`. Файл CMakeLists.txt сохранён из исходного архива; для приведённых команд CMake не нужен.
