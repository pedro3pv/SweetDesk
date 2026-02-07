# Contributing to SweetDesk

Thank you for your interest in contributing to SweetDesk! 🍬✨

This document provides guidelines for contributing to the project.

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment for everyone.

## How Can I Contribute?

### Reporting Bugs 🐛

Before creating bug reports, please check existing issues to avoid duplicates.

**Good Bug Report Includes:**
- Clear, descriptive title
- Steps to reproduce the issue
- Expected vs actual behavior
- Screenshots (if applicable)
- Environment details:
  - OS and version
  - SweetDesk version
  - Hardware specs

**Template:**
```markdown
**Describe the bug**
A clear description of what the bug is.

**To Reproduce**
Steps to reproduce the behavior:
1. Go to '...'
2. Click on '....'
3. See error

**Expected behavior**
What you expected to happen.

**Screenshots**
If applicable, add screenshots.

**Environment:**
 - OS: [e.g. macOS 14.0]
 - Version: [e.g. 0.1.0]
 - Hardware: [e.g. M1 Mac, 16GB RAM]
```

### Suggesting Features 💡

Feature requests are welcome! Please:
1. Check if the feature is already requested
2. Clearly describe the feature and its benefits
3. Provide examples or mockups if possible

**Template:**
```markdown
**Is your feature request related to a problem?**
A clear description of the problem.

**Describe the solution you'd like**
A clear description of what you want to happen.

**Describe alternatives you've considered**
Other solutions you've thought about.

**Additional context**
Any other context or screenshots.
```

### Contributing Code 💻

#### Getting Started

1. **Fork the repository**
   ```bash
   # Click "Fork" on GitHub
   git clone https://github.com/YOUR_USERNAME/SweetDesk.git
   cd SweetDesk
   ```

2. **Set up development environment**
   ```bash
   # Install dependencies
   cd frontend && npm install && cd ..
   go mod download
   pip install -r python/requirements.txt
   
   # Create .env file
   cp .env.example .env
   # Add your API keys
   ```

3. **Create a branch**
   ```bash
   git checkout -b feature/your-feature-name
   # or
   git checkout -b fix/your-bug-fix
   ```

#### Making Changes

1. **Write clean code**
   - Follow existing code style
   - Add comments for complex logic
   - Keep functions small and focused
   - Use meaningful variable names

2. **Test your changes**
   ```bash
   # Run in dev mode
   wails dev
   
   # Test thoroughly
   # - Try different scenarios
   # - Test error cases
   # - Check performance
   ```

3. **Commit your changes**
   ```bash
   # Stage your changes
   git add .
   
   # Commit with a clear message
   git commit -m "Add: Feature description"
   ```

   **Commit Message Format:**
   ```
   Type: Brief description
   
   Longer explanation if needed.
   Explain why, not what.
   ```

   **Types:**
   - `Add:` New feature
   - `Fix:` Bug fix
   - `Update:` Modify existing feature
   - `Remove:` Delete code/feature
   - `Refactor:` Code restructuring
   - `Docs:` Documentation changes
   - `Style:` Code style changes
   - `Test:` Add or update tests

#### Submitting Pull Requests

1. **Push to your fork**
   ```bash
   git push origin feature/your-feature-name
   ```

2. **Create pull request**
   - Go to your fork on GitHub
   - Click "New Pull Request"
   - Fill in the template
   - Link related issues

**PR Checklist:**
- [ ] Code follows project style
- [ ] Changes are tested
- [ ] No console errors
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
- [ ] Commit messages are clear
- [ ] PR description is detailed

**PR Template:**
```markdown
## Description
Brief description of changes.

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
How was this tested?

## Screenshots
If applicable, add screenshots.

## Checklist
- [ ] Code compiles without errors
- [ ] Tests pass
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
```

### Improving Documentation 📚

Documentation improvements are always welcome!

**Areas to improve:**
- Fix typos or unclear instructions
- Add examples or tutorials
- Improve API documentation
- Translate to other languages
- Add diagrams or screenshots

**How to contribute:**
1. Edit the markdown files
2. Follow the existing style
3. Submit a PR with your changes

## Code Style Guidelines

### Go

```go
// Use meaningful names
func ProcessImage(data []byte) ([]byte, error) {
    // Add comments for complex logic
    if err != nil {
        return nil, fmt.Errorf("processing failed: %w", err)
    }
    return result, nil
}

// Group imports
import (
    "context"
    "fmt"
    
    "SweetDesk/internal/services"
)
```

### TypeScript/React

```typescript
// Use functional components
export default function MyComponent({ prop }: Props) {
    // Use hooks for state
    const [state, setState] = useState<string>('');
    
    // Use callbacks for handlers
    const handleClick = useCallback(() => {
        // Handler logic
    }, []);
    
    return (
        <div className="container">
            {/* JSX */}
        </div>
    );
}
```

### Python

```python
"""
Module docstring.
"""

def process_image(image_path: str) -> dict:
    """
    Function docstring.
    
    Args:
        image_path: Path to the image file.
    
    Returns:
        Dictionary with results.
    """
    # Implementation
    return result
```

## Development Workflow

### Branch Strategy

- `main` - Stable, production-ready code
- `develop` - Development branch
- `feature/*` - New features
- `fix/*` - Bug fixes
- `docs/*` - Documentation updates

### Release Process

1. Merge approved PRs to `develop`
2. Test thoroughly
3. Update version in `package.json` and `wails.json`
4. Update CHANGELOG.md
5. Merge `develop` to `main`
6. Tag release: `git tag v0.1.0`
7. Build and distribute

## Areas Needing Help

### High Priority
- [ ] Download and integrate actual AI binaries
- [ ] Implement "Set as Wallpaper" for each OS
- [ ] Add comprehensive error handling
- [ ] Write unit tests
- [ ] Add integration tests

### Medium Priority
- [ ] Add Unsplash API integration
- [ ] Implement batch processing
- [ ] Add user preferences system
- [ ] Improve UI/UX
- [ ] Add more image filters

### Low Priority
- [ ] Translate documentation (PT-BR, ES, JA)
- [ ] Add more example images
- [ ] Create video tutorials
- [ ] Improve performance
- [ ] Add telemetry (optional)

## Questions?

- 💬 GitHub Discussions: https://github.com/Molasses-Co/SweetDesk/discussions
- 🐛 GitHub Issues: https://github.com/Molasses-Co/SweetDesk/issues

## Recognition

Contributors will be:
- Listed in CONTRIBUTORS.md
- Mentioned in release notes
- Thanked in the README

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to SweetDesk! 🍬✨
