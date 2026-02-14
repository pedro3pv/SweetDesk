# 🍬 SweetDesk

> **Wallpapers em 4K automático para macOS** — Baixe, upscale e use em segundos.

![Version](https://img.shields.io/badge/version-0.0.1-blue)
![License](https://img.shields.io/badge/license-MIT-green)
![macOS](https://img.shields.io/badge/macOS-11.0+-lightgrey)
![Node](https://img.shields.io/badge/Node-18+-green)

---

## 📸 O Que É SweetDesk?

**SweetDesk** é uma aplicação nativa de macOS que transforma wallpapers de **baixa/média resolução** em **imagens perfeitas em 4K (3840×2160)** usando **inteligência artificial**. O projeto utiliza o [**SweetDesk-core**](https://github.com/pedro3pv/SweetDesk-core) como engine de processamento, combinando upscaling inteligente, classificação automática e ajuste de aspect ratio.

### Principais Recursos:

✅ **Upscale Automático** — De qualquer resolução para 4K com AI (RealCUGAN + LSDIR)  
✅ **Classificação Inteligente** — Detecta automaticamente anime vs fotografia  
✅ **Content-Aware Crop** — Ajusta aspect ratio preservando conteúdo importante (Seam Carving)  
✅ **Múltiplas Fontes** — Integração com Pixabay, Unsplash, Wallhaven  
✅ **Interface macOS Nativa** — Design consistente com apps do sistema  
✅ **Batch Processing** — Processa múltiplas imagens em background  
✅ **Sem Perdas** — Upscale local no seu Mac (sem enviar para nuvem)  
✅ **Dark/Light Mode** — Segue preferências do sistema  
✅ **Aceleração por Hardware** — Suporte CoreML (Apple Silicon) e CUDA

---

## 🚀 Começando Rápido

### Pré-requisitos

- **macOS 11.0+** (Big Sur ou superior)
- **Apple Silicon (M1/M2/M3)** ou Intel x86-64
- **4GB RAM mínimo** (8GB+ recomendado para upscaling de lotes)
- **Node.js 18+** (se compilar do source)

### Instalação

#### Opção 1: Download DMG (Recomendado)

```bash
# Baixe a última versão de Releases
# https://github.com/Molasses-Co/SweetDesk/releases

# Arraste SweetDesk.app para Applications
# Abra Launchpad → SweetDesk
```

#### Opção 2: Homebrew

```bash
brew install molasses-co/sweetdesk/sweetdesk
```

#### Opção 3: Compilar do Source

```bash
# Clone o repositório
git clone https://github.com/Molasses-Co/SweetDesk.git
cd SweetDesk

# Instale dependências
npm install

# Build para macOS
npm run build:mac

# O app estará em dist/SweetDesk.app
# Mova para Applications: cp -r dist/SweetDesk.app /Applications/
```

---

## 📖 Como Usar

### Fluxo Básico

```
1️⃣  Abra SweetDesk
        ↓
2️⃣  Cole URL de wallpaper OU selecione imagem local
        ↓
3️⃣  Escolha resolução final (4K/5K/8K - padrão é 4K)
        ↓
4️⃣  App detecta: anime? foto? arte?
        ↓
5️⃣  Escolhe modelo de upscale automático
        ↓
6️⃣  Processa (30s-2min dependendo do tamanho)
        ↓
7️⃣  Preview do resultado
        ↓
8️⃣  "Set as Desktop" com 1 clique ✅
```

### Exemplos de Uso

#### Cenário 1: Foto do Unsplash → 4K

```
1. Abra SweetDesk
2. Clique "Paste from Clipboard" (após copiar URL do Unsplash)
3. Sistema detecta: "📷 Photo"
4. Aplica: LSDIR (Real-ESRGAN 4x)
5. Resultado: 3840×2160 em 4K puro
6. Clique "Set as Wallpaper" → Done!
```

#### Cenário 2: Anime de Wallhaven → 4K

```
1. Abra SweetDesk
2. Clique "Choose File" → selecione PNG anime
3. Sistema detecta: "🎨 Anime"
4. Aplica: RealCUGAN-pro (mantém linhas nítidas)
5. Resultado: 3840×2160 com anime limpo
6. Clique "Set as Wallpaper" → Done!
```

#### Cenário 3: Batch Processing (10+ imagens)

```
1. Crie pasta: ~/Pictures/ToUpscale
2. Coloque 20 imagens lá
3. Abra SweetDesk → "Batch Mode"
4. Selecione ~/Pictures/ToUpscale
5. Define output: ~/Pictures/Upscaled4K
6. Deixe rodar em background (mostra progresso)
7. Wallpapers aparecem em ~/Pictures/Upscaled4K
```

---

## 🎯 Funcionalidades Detalhadas

### 1. **Detecção Automática (Anime vs Foto)**

O app usa classificação baseada em IA para identificar o tipo de conteúdo:

- **Foto**: Rua, natureza, retrato, objeto real
  - **Modelo**: LSDIR (Real-ESRGAN 4x)
  - **Melhor para**: Preservar detalhes, texturas naturais

- **Anime**: Desenho, manga, ilustração
  - **Modelo**: RealCUGAN-pro
  - **Melhor para**: Manter linhas nítidas, cores vibrantes

- **Arte Digital**: Renderização 3D, design, abstrato
  - **Modelo**: Real-ESRGAN (UltraSharp)
  - **Melhor para**: Aumentar definição, preservar cores

### 2. **Upscaling de Resoluções**

Você escolhe o **fator de escala**:

| Tamanho Original | Escala | Resultado |
|---|---|---|
| 960×540 | 4x | **3840×2160** (4K) |
| 1920×1080 | 2x | **3840×2160** (4K) |
| 2560×1440 | 1.5x | **3840×2160** (4K) |
| 2560×1600 | 1.5x | **3840×2400** (~4K ultrawide) |
| 2560×1440 | 2x | **5120×2880** (5K) |

**Nota**: Upscale 4x + resolução arbitrária = processamento mais longo.

### 3. **Ajuste de Aspect Ratio (Content-Aware)**

Se a imagem não for 16:9 exato, o SweetDesk pode:

- **Crop (Rápido)**: Remove bordas, mantém centro
- **Seam Carving (Inteligente)**: Expande/reduz sem distorcer conteúdo importante
- **Pillar Box (Seguro)**: Adiciona fundo uniforme (menos comum)

Exemplo:
```
Original: 3840×2400 (16:10)
      ↓ (Seam Carving)
Resultado: 3840×2160 (16:9)
Conteúdo preservado, sem distorção
```

### 4. **Set as Wallpaper com 1 Clique**

Após upscale:

```
Clique "Set as Wallpaper"
    ↓
SweetDesk salva em:
~/Library/Application Support/SweetDesk/Wallpapers/
    ↓
Chama System Preferences via AppleScript
    ↓
Desktop & Screen Saver → Seleciona a imagem
    ↓
✅ Wallpaper aplicado em todos os desktops
```

---

## 🔧 Arquitetura Técnica

```
┌─────────────────────────────────────────────┐
│           SweetDesk (Frontend)              │
│  React + TypeScript + Electron/Tauri       │
│  (UI, preview, file picker)                 │
└──────────────┬──────────────────────────────┘
               │
               ↓
┌─────────────────────────────────────────────┐
│       Backend (Next.js API Route)           │
│  Node.js + TypeScript                       │
│  (Orquestração, classificação)              │
└──────┬─────────────────────────────┬────────┘
       │                             │
       ↓                             ↓
┌──────────────────┐    ┌─────────────────────────┐
│  Classificador   │    │   SweetDesk-core (Go)   │
│  (IA/ML)         │    │   Engine de Processing  │
│  anime vs foto   │    │   - Upscaling (ONNX)    │
└──────────────────┘    │   - RealCUGAN / LSDIR   │
                        │   - Seam Carving        │
                        │   - Tiling Automático   │
                        └─────────┬───────────────┘
                                  │
                                  ↓
                     ┌─────────────────────────┐
                     │  ONNX Runtime           │
                     │  (Bibliotecas embarcadas)│
                     │  + CoreML/CUDA          │
                     └─────────────────────────┘
```

### Stack Técnico

| Layer | Tecnologia | Propósito |
|-------|-----------|----------|
| **Frontend** | React 18 + TypeScript | UI interativa |
| **Runtime** | Electron ou Tauri | App nativa macOS |
| **Backend** | Next.js 14 API Routes | Orquestração |
| **Core Engine** | [SweetDesk-core](https://github.com/pedro3pv/SweetDesk-core) (Go) | Processamento de imagens |
| **Upscaling** | ONNX Runtime + RealCUGAN/LSDIR | IA local, sem nuvem |
| **Content-Aware** | Seam Carving | Ajuste inteligente |
| **Aceleração** | CoreML (macOS) / CUDA | Hardware acceleration |
| **Storage** | Sistema de arquivos local | Processamento offline |
| **OS Integration** | AppleScript + Foundation | Set as Wallpaper |

### Como Funciona o SweetDesk-core

O [**SweetDesk-core**](https://github.com/pedro3pv/SweetDesk-core) é o motor de processamento escrito em Go que:

1. **Embarca bibliotecas ONNX Runtime** no executável durante build
2. **Classifica automaticamente** imagens (anime vs foto) usando modelos ML
3. **Aplica upscaling** com modelos apropriados:
   - **RealCUGAN**: Para anime/ilustrações
   - **LSDIR**: Para fotografias realísticas
4. **Processa em tiles** para imagens grandes (evita sobrecarga de memória)
5. **Aplica seam carving** quando necessário ajustar aspect ratio
6. **Acelera via hardware** usando CoreML (Apple Silicon) ou CUDA

**Vantagens da Integração:**
- ✅ **Sem downloads em runtime** — bibliotecas embarcadas
- ✅ **Cross-platform** — suporta macOS (Intel + ARM), Linux e Windows
- ✅ **Performance nativa** — escrito em Go com ONNX otimizado
- ✅ **API pública** — reutilizável em outros projetos

---

## 📦 Estrutura do Projeto

```
SweetDesk/
├── src/
│   ├── components/          # React components (UI)
│   │   ├── UploadZone.tsx
│   │   ├── Preview.tsx
│   │   ├── SettingsPanel.tsx
│   │   └── BatchMode.tsx
│   ├── pages/
│   │   ├── api/
│   │   │   ├── classify.ts        # Detecta anime vs foto
│   │   │   ├── upscale.ts         # Chama SweetDesk-core
│   │   │   ├── crop.ts            # Ajusta aspect ratio
│   │   │   └── set-wallpaper.ts   # AppleScript bridge
│   │   └── index.tsx              # Home page
│   ├── lib/
│   │   ├── core-integration.ts    # Wrapper para SweetDesk-core
│   │   ├── classifier.ts          # Classificador de imagens
│   │   └── macos-integration.ts   # AppleScript, System Prefs
│   └── types/
│       └── index.ts               # TypeScript types
├── public/
│   └── icons/                     # App icons (icns)
├── scripts/
│   ├── build-mac.sh               # Build para .dmg
│   ├── download-core.sh           # Download SweetDesk-core binary
│   └── setup-env.sh               # Setup inicial
├── next.config.js
├── tsconfig.json
├── package.json
└── README.md (este arquivo)
```

---

## ⚙️ Configuração Avançada

### Modelos de Upscaling

O SweetDesk-core gerencia automaticamente os modelos ONNX:

```bash
# Modelos são embarcados no core ou baixados na primeira execução:
- RealCUGAN-pro (anime)
- LSDIR (fotografias)

# Localização:
~/.cache/sweetdesk/models/
```

### Customizar Threshold de Classificação

Arquivo: `src/lib/classifier.ts`

```typescript
const CLASSIFICATION_THRESHOLD = 0.7; // 0-1, default 0.7
// Valores menores = mais sensível em detectar anime
```

### Ativar Debug Mode

```bash
# Terminal
export SWEETDESK_DEBUG=1
open /Applications/SweetDesk.app

# Mostra logs completos no console
```

### Configurar Aceleração por Hardware

O SweetDesk-core automaticamente detecta e usa:
- **CoreML** em Apple Silicon (M1/M2/M3)
- **CUDA** em GPUs NVIDIA (se disponível)
- **CPU** como fallback

Para forçar CPU-only:
```bash
export SWEETDESK_FORCE_CPU=1
```

---

## 🖥️ Sistema de Requisitos

### Mínimo

- macOS 11 Big Sur
- 4GB RAM
- 2GB espaço em disco (modelos + cache)
- Processador com suporte a ONNX Runtime

### Recomendado

- macOS 13+ Ventura/Sonoma
- Apple Silicon (M1/M2/M3+) ou Intel i7+
- 16GB RAM
- 5GB SSD (processamento mais rápido)

### Performance por Hardware

| Hardware | Upscale 4x (1080p→4K) | 2x (2K→4K) |
|---|---|---|
| M1 Pro | ~45s | ~20s |
| M2 Max | ~35s | ~15s |
| M3 Max | ~30s | ~12s |
| Intel i9 (10th Gen) | ~60s | ~28s |

---

## 🔒 Segurança & Privacidade

✅ **Sem servidor externo** — Upscaling ocorre 100% localmente no seu Mac  
✅ **Sem coleta de dados** — Nenhuma telemetria enviada  
✅ **Open source** — Código auditável no GitHub  
✅ **Modelos compactados** — ONNX Runtime otimizado

**Armazenamento**:
- Imagens temporárias: `~/Library/Application Support/SweetDesk/temp/` (limpas após uso)
- Wallpapers finais: `~/Library/Application Support/SweetDesk/wallpapers/` (sua propriedade)
- Modelos IA: `~/.cache/sweetdesk/models/` (somente leitura)
- Cache do core: `./cache/onnxruntime/` (bibliotecas extraídas)

---

## 📥 Integrações

### Import de URLs

Suporta direto:

```
✅ Unsplash.com (copie a URL de download)
✅ Wallhaven.cc (download direto)
✅ Pexels.com (download direto)
✅ Pixabay.com (download direto)
✅ Qualquer URL JPEG/PNG público
```

**Exemplo**:
```
1. Unsplash → Imagem → Clique "Download" → Copy Link
2. SweetDesk → Paste URL
3. App baixa e processa
```

### Set as Wallpaper

Integra com **System Preferences** via AppleScript:

```applescript
tell application "System Preferences"
    activate
    set current pane to pane id "com.apple.preference.desktopscreeneffect"
    # E define a imagem no sistema
end tell
```

### Batch Export

Exporta para:

```
✅ Pasta local (~/Pictures/Upscaled4K/)
✅ iCloud Drive (~/Library/Mobile Documents/)
```

---

## 🐛 Troubleshooting

### Problema: "App não abre no macOS Sonoma"

**Solução**:
```bash
# Remova quarentena (se fez download manual)
xattr -d com.apple.quarantine /Applications/SweetDesk.app

# Ou: System Preferences → Security & Privacy → Allow SweetDesk
```

### Problema: Upscaling muito lento

**Checklist**:
1. Verificar RAM disponível
2. Fechar apps pesados (Chrome, Photoshop, etc.)
3. Apple Silicon usa CoreML automaticamente
4. Tentar resolução menor (2K em vez de 4K)

**Se lento demais**:
```bash
# Verificar se CoreML está ativo (Apple Silicon)
# Ou usar resolução intermediária
```

### Problema: "Biblioteca embarcada não encontrada"

**Solução**:
```bash
# Reinstalar o SweetDesk-core
npm run download-core

# Ou baixar manualmente:
# https://github.com/pedro3pv/SweetDesk-core/releases
```

### Problema: "Set as Wallpaper" não funciona

**Solução**:
```bash
# Verificar permissões AppleScript
System Preferences → Privacy & Security → Automation
    ↓
Procure "SweetDesk" → Marque todas as permissões
```

---

## 🤝 Contribuindo

Adoramos contribuições! Aqui está como ajudar:

### Setup Desenvolvimento

```bash
# Clone e instale
git clone https://github.com/Molasses-Co/SweetDesk.git
cd SweetDesk
npm install

# Rode em dev mode
npm run dev

# O app abre em Electron/Tauri com hot reload
```

### Estrutura de PR

1. **Fork** o repo
2. **Branch**: `git checkout -b feature/minha-feature`
3. **Commit**: `git commit -m "Add: descrição clara"`
4. **Push**: `git push origin feature/minha-feature`
5. **PR**: Abra no GitHub com descrição

### Áreas Procurando Help

- [ ] Suporte a **Windows / Linux** (via SweetDesk-core)
- [ ] Integração **Apple Shortcuts** (automação)
- [ ] **Performance optimization** para Intel chips
- [ ] Documentação em **outras linguagens** (pt-BR, es, ja, etc.)
- [ ] **Testes unitários** (Jest + React Testing Library)

---

## 📄 Licença

**SweetDesk** é distribuído sob a **MIT License**.

### Componentes e Dependências

| Componente | Licença | Comercial OK? | Notas |
|---|---|---|---|
| **SweetDesk-core** | MIT | ✅ Sim | Engine de processamento |
| **ONNX Runtime** | MIT | ✅ Sim | Inferência de modelos |
| **RealCUGAN** | MIT-like | ✅ Sim | Upscaling anime |
| **LSDIR (Real-ESRGAN)** | BSD | ✅ Sim | Upscaling fotográfico |

**IMPORTANTE**: Se você modificar ou redistribuir este software, **mantenha a licença MIT intacta** e inclua aviso de copyright.

---

## 🔗 Links Relacionados

- **[SweetDesk-core](https://github.com/pedro3pv/SweetDesk-core)** — Engine de processamento (Go)
- **[ONNX Runtime](https://github.com/microsoft/onnxruntime)** — Runtime de ML
- **[RealCUGAN](https://github.com/bilibili/ailab)** — Upscaling de anime
- **[LSDIR](https://github.com/cszn/LSDIR)** — Upscaling realístico

---

## 📞 Suporte

- **Issues**: [GitHub Issues](https://github.com/Molasses-Co/SweetDesk/issues)
- **Discussions**: [GitHub Discussions](https://github.com/Molasses-Co/SweetDesk/discussions)

---

## 🗺️ Roadmap

### v1.0 (Atual)
- [x] Integração com SweetDesk-core
- [x] Upscale 4K para foto/anime
- [x] Classificação automática
- [x] Set as Wallpaper integrado
- [x] Batch processing básico
- [x] Dark/Light mode

### v1.1 (Planejado)
- [ ] Suporte a **5K/8K explícito**
- [ ] **Color correction** pós-upscale
- [ ] **Smart crop** com detecção de faces
- [ ] **Multiple display setup** (diferentes resoluções por monitor)
- [ ] **Scheduled wallpaper rotation** (trocar a cada hora/dia)

### v2.0 (Futuro)
- [ ] **Windows & Linux** support (via SweetDesk-core)
- [ ] **AI-powered wallpaper generation** (Text-to-Image)
- [ ] **Wallpaper marketplace integrado** (Unsplash + Wallhaven APIs)
- [ ] **Local AI model training** (seu próprio estilo)
- [ ] **Cloud sync** (sincronizar wallpapers entre Macs)

---

## 🎨 Créditos

Desenvolvido por **[Molasses Co.](https://molasses.co)** com ❤️ para a comunidade macOS.

### Agradecimentos Especiais

- **Pedro Augusto ([@pedro3pv](https://github.com/pedro3pv))** — Desenvolvedor do SweetDesk-core
- **RealCUGAN Team** — Upscaling de anime
- **LSDIR/Real-ESRGAN** — Upscaling fotográfico
- **Microsoft ONNX Runtime** — Runtime de ML
- **Tauri/Electron** — Framework nativo
- **Community** — Feedback e PRs

---

## 📖 Documentação Adicional

- [**Quick Start Guide**](./docs/QUICKSTART.md)
- [**Advanced Configuration**](./docs/ADVANCED.md)
- [**Architecture Overview**](./docs/ARCHITECTURE.md)
- [**Contributing Guide**](./CONTRIBUTING.md)
- [**Changelog**](./CHANGELOG.md)

---

## 📊 Status do Projeto

![Build Status](https://img.shields.io/badge/build-passing-brightgreen)
![macOS](https://img.shields.io/badge/platform-macOS-lightgrey)
![License](https://img.shields.io/badge/license-MIT-green)

---

## 🌟 Dê uma Star!

Se SweetDesk foi útil, considere dar uma ⭐ no GitHub!

```
https://github.com/Molasses-Co/SweetDesk
```

---

**SweetDesk** — *Wallpapers Lindos em 4K, Automático* 🍬✨

**Última atualização**: Fevereiro 2026  
**Versão**: 0.0.1  
**Mantenedor**: [@molassesco](https://github.com/Molasses-Co)  
**Core Engine**: [SweetDesk-core](https://github.com/pedro3pv/SweetDesk-core) by [@pedro3pv](https://github.com/pedro3pv)
