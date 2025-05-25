# Desemprego Zero - Backend

Backend para o sistema de conexão entre membros da igreja e oportunidades de trabalho.

## Requisitos

- Go 1.21 ou superior
- PostgreSQL
- Docker (opcional)

## Configuração

1. Clone o repositório
2. Copie o arquivo `.env.example` para `.env` e configure as variáveis de ambiente
3. Instale as dependências:

```bash
go mod download
```

## Executando Localmente

### Com Go

1. Certifique-se de que o PostgreSQL está rodando
2. Execute o servidor:

```bash
go run main.go
```

### Com Docker

1. Construa e inicie os containers:

```bash
docker-compose up --build
```

2. Para executar em background:

```bash
docker-compose up -d
```

3. Para parar os containers:

```bash
docker-compose down
```

## Migrações e Setup Inicial

### Localmente

1. Execute as migrações:

```bash
go run main.go
```

2. Crie o administrador inicial:

```bash
go run scripts/create_admin.go
```

### Com Docker

1. Execute as migrações e crie o admin:

```bash
docker-compose exec app go run scripts/create_admin.go
```

## Deploy

### Railway

1. Crie uma conta no [Railway](https://railway.app)
2. Instale o CLI do Railway:

```bash
npm i -g @railway/cli
```

3. Login no Railway:

```bash
railway login
```

4. Inicialize o projeto:

```bash
railway init
```

5. Adicione as variáveis de ambiente no dashboard do Railway:

```
DB_HOST=containers-us-west-XX.railway.app
DB_PORT=XXXX
DB_USER=postgres
DB_PASSWORD=seu_password
DB_NAME=railway
JWT_SECRET=sua_chave_secreta
```

6. Faça o deploy:

```bash
railway up
```

### Render

1. Crie uma conta no [Render](https://render.com)
2. Crie um novo Web Service
3. Conecte seu repositório GitHub
4. Configure as variáveis de ambiente
5. Use os seguintes comandos de build e start:
   - Build Command: `go build -o main .`
   - Start Command: `./main`

## Estrutura do Projeto

```
.
├── internal/
│   ├── models/      # Modelos do banco de dados
│   ├── middleware/  # Middlewares (auth, error handling)
│   ├── handlers/    # Handlers das rotas
│   └── routes/      # Configuração das rotas
├── scripts/         # Scripts de utilidade
├── main.go         # Ponto de entrada da aplicação
├── go.mod          # Dependências
├── Dockerfile      # Configuração do container
├── docker-compose.yml # Configuração dos serviços
└── .env           # Variáveis de ambiente
```

## Rotas da API

### Rotas Públicas

- GET /jobs - Lista todas as vagas
- GET /jobs/:id - Visualiza detalhes de uma vaga
- POST /candidates - Cadastro de candidatos

### Rotas Protegidas (Admin)

- POST /admin/login - Login de administrador
- POST /jobs - Criar vaga
- PUT /jobs/:id - Editar vaga
- DELETE /jobs/:id - Excluir vaga
- GET /candidates - Listar candidatos

## Health Check

- GET /health - Verifica o status da aplicação e conexão com o banco
