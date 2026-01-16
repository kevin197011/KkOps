# Simplify License Header

## Why

The current MIT license header in source files is verbose (19 lines), making files longer and harder to scan. A simplified format (4 lines) maintains legal compliance while improving code readability.

## What Changes

Replace the full MIT license text header with a concise format:

**Current format:**
```
// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT
```

**New format:**
```
// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT
```

## Impact

- **Files affected:** All source files with the current license header (approximately 79 files across backend and frontend)
- **Scope:** Go files (`.go`), TypeScript files (`.ts`, `.tsx`)
- **Breaking changes:** None
- **Legal compliance:** Maintained (MIT license reference with URL)
