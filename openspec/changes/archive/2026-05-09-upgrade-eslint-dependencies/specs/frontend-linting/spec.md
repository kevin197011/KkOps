# Frontend Linting Specification

## MODIFIED Requirements

### Requirement: ESLint Configuration Uses Flat Config Format
The frontend project SHALL use ESLint 9 with flat config format (`eslint.config.js`) instead of the legacy `.eslintrc.json` format. ESLint 9 requires flat config format, provides better performance, and resolves deprecated dependency warnings.

#### Scenario: ESLint 9 Flat Config is Used
- **Given** the frontend project has ESLint installed
- **When** the developer runs `npm run lint`
- **Then** ESLint should use `eslint.config.js` (flat config format)
- **And** the old `.eslintrc.json` file should not exist
- **And** all linting rules should work as before

#### Scenario: Deprecated Dependencies Are Resolved
- **Given** npm dependencies are installed
- **When** the developer runs `npm install`
- **Then** no deprecated warnings should appear for:
  - `rimraf@3.0.2`
  - `glob@7.2.3`
  - `inflight@1.0.6`
  - `@humanwhocodes/config-array@0.13.0`
  - `@humanwhocodes/object-schema@2.0.3`

#### Scenario: ESLint Plugins Are Compatible with ESLint 9
- **Given** ESLint 9 is installed
- **When** TypeScript and React plugins are loaded
- **Then** all plugins should be compatible with ESLint 9:
  - `@typescript-eslint/eslint-plugin@8.x` or `typescript-eslint@8.x`
  - `@typescript-eslint/parser@8.x` (if using separate packages)
  - `eslint-plugin-react-hooks@^5.x`
  - `eslint-plugin-react-refresh@latest`

#### Scenario: Existing Linting Rules Are Preserved
- **Given** the project had specific ESLint rules configured
- **When** migrating to ESLint 9 flat config
- **Then** all existing rules should be preserved:
  - TypeScript recommended rules
  - React Hooks recommended rules
  - React Refresh rules (only-export-components)
  - Custom rules (no-unused-vars, no-explicit-any)
  - Custom settings (React version detection)

#### Scenario: IDE Integration Works with Flat Config
- **Given** VSCode with ESLint extension is used
- **When** opening TypeScript/React files
- **Then** ESLint should provide real-time linting
- **And** errors and warnings should be displayed correctly
- **And** the flat config format should be recognized

## ADDED Requirements

### Requirement: ESLint 9 Dependency Management
The project SHALL use ESLint 9.x with compatible plugin versions that resolve all deprecated dependency warnings. All npm deprecated warnings for rimraf, glob, inflight, and @humanwhocodes packages SHALL be eliminated.

#### Scenario: ESLint 9 is Installed
- **Given** the frontend dependencies are up to date
- **When** checking `package.json`
- **Then** `eslint` should be version `^9.x`
- **And** all ESLint-related plugins should be compatible versions

#### Scenario: No Deprecated Warnings During Installation
- **Given** the frontend project dependencies
- **When** running `npm install`
- **Then** no deprecated package warnings should appear
- **And** all dependencies should use maintained, up-to-date versions
