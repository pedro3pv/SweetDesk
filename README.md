# 🍬 SweetDesk

> **Aplicação desktop para upscaling inteligente de wallpapers** — De qualquer resolução para 4K usando IA.

![Version](https://img.shields.io/badge/version-0.0.1-blue)
![License](https://img.shields.io/badge/license-MIT-green)
![Go](https://img.shields.io/badge/Go-1.25.7-00ADD8)
![Next.js](https://img.shields.io/badge/Next.js-16.1.6-black)

---

## 📸 O Que É SweetDesk?

**SweetDesk** é uma aplicação desktop multiplataforma que transforma wallpapers de **baixa/média resolução** em **imagens 4K (3840×2160)** usando **inteligência artificial**. O projeto combina [**SweetDesk-core**](https://github.com/Molasses-Co/SweetDesk-core) (motor de processamento em Go) com uma interface moderna em Next.js + React.

### Principais Recursos Implementados:

✅ **Upscale Automático** — De qualquer resolução para 4K/customizado com IA  
✅ **Classificação Inteligente** — Detecta automaticamente anime vs fotografia (via SweetDesk-core)  
✅ **Múltiplas Fontes** — Busca e download integrado com Pixabay  
✅ **Batch Processing** — Processa múltiplas imagens em lote com progresso em tempo real  
✅ **Interface Moderna** — Next.js 16 + React 19 + TailwindCSS  
✅ **Processamento Local** — Todo upscale roda no seu computador (sem enviar para nuvem)  
✅ **Cross-Platform** — macOS, Windows e Linux (via Wails)

---

## 🏗️ Arquitetura

```
┌─────────────────────────────────────────────┐
│           Frontend (Next.js 16)             │
│  React 19 + TypeScript + TailwindCSS        │
│  - Busca de imagens (Pixabay)              │
│  - Upload e preview                         │
│  - Processamento batch                      │
└──────────────┬──────────────────────────────┘
               │ (Wails Bridge)
               ↓
┌─────────────────────────────────────────────┐
│       Backend Go (Wails v2.11.0)            │
│  - ImageProcessor                           │
│  - CoreBridge (integração)                  │
│  - API providers (Pixabay)                  │
└──────┬──────────────────────────────────────┘
       │
       ↓
┌─────────────────────────────────────────────┐
│   SweetDesk-core v0.0.4 (Go Module)         │
│   github.com/Molasses-Co/SweetDesk-core        │
│   - Classificação automática (anime/foto)   │
│   - Upscaling (RealCUGAN / LSDIR)           │
│   - Seam carving (aspect ratio)             │
│   - ONNX Runtime embarcado                  │
└─────────────────────────────────────────────┘
```

### Stack Técnico

| Camada | Tecnologia | Versão | Propósito |
|--------|-----------|--------|-----------|
| **Frontend** | Next.js + React | 16.1.6 / 19.2.4 | Interface do usuário |
| **UI Framework** | TailwindCSS | 4.1.18 | Estilização |
| **Runtime** | Wails | 2.11.0 | App nativa multiplataforma |
| **Backend** | Go | 1.25.7 | Lógica de negócios |
| **Core Engine** | SweetDesk-core | 0.0.4 | Processamento de imagens |
| **AI Runtime** | ONNX Runtime | 1.25.0 | Inferência de modelos |

---

## 🚀 Começando

### Pré-requisitos

- **Go 1.25.7+**
- **Node.js 18+** ou **Bun**
- **Wails v2.11.0** ([Instalar Wails](https://wails.io/docs/gettingstarted/installation))
- **4GB RAM mínimo** (8GB+ recomendado para batch processing)

### Instalação para Desenvolvimento

```bash
# 1. Clone o repositório
git clone https://github.com/Molasses-Co/SweetDesk.git
cd SweetDesk

# 2. Instale dependências do frontend
cd frontend
npm install  # ou: bun install / pnpm install
cd ..

# 3. Baixe dependências do Go
go mod download

# 4. Rode em modo desenvolvimento
wails dev

# O app abrirá com hot reload ativado
```

### Build para Produção

```bash
# Build para o sistema operacional atual
wails build

# Binário estará em: build/bin/
```

### Configuração da API Pixabay (Opcional)

Para usar a busca de imagens integrada:

```bash
# Crie um arquivo .env na raiz do projeto
echo "PIXABAY_API_KEY=sua_chave_aqui" > .env

# Obtenha sua chave em: https://pixabay.com/api/docs/
```

---

## 📖 Como Usar

### Fluxo Básico

```
1️⃣  Abra SweetDesk
        ↓
2️⃣  Escolha uma opção:
    • Fazer upload de imagem local
    • Buscar no Pixabay
        ↓
3️⃣  Selecione resolução de saída (ex: 3840x2160)
        ↓
4️⃣  Clique "Process"
        ↓
5️⃣  Aguarde processamento (30s-2min)
        ↓
6️⃣  Salve a imagem upscalada ✅
```

### Processamento em Lote

```bash
1. Abra a aba "Batch"
2. Adicione múltiplas imagens
3. Configure resolução de saída
4. Selecione pasta de destino
5. Clique "Start Batch Processing"
6. Acompanhe progresso em tempo real
```

---

## 📦 Estrutura do Projeto

```
SweetDesk/
├── frontend/                    # Next.js 16 + React 19
│   ├── src/
│   │   ├── app/                # App Router do Next.js
│   │   ├── components/         # Componentes React
│   │   │   ├── ImageUpload.tsx
│   │   │   ├── ImagePreview.tsx
│   │   │   ├── SearchPanel.tsx
│   │   │   ├── DownloadList.tsx
│   │   │   ├── ProcessingView.tsx
│   │   │   ├── ProcessingPanel.tsx
│   │   │   ├── FolderSelect.tsx
│   │   │   ├── ImageDetail.tsx
│   │   │   └── CompleteView.tsx
│   │   └── lib/                # Utilitários
│   ├── wailsjs/                # Bindings gerados pelo Wails
│   ├── package.json
│   ├── next.config.ts
│   └── tailwind.config.ts
│
├── internal/                    # Backend Go
│   └── services/
│       ├── core_bridge.go      # Integração com SweetDesk-core
│       ├── image_processor.go  # Processamento de imagens
│       └── api_provider.go     # Pixabay API
│
├── app.go                       # Aplicação principal Wails
├── main.go                      # Entry point
├── go.mod                       # Dependências Go
├── wails.json                   # Configuração Wails
└── README.md
```

---

## 🎯 Funcionalidades Detalhadas

### 1. Upscaling Inteligente

O SweetDesk usa o [SweetDesk-core](https://github.com/Molasses-Co/SweetDesk-core) que aplica automaticamente o modelo mais adequado:

- **Fotos/Realismo**: LSDIR (Real-ESRGAN 4x)
- **Anime/Ilustrações**: RealCUGAN-pro

**Opções de Processamento:**

```go
type ProcessingOptions struct {
    TargetWidth     int     // Largura final (ex: 3840)
    TargetHeight    int     // Altura final (ex: 2160)
    ScaleFactor     float64 // Fator de escala (ex: 4.0)
    MaxResolution   int     // Limite máximo (padrão: 16384)
    KeepAspectRatio bool    // Manter proporção
}
```

### 2. Busca de Imagens (Pixabay)

Interface integrada para buscar wallpapers:

- Busca por palavras-chave
- Filtros de resolução mínima
- Preview antes do download
- Download direto para processamento

### 3. Processamento em Lote

Sistema robusto para processar múltiplas imagens:

- Adicione quantas imagens quiser
- Progresso individual por imagem
- Estados: `pending` → `processing` → `done` / `error`
- Notificações em tempo real via Wails Events

**Exemplo de uso interno:**

```go
items := []BatchItem{
    {ID: "1", Base64Data: "...", Dimension: "3840x2160"},
    {ID: "2", DownloadURL: "https://...", Dimension: "3840x2160"},
}
ProcessBatch(items, "/caminho/destino")
```

### 4. Seam Carving (via SweetDesk-core)

Ajuste inteligente de aspect ratio que preserva conteúdo importante ao invés de apenas cortar/distorcer a imagem.

---

## ⚙️ Configuração Avançada

### Variáveis de Ambiente

```bash
# .env (raiz do projeto)
PIXABAY_API_KEY=sua_chave_aqui
```

### Limites de Processamento

Configurados em `app.go`:

```go
const (
    defaultMaxResolution = 16384  // Máximo 16K
    maxPixels = 16384 * 16384     // ~268 megapixels
)
```

### Debug Mode

```bash
# Ativar logs detalhados
export SWEETDESK_DEBUG=1
wails dev
```

---

## 🔧 Desenvolvimento

### Rodar Frontend Isolado (Next.js)

```bash
cd frontend
npm run dev
# Abre em http://localhost:3000
```

### Testar Backend

```bash
# Rodar testes
go test ./internal/services/...

# Testar CoreBridge
go test -v ./internal/services -run TestCoreBridge
```

### Atualizar SweetDesk-core

```bash
# Atualizar para versão específica
go get github.com/Molasses-Co/SweetDesk-core@v0.0.5

# Ou usar versão local para desenvolvimento
# Edite go.mod e descomente:
# replace github.com/Molasses-Co/SweetDesk-core => ../SweetDesk-core

go mod tidy
```

---

## 🐛 Troubleshooting

### "CoreBridge not initialized"

**Causa**: SweetDesk-core não foi carregado corretamente.

**Solução**:
```bash
# Verificar se o módulo está baixado
go mod download
go mod verify

# Rebuild o projeto
wails build
```

### "Pixabay API key not configured"

**Solução**:
```bash
# Criar arquivo .env na raiz
echo "PIXABAY_API_KEY=sua_chave" > .env
```

### Performance Lenta

**Dicas**:
1. Usar resoluções intermediárias (2K ao invés de 8K)
2. Processar lotes menores (5-10 imagens por vez)
3. Fechar outros apps pesados
4. SweetDesk-core usa automaticamente aceleração de hardware quando disponível

---

## 🤝 Contribuindo

Contribuições são bem-vindas! 

### Como Contribuir

```bash
# 1. Fork o repositório
# 2. Crie uma branch
git checkout -b feature/minha-feature

# 3. Commit suas mudanças
git commit -m "feat: adiciona nova funcionalidade"

# 4. Push para o branch
git push origin feature/minha-feature

# 5. Abra um Pull Request
```

### Áreas que Precisam de Ajuda

- [ ] Integração com mais providers (Unsplash, Wallhaven)
- [ ] Preset de resoluções comuns (1080p, 1440p, 4K, etc.)
- [ ] Suporte a arrastar e soltar
- [ ] Testes unitários do frontend
- [ ] Documentação em outras línguas

---

## 📄 Licença

**MIT License** — Veja [LICENSE](./LICENSE) para detalhes.

### Dependências

| Componente | Licença | Link |
|-----------|---------|------|
| SweetDesk-core | MIT | [GitHub](https://github.com/Molasses-Co/SweetDesk-core) |
| Wails | MIT | [wails.io](https://wails.io) |
| Next.js | MIT | [nextjs.org](https://nextjs.org) |
| ONNX Runtime | MIT | [onnxruntime.ai](https://onnxruntime.ai) |
| RealCUGAN | BSD-3 | [GitHub](https://github.com/bilibili/ailab) |

---

## 🔗 Links Relacionados

- **[SweetDesk-core](https://github.com/Molasses-Co/SweetDesk-core)** — Motor de processamento em Go
- **[Wails](https://wails.io)** — Framework para apps desktop com Go
- **[Next.js](https://nextjs.org)** — Framework React
- **[Pixabay API](https://pixabay.com/api/docs/)** — API de imagens grátis

---

## 📞 Suporte

- **Issues**: [GitHub Issues](https://github.com/Molasses-Co/SweetDesk/issues)
- **Discussions**: [GitHub Discussions](https://github.com/Molasses-Co/SweetDesk/discussions)

---

## 🗺️ Status do Projeto

### ✅ Implementado (v0.0.1)

- [x] Aplicação Wails com Next.js + React
- [x] Integração com SweetDesk-core v0.0.4
- [x] Upload de imagens locais
- [x] Busca e download via Pixabay
- [x] Upscaling individual (resolução customizável)
- [x] Processamento em lote
- [x] Progresso em tempo real
- [x] Preview de imagens

### 🚧 Em Desenvolvimento

- [ ] Mais providers (Unsplash, Wallhaven)
- [ ] Presets de resolução (1-click para 4K, 5K, etc.)
- [ ] Histórico de processamentos
- [ ] Configurações persistentes

### 💡 Planejado para Futuras Versões

- [ ] Modo "Set as Wallpaper" automático (macOS/Windows/Linux)
- [ ] Suporte a múltiplos monitores
- [ ] Rotação agendada de wallpapers
- [ ] Color correction pós-upscale
- [ ] Integração com serviços de nuvem

---

## 🎨 Créditos

Desenvolvido por **[@pedro3pv](https://github.com/pedro3pv)** para **[Molasses Co.](https://github.com/Molasses-Co)**

### Agradecimentos

- **SweetDesk-core** — Motor de processamento
- **RealCUGAN Team** — Upscaling de anime
- **Real-ESRGAN/LSDIR** — Upscaling fotográfico
- **Microsoft ONNX Runtime** — Inferência de ML
- **Wails Team** — Framework desktop
- **Vercel** — Next.js e React

---

**SweetDesk** — *Wallpapers em Alta Resolução com IA* 🍬✨

**Versão**: 0.0.1  
**Última atualização**: Fevereiro 2026  
**Mantenedor**: [@pedro3pv](https://github.com/pedro3pv)  
**Core Engine**: [SweetDesk-core v0.0.4](https://github.com/Molasses-Co/SweetDesk-core)
