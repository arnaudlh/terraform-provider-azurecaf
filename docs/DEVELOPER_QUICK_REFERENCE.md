# Developer Quick Reference - CI/CD Pipeline

This is a quick reference guide for developers working with the modernized CI/CD pipeline.

## 🚀 Quick Start

### Essential Commands
```bash
# Set up your development environment
make dev_setup

# Before committing changes
make format lint unittest

# Before pushing to PR
make ci_local

# Run all quality checks
make qa

# Clean up build artifacts
make clean
```

## 🔄 Workflow Triggers

| Event | Workflows Triggered | Purpose |
|-------|-------------------|---------|
| PR to main | `go.yml`, `e2e.yml`, `security.yml` | Full validation |
| Push to main | `go.yml`, `performance.yml` | Integration validation |
| Tagged release (`v*`) | `go.yml`, `release.yml` | Release process |
| Daily schedule | `e2e.yml`, `security.yml`, `dependencies.yml` | Maintenance checks |
| Manual dispatch | Any workflow | On-demand testing |

## 🧪 Testing Strategy

### Local Testing Hierarchy
```bash
# Level 1: Basic checks (< 1 min)
make format lint

# Level 2: Unit tests (< 3 min)
make unittest

# Level 3: Integration tests (< 10 min)
make ci_local

# Level 4: Full validation (< 20 min)
make ci_full
```

### CI Test Matrix
- **Unit Tests**: Fast, isolated component testing
- **Integration Tests**: Provider functionality validation
- **E2E Tests**: End-to-end scenarios with Terraform
- **Performance Tests**: Benchmark and profiling

## 🔒 Security Checks

### Automated Scans
- **Gosec**: Go security analyzer
- **Nancy**: Dependency vulnerability scanner
- **govulncheck**: Official Go vulnerability checker
- **go-licenses**: License compliance verification

### Security Workflow
1. Runs on every PR and daily
2. Results appear in GitHub Security tab
3. SARIF reports provide detailed findings
4. Failed scans block releases

## 📈 Performance Monitoring

### What's Monitored
- **Benchmark Results**: Performance regression detection
- **Memory Usage**: Memory leak detection
- **CPU Profiling**: Performance bottleneck identification
- **Test Timing**: Execution speed tracking

### Accessing Results
- Check workflow artifacts for detailed reports
- Performance data retained for 90 days
- Compare against baseline metrics

## 🔧 Common Workflows

### Making a Code Change
```bash
# 1. Create feature branch
git checkout -b feature/my-feature

# 2. Make changes and test locally
make format lint unittest

# 3. Run integration tests
make ci_local

# 4. Commit and push
git add .
git commit -m "feat: my new feature"
git push origin feature/my-feature

# 5. Create PR (triggers CI automatically)
```

### Fixing a Test Failure
```bash
# 1. Reproduce locally
make unittest  # or specific test target

# 2. Check coverage impact
make test_coverage

# 3. Run full validation
make ci_local

# 4. For E2E issues
make test_e2e_quick
```

### Performance Investigation
```bash
# 1. Run benchmarks
make benchmark

# 2. Generate memory profile
make profile_mem

# 3. Generate CPU profile  
make profile_cpu

# 4. Analyze profiles
make profile_analyze

# 5. Check results in *.txt files
```

## 📦 Release Process

### Manual Release
```bash
# 1. Ensure main branch is ready
git checkout main
git pull origin main

# 2. Validate locally
make qa_full

# 3. Create and push tag
git tag v1.2.3
git push origin v1.2.3

# 4. Release workflow runs automatically
```

### Pre-release
Use GitHub UI for workflow dispatch with pre-release option.

## 🛠 Troubleshooting

### Common Issues

#### "Tests pass locally but fail in CI"
```bash
# Use exact CI environment
make ci_local

# Check Go version consistency
go version

# Verify dependencies
go mod verify
```

#### "Cache issues slowing CI"
- Check if `go.sum` changed significantly
- Cache keys are based on `go.sum` hash
- New dependencies invalidate cache

#### "Security scan failures"
1. Check GitHub Security tab
2. Review SARIF reports in artifacts
3. Update dependencies if needed:
   ```bash
   make dependency_check
   ```

#### "Performance regression detected"
1. Download benchmark artifacts
2. Compare with previous results
3. Profile the issue:
   ```bash
   make profile_cpu profile_mem
   make profile_analyze
   ```

## 📋 Make Target Reference

### Build & Test
| Target | Purpose | Duration |
|--------|---------|----------|
| `build` | Build provider | < 1m |
| `unittest` | Unit tests only | < 3m |
| `test_coverage` | With coverage | < 5m |
| `test_integration` | Integration tests | < 10m |

### Quality & Security
| Target | Purpose | Duration |
|--------|---------|----------|
| `lint` | Code linting | < 1m |
| `format` | Code formatting | < 30s |
| `security_scan` | Security analysis | < 3m |
| `dependency_check` | Dependency audit | < 2m |

### Performance
| Target | Purpose | Duration |
|--------|---------|----------|
| `benchmark` | Performance tests | < 5m |
| `profile_mem` | Memory profiling | < 10m |
| `profile_cpu` | CPU profiling | < 10m |
| `profile_analyze` | Analyze profiles | < 1m |

### Development
| Target | Purpose | Duration |
|--------|---------|----------|
| `dev_setup` | Setup environment | < 3m |
| `dev_build` | Local install | < 2m |
| `watch_tests` | Continuous testing | Continuous |
| `clean` | Clean artifacts | < 30s |

### Composite
| Target | Purpose | Duration |
|--------|---------|----------|
| `qa` | Quality checks | < 10m |
| `qa_full` | Full QA suite | < 20m |
| `ci_local` | Local CI | < 15m |
| `ci_full` | Full CI suite | < 30m |

## 🔗 Useful Links

- **GitHub Actions**: [Repository Actions Tab](../../actions)
- **Security Findings**: [Repository Security Tab](../../security)
- **Release History**: [Repository Releases](../../releases)
- **Performance Data**: Check workflow artifacts
- **Documentation**: [Full CI/CD Documentation](CI_CD_PIPELINE.md)

## 💡 Tips

### Development Efficiency
- Use `make watch_tests` during development
- Run `make ci_local` before pushing
- Check Security tab regularly for issues
- Review performance artifacts periodically

### CI Optimization
- Keep PRs focused to minimize CI time
- Use draft PRs to skip CI until ready
- Monitor cache hit rates in workflow logs
- Report CI issues early

### Best Practices
- Always format code before committing
- Write tests for new functionality
- Update documentation with changes
- Monitor performance impact of changes

---

*For detailed information, see the [full CI/CD documentation](CI_CD_PIPELINE.md).*
