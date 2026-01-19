# WebSSH ZMODEM File Transfer Specification

## ADDED Requirements

### Requirement: ZMODEM File Upload Support

The WebSSH terminal **SHALL** support file upload via the `rz` command using the ZMODEM protocol.

#### Scenario: User uploads file using rz command
1. User runs `rz` command in the WebSSH terminal
2. System detects ZMODEM upload sequence (`**B`)
3. **Expected**: File picker dialog automatically opens in the browser
4. User selects a file
5. **Expected**: File is uploaded in chunks via WebSocket binary messages
6. **Expected**: Upload progress is displayed (filename, progress bar, bytes transferred)
7. **Expected**: Upon completion, file is saved on the remote server
8. **Expected**: User can cancel the upload during transfer

### Requirement: ZMODEM File Download Support

The WebSSH terminal **SHALL** support file download via the `sz` command using the ZMODEM protocol.

#### Scenario: User downloads file using sz command
1. User runs `sz <filename>` command in the WebSSH terminal
2. System detects ZMODEM download sequence (`**G`)
3. **Expected**: Download progress modal appears
4. **Expected**: File data is received via WebSocket binary messages
5. **Expected**: Download progress is displayed (filename, progress bar, bytes transferred)
6. **Expected**: Upon completion, file is automatically downloaded by the browser
7. **Expected**: User can cancel the download during transfer

### Requirement: ZMODEM Sequence Detection

The system **SHALL** detect ZMODEM protocol initiation sequences in SSH output streams.

#### Scenario: ZMODEM sequence detection in stdout
1. Remote server outputs ZMODEM sequence (`**B` for upload, `**G` for download)
2. **Expected**: Backend detects the sequence before sending to frontend
3. **Expected**: Sequence is intercepted and not displayed in terminal
4. **Expected**: `zmodem_start` message is sent to frontend with direction and mode information

#### Scenario: ZMODEM sequence in stderr
1. Remote server outputs ZMODEM sequence in stderr stream
2. **Expected**: Backend detects the sequence in stderr as well
3. **Expected**: Same handling as stdout (intercept, send message to frontend)

### Requirement: Binary Data Transfer

The system **SHALL** handle binary file data separately from text terminal output.

#### Scenario: Binary data during file upload
1. File upload is in progress
2. Frontend sends file chunks as WebSocket BinaryMessage
3. **Expected**: Backend receives binary messages and forwards to SSH stdin
4. **Expected**: Binary data is not processed as text

#### Scenario: Binary data during file download
1. File download is in progress
2. Backend receives binary data from SSH stdout
3. **Expected**: Backend sends binary data as WebSocket BinaryMessage to frontend
4. **Expected**: Frontend receives and buffers binary data chunks
5. **Expected**: Binary data is assembled into complete file

### Requirement: Transfer Progress Display

The system **SHALL** display file transfer progress to users.

#### Scenario: Upload progress display
1. File upload starts
2. **Expected**: Modal dialog appears showing:
   - Upload direction icon
   - Filename
   - Progress bar (percentage)
   - Bytes transferred / Total bytes
   - Cancel button
3. **Expected**: Progress updates in real-time as chunks are sent

#### Scenario: Download progress display
1. File download starts
2. **Expected**: Modal dialog appears with same information as upload
3. **Expected**: Progress updates in real-time as chunks are received

### Requirement: Transfer Cancellation

The system **SHALL** allow users to cancel file transfers.

#### Scenario: Cancel file upload
1. User clicks cancel button during upload
2. **Expected**: Upload is stopped
3. **Expected**: WebSocket sends `zmodem_cancel` message to backend
4. **Expected**: Backend stops forwarding data to SSH
5. **Expected**: Progress modal is closed
6. **Expected**: Transfer state is reset

### Requirement: Protocol Sequence Filtering

The system **SHALL** prevent ZMODEM protocol sequences from appearing in the terminal.

#### Scenario: ZMODEM sequence not displayed
1. Remote server outputs `**B0100000023be50` (ZMODEM upload sequence)
2. **Expected**: Sequence is intercepted by backend
3. **Expected**: Sequence does not appear in terminal output
4. **Expected**: Only text output before/after sequence is displayed

### Requirement: Transfer Error Handling

The system **SHALL** handle and display errors during file transfer.

#### Scenario: Upload error handling
1. Upload fails (e.g., network error, file system error on server)
2. **Expected**: Error message is displayed to user
3. **Expected**: Transfer state is reset
4. **Expected**: Progress modal is closed

#### Scenario: Download error handling
1. Download fails (e.g., file not found, permission denied)
2. **Expected**: Error message is displayed to user
3. **Expected**: Transfer state is reset
4. **Expected**: Progress modal is closed
