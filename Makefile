.PHONY: build test clean release tag patch minor major help

# Default version
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

# Build the binary
build:
	go build -o sjq

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -f sjq

# Show current version
version:
	@echo "Current version: $(VERSION)"

# Create a new release
# Usage: make release VERSION=1.0.1
release: test
	@if [ -z "$(VERSION)" ] || [ "$(VERSION)" = "dev" ]; then \
		echo "Error: VERSION is required. Usage: make release VERSION=1.0.1"; \
		exit 1; \
	fi
	@echo "Releasing version $(VERSION)..."
	@# Update version in version.go
	@sed -i.bak 's/const Version = ".*"/const Version = "$(VERSION)"/' version.go && rm version.go.bak
	@# Commit version change
	@git add version.go
	@git commit -m "Release version $(VERSION)"
	@# Create and push tag
	@git tag -a v$(VERSION) -m "Release version $(VERSION)"
	@echo "Version $(VERSION) released!"
	@echo "Run 'git push && git push --tags' to push the release"

# Create a tag without updating version.go
# Usage: make tag VERSION=1.0.1
tag:
	@if [ -z "$(VERSION)" ] || [ "$(VERSION)" = "dev" ]; then \
		echo "Error: VERSION is required. Usage: make tag VERSION=1.0.1"; \
		exit 1; \
	fi
	@git tag -a v$(VERSION) -m "Release version $(VERSION)"
	@echo "Tag v$(VERSION) created!"
	@echo "Run 'git push --tags' to push the tag"

# Semantic versioning shortcuts
# Usage: make patch/minor/major
patch:
	@./release.sh patch

minor:
	@./release.sh minor

major:
	@./release.sh major

# Show help
help:
	@echo "Available targets:"
	@echo "  build    - Build the sjq binary"
	@echo "  test     - Run tests"
	@echo "  clean    - Remove build artifacts"
	@echo "  version  - Show current version"
	@echo "  release  - Create a new release (updates version.go, commits, and tags)"
	@echo "             Usage: make release VERSION=1.0.1"
	@echo "  patch    - Create a patch release (e.g., 1.0.0 -> 1.0.1)"
	@echo "  minor    - Create a minor release (e.g., 1.0.0 -> 1.1.0)"
	@echo "  major    - Create a major release (e.g., 1.0.0 -> 2.0.0)"
	@echo "  tag      - Create a tag without updating version.go"
	@echo "             Usage: make tag VERSION=1.0.1"
	@echo "  help     - Show this help message"

# Default target
.DEFAULT_GOAL := build