# Upgrade ESLint and Resolve Deprecated Dependencies

## Why

The frontend project currently uses ESLint 8.57.1, which includes several deprecated dependencies:
- `rimraf@3.0.2` (via flat-cache → file-entry-cache → eslint)
- `glob@7.2.3` (via rimraf)
- `inflight@1.0.6` (via glob)
- `@humanwhocodes/config-array@0.13.0` (deprecated, use `@eslint/config-array`)
- `@humanwhocodes/object-schema@2.0.3` (deprecated, use `@eslint/object-schema`)

These deprecated packages:
1. Generate npm warnings during installation
2. May have security vulnerabilities
3. Are no longer maintained
4. Will be automatically replaced by upgrading ESLint to v9

ESLint 9 introduces:
- Flat config format (modern configuration)
- Better performance
- Native support for `@eslint/config-array` and `@eslint/object-schema`
- Updated dependencies that resolve all deprecated warnings

## What Changes

1. **Upgrade ESLint** from 8.57.1 to 9.x (latest stable)
2. **Migrate ESLint configuration** from `.eslintrc.json` to flat config format (`eslint.config.js`)
3. **Upgrade related plugins** to versions compatible with ESLint 9:
   - `@typescript-eslint/eslint-plugin` → 8.x
   - `@typescript-eslint/parser` → 8.x
   - `eslint-plugin-react-hooks` → latest (verify compatibility)
   - `eslint-plugin-react-refresh` → latest (verify compatibility)
4. **Update lint script** if needed for flat config format
5. **Verify all linting rules** work correctly after migration

## Impact

### Files Affected
- `frontend/package.json` - Update dependencies
- `frontend/.eslintrc.json` - Remove (replaced by flat config)
- `frontend/eslint.config.js` - New flat config file
- `frontend/package-lock.json` - Auto-updated by npm

### Breaking Changes
- ESLint 9 requires flat config format (breaking change from `.eslintrc.*`)
- Some plugin APIs may have changed
- Need to verify all existing rules are compatible

### Testing Required
- Run `npm run lint` to verify no regressions
- Check that existing code passes linting
- Verify IDE integration (VSCode ESLint extension) works correctly
