# Release Guide

This project uses [GoReleaser](https://goreleaser.com/) paired with GitHub Actions to automate the creation of releases across multiple platforms (Linux, macOS, Windows) and architectures (amd64, arm64).

## How to Create a Release

Releases are triggered automatically whenever you push a Git tag to the repository.

1. **Ensure your code is ready for release:**
   Make sure all your changes are committed and pushed to the `main` branch.

2. **Create a new Git tag:**
   Tags should follow semantic versioning (e.g., `v1.0.0`, `v0.2.1`).

   ```bash
   git tag v1.0.0
   ```

3. **Push the tag to GitHub:**
   ```bash
   git push origin v1.0.0
   ```

4. **Wait for GitHub Actions to complete:**
   - Head over to the **Actions** tab in your GitHub repository.
   - You will see a workflow named `Release` running.
   - Once the workflow finishes successfully, a new release will be published on the **Releases** page of your repository.
   - The release will include auto-generated release notes (based on your commit history) and pre-compiled binaries for:
     - Linux (amd64, arm64)
     - macOS (amd64, arm64)
     - Windows (amd64, arm64)

## Configuration Details

- **`.goreleaser.yaml`**: Contains the build and archiving configuration for GoReleaser. It specifies the target operating systems and architectures, format of the archives (tar.gz for Linux/macOS, zip for Windows), and rules for changelog generation.
- **`.github/workflows/release.yml`**: The GitHub Actions workflow that installs Go, checks out the code, and runs the `goreleaser` action when a new tag is pushed.

## Testing Locally (Optional)

If you want to test the release process locally before pushing a tag, you can install the `goreleaser` CLI and run it in snapshot mode:

```bash
# This builds the binaries and archives them in the `dist/` folder without publishing.
goreleaser release --snapshot --clean
```
