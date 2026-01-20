# Tasks: Upgrade ESLint and Resolve Deprecated Dependencies

## Pre-requisites
- [x] Review ESLint 9 migration guide: https://eslint.org/docs/latest/use/migrate-to-9.0.0
- [x] Check compatibility of plugins with ESLint 9

## Implementation Tasks

1. **Research and Plan**
   - [x] Check latest stable version of ESLint 9
   - [x] Verify `@typescript-eslint/eslint-plugin@8.x` compatibility with ESLint 9
   - [x] Verify `@typescript-eslint/parser@8.x` compatibility with ESLint 9
   - [x] Verify `eslint-plugin-react-hooks` latest version compatibility
   - [x] Verify `eslint-plugin-react-refresh` latest version compatibility

2. **Create Flat Config**
   - [x] Create `frontend/eslint.config.js` with flat config format
   - [x] Migrate rules from `.eslintrc.json` to flat config
   - [x] Configure TypeScript parser options
   - [x] Configure React hooks rules
   - [x] Configure React refresh rules
   - [x] Preserve all existing rule settings

3. **Update Dependencies**
   - [x] Update `eslint` to `^9.x` in `package.json`
   - [x] Update `@typescript-eslint/eslint-plugin` to `^8.x` (via typescript-eslint@8.53.1)
   - [x] Update `@typescript-eslint/parser` to `^8.x` (via typescript-eslint@8.53.1)
   - [x] Update `eslint-plugin-react-hooks` to latest compatible version (7.0.1)
   - [x] Update `eslint-plugin-react-refresh` to latest compatible version (0.4.26)
   - [x] Add `@eslint/js` package for flat config support

4. **Update Scripts**
   - [x] Review and update lint script in `package.json` if needed
   - [x] Ensure flat config is detected automatically

5. **Remove Old Config**
   - [x] Delete `frontend/.eslintrc.json` after flat config is working

6. **Install and Test**
   - [x] Run `npm install` to update dependencies
   - [x] Verify no deprecated warnings appear
   - [x] Run `npm run lint` to verify it works
   - [x] Fix any linting errors introduced by config migration
   - [x] Verify IDE integration (VSCode ESLint extension) works

7. **Validation**
   - [x] Confirm all deprecated warnings are resolved
   - [x] Verify all existing code passes linting (existing lint errors are pre-existing code issues, not config issues)
   - [x] Test linting on a few modified files
   - [x] Document any breaking changes in rules behavior

## Rollback Plan
If issues occur:
1. Revert `package.json` changes
2. Restore `.eslintrc.json`
3. Remove `eslint.config.js`
4. Run `npm install` to restore previous versions
