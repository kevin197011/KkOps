# Design: ESLint 9 Migration

## Overview

This document outlines the migration from ESLint 8 to ESLint 9, including the transition from `.eslintrc.json` to flat config format.

## ESLint 9 Flat Config Format

ESLint 9 introduces the flat config format (`eslint.config.js`) as the default configuration method. The new format:
- Uses JavaScript instead of JSON
- Is more explicit and flexible
- Provides better performance
- Supports better IDE integration

### Configuration Mapping

**Old Format (`.eslintrc.json`)**:
```json
{
  "root": true,
  "env": { "browser": true, "es2020": true, "node": true },
  "extends": [
    "eslint:recommended",
    "plugin:@typescript-eslint/recommended",
    "plugin:react-hooks/recommended"
  ],
  "parser": "@typescript-eslint/parser",
  "parserOptions": {
    "ecmaVersion": "latest",
    "sourceType": "module",
    "ecmaFeatures": { "jsx": true },
    "project": "./tsconfig.json"
  },
  "plugins": ["react-refresh", "@typescript-eslint"],
  "rules": { ... },
  "settings": { "react": { "version": "detect" } }
}
```

**New Format (`eslint.config.js`)**:
```javascript
import js from '@eslint/js';
import tseslint from 'typescript-eslint';
import reactHooks from 'eslint-plugin-react-hooks';
import reactRefresh from 'eslint-plugin-react-refresh';

export default tseslint.config(
  { ignores: ['dist', 'node_modules'] },
  {
    files: ['**/*.{ts,tsx}'],
    languageOptions: {
      ecmaVersion: 'latest',
      sourceType: 'module',
      parserOptions: {
        ecmaFeatures: { jsx: true },
        project: './tsconfig.json',
      },
    },
    plugins: {
      'react-hooks': reactHooks,
      'react-refresh': reactRefresh,
    },
    rules: {
      ...js.configs.recommended.rules,
      ...tseslint.configs.recommended.rules,
      ...reactHooks.configs.recommended.rules,
      'react-refresh/only-export-components': [
        'warn',
        { allowConstantExport: true },
      ],
      '@typescript-eslint/no-unused-vars': [
        'error',
        { argsIgnorePattern: '^_', varsIgnorePattern: '^_' },
      ],
      '@typescript-eslint/no-explicit-any': 'warn',
    },
    settings: {
      react: { version: 'detect' },
    },
  },
);
```

## Plugin Compatibility

### TypeScript ESLint
- **ESLint 8**: `@typescript-eslint/eslint-plugin@7.x` and `@typescript-eslint/parser@7.x`
- **ESLint 9**: `typescript-eslint@8.x` (bundled package) or separate `@typescript-eslint/eslint-plugin@8.x` and `@typescript-eslint/parser@8.x`

### React Hooks
- `eslint-plugin-react-hooks@^5.x` supports ESLint 9

### React Refresh
- `eslint-plugin-react-refresh@latest` should work with ESLint 9

## Migration Strategy

1. **Create flat config alongside old config** - Keep both during transition
2. **Test thoroughly** - Ensure all rules work correctly
3. **Remove old config** - Only after verification

## Potential Issues

1. **Rule naming changes** - Some rules may have been renamed or removed
2. **Plugin initialization** - Flat config requires explicit plugin imports
3. **IDE support** - VSCode ESLint extension needs to support flat config (v3.0.0+)

## References

- [ESLint 9 Migration Guide](https://eslint.org/docs/latest/use/migrate-to-9.0.0)
- [Flat Config Documentation](https://eslint.org/docs/latest/use/configure/configuration-files-new)
- [TypeScript ESLint Flat Config](https://typescript-eslint.io/users/configs)
