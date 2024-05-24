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

Para gerar uma imagem Docker que tem a aplicacao pronta pra uso, entre no diretorio raiz da aplicacao e rode o seguinte comando para gerar a 
imagem

```
docker build -t nome_da_imagem .
```

Para executar

```
docker run -it nome_da_imagem --url=http://www.site.com --requests=50 --concurrency=2
```

ou

```
docker run -it nome_da_imagem -u http://www.site.com -r 50 -c 2
```