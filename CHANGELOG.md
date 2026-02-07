# Changelog

All notable changes to SweetDesk will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Complete application architecture with Wails + Next.js
- Pixabay API integration for image search
- Image upload via drag & drop or file picker
- AI-powered image classification (anime vs photo detection)
- AI upscaling using Real-ESRGAN and RealCUGAN models
- Content-aware image resizing with seam carving
- Fast center crop for aspect ratio adjustment
- In-memory image processing (no disk I/O until save)
- Cross-platform support (macOS, Windows, Linux)
- Modern UI with dark/light mode support
- Real-time processing progress feedback
- Image preview with before/after comparison
- Export to 4K, 5K, or 8K resolutions
- Multi-API provider architecture (ready for Unsplash, Pexels, etc.)

### Backend (Go/Wails)
- Image processor service for encoding/decoding
- API provider service with Pixabay implementation
- Python bridge for AI model integration
- Upscaler service for managing AI binaries
- Complete Wails bindings for frontend integration

### Frontend (Next.js + React)
- Search panel with Pixabay integration
- Image upload component with drag & drop
- Processing panel with resolution and options
- Image preview with side-by-side comparison
- Responsive UI with Tailwind CSS
- TypeScript for type safety

### Python Integration
- Image classification script (heuristic-based)
- Seam carving implementation
- Ready for DeepGHS/imgutils integration
- NumPy and Pillow for image processing

### Build System
- Cross-platform build scripts
- AI binary download automation
- Platform-specific model distribution
- Development and production build modes

### Documentation
- Comprehensive README
- Quick start guide (QUICKSTART.md)
- Development guide (DEVELOPMENT.md)
- API documentation (API.md)
- Build instructions (BUILD.md)
- Environment configuration (.env.example)

## [0.0.1] - 2026-02-07

### Initial Release
- Project structure setup
- Wails + Next.js integration
- Basic application scaffolding

---

## Future Roadmap

### v0.1.0 (Next)
- [ ] Download actual AI binaries (Real-ESRGAN, RealCUGAN)
- [ ] Implement "Set as Wallpaper" functionality
- [ ] Add batch processing support
- [ ] Improve error handling and user feedback
- [ ] Add processing queue for multiple images
- [ ] Implement image caching

### v0.2.0
- [ ] Add Unsplash API integration
- [ ] Add Pexels API integration
- [ ] Implement Supabase storage integration
- [ ] Add user preferences/settings
- [ ] Multi-language support (PT-BR, ES, JA)
- [ ] Add image filters and adjustments

### v0.3.0
- [ ] GPU acceleration for upscaling
- [ ] Advanced classification using DeepGHS/imgutils
- [ ] Face detection for smart cropping
- [ ] Color correction and enhancement
- [ ] HDR processing

### v1.0.0
- [ ] Stable release with all core features
- [ ] Complete test coverage
- [ ] Production-ready AI binaries
- [ ] Installer packages for all platforms
- [ ] macOS App Store release
- [ ] Windows Store release

### v2.0.0 (Future)
- [ ] AI-powered wallpaper generation (Text-to-Image)
- [ ] Wallpaper marketplace integration
- [ ] Cloud sync between devices
- [ ] Mobile app (iOS/Android)
- [ ] Scheduled wallpaper rotation
- [ ] Custom AI model training

---

**Legend:**
- `Added` - New features
- `Changed` - Changes in existing functionality
- `Deprecated` - Soon-to-be removed features
- `Removed` - Removed features
- `Fixed` - Bug fixes
- `Security` - Security improvements
