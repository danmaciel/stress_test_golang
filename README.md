# Teste de stress com Golang


## Como usar?

Entre na raiz do projeto e execute o seguinte comando para criar e subir o banco de dados Mysql e o servidor do rabbitMq:

```
go run main.go exec --url "urlDoSite" --requests "quantidadeDeRequests" --concurrency "quantideRequisicoesSimultanea"
```

Sendo os parametros:

```
url: url completa do site que sera feito o teste
requests: quantidade total de requisicoes que serao executadas
concurrency: quantidade de requisicoes executadas simultaneamente
```
