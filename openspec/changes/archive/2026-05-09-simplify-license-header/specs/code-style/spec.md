## MODIFIED Requirements

### Requirement: License Header Format
All source files SHALL use a simplified MIT license header format.

#### Scenario: New file header format
- **WHEN** a source file contains a license header
- **THEN** the header SHALL use the format:
  ```
  // Copyright (c) 2025 kk
  //
  // This software is released under the MIT License.
  // https://opensource.org/licenses/MIT
  ```
- **AND** the header SHALL replace the verbose MIT license text (19 lines) with the simplified format (4 lines)
- **AND** the copyright year SHALL remain as 2025
- **AND** the license reference URL SHALL be included
