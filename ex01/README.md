# 42 School — Cybersecurity Piscine — Arachnida — Exercise 01.

## Description

Recursively downloads images from a website. Built in Go using a Standard
Project Layout combined with a pragmatic Hexagonal / Clean Architecture:
the domain layer defines entities and port interfaces with zero external
dependencies, the application layer implements use cases against those
ports, and the adapter layer provides the concrete infrastructure
(HTTP, HTML scanning, filesystem) injected by `cmd/spider/main.go`.
