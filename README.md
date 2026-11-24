# DevOps Lab 7

[![CI](https://github.com/AleksandrPleskachev/devops_lab7/actions/workflows/ci.yml/badge.svg)](https://github.com/AleksandrPleskachev/devops_lab7/actions/workflows/ci.yml)

CI/CD для сервиса «Курсы» из лабораторной №6. Вариант 7 — ежедневный
запуск матричной проверки в 03:15 UTC. Параллельно проверяются Go
1.21.13 и 1.22.12, артефакты coverage сохраняются отдельно.

    go build ./...
    go test -covermode=atomic -coverpkg=./pkg/... -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html

Dockerfile собирает сервис и запускает его от непривилегированного
пользователя. HTTP доступен на 8080, gRPC — на 9090.
