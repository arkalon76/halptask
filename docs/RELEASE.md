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
    - The release will include auto-generated release notes (based on your commit history) and artifacts:
      - **Standalone Uncompressed Binaries**: Direct executables (`halptask_Linux_x86_64`, `halptask_Darwin_arm64`, `halptask_Windows_x86_64.exe`, etc.) for instant use without unpacking.
      - **Linux Packages**: `.deb` (Debian/Ubuntu) and `.rpm` (Fedora/RHEL) packages for native package manager installation.
      - **Compressed Archives**: Standard `.tar.gz` and `.zip` archives for traditional distribution.
      - **Installer Script**: Automated one-line installation via `curl -fsSL ... | bash`.

## Configuration Details

- **`.goreleaser.yaml`**: Configured with GoReleaser v2 specifications. Generates both raw standalone binary assets (`formats: [binary]`) and compressed archives (`formats: [tar.gz, zip]`), as well as Debian (`.deb`) and RedHat (`.rpm`) packages (`nfpms`).
- **`scripts/install.sh`**: Portable shell installer script that auto-detects OS/architecture, retrieves the raw binary asset from GitHub Releases, and places it into `/usr/local/bin` (or `~/.local/bin`).
- **`.github/workflows/release.yml`**: GitHub Actions workflow that executes GoReleaser on new version tag pushes.

## Testing Locally (Optional)

If you want to test the release process locally before pushing a tag, you can install the `goreleaser` CLI and run it in snapshot mode:

```bash
# This builds the binaries and archives them in the `dist/` folder without publishing.
goreleaser release --snapshot --clean
```
