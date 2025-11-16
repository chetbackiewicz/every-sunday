# CSV Parsing & File Upload Research

**Status**: ✅ Complete  
**Date**: October 19, 2025  
**Decision**: PapaParse v5 + react-dropzone

---

## Recommendation Summary

**CSV Parsing**: **PapaParse v5**
- Streaming support for large files (10MB chunks)
- Header detection and type conversion
- Worker thread support prevents UI blocking

**File Upload**: **react-dropzone**
- Drag & drop interface with 10-file limit via `maxFiles` prop
- Built-in validation and error handling
- Clean React hooks API

---

## PapaParse - CSV Parsing Library

**Source**: https://www.papaparse.com/docs, https://github.com/mholt/PapaParse

### Key Features

- ✅ **Local file parsing** with streaming support for large files
- ✅ **Multiple file support** - can parse multiple files sequentially or in parallel
- ✅ **Header row detection** - auto-maps CSV columns to object keys
- ✅ **Type conversion** - `dynamicTyping` option converts numbers/booleans automatically
- ✅ **Error handling** - detailed error reporting with row numbers
- ✅ **Streaming mode** - `step` callback for row-by-row processing (prevents memory issues)
- ✅ **Worker thread support** - prevents UI blocking for large files
- ✅ **Configurable chunk size** - `Papa.LocalChunkSize = 10MB` (default)
- ✅ **Custom validation** - `validator` function for per-row validation
- ✅ **Transform functions** - modify data during parsing

### Implementation Pattern

```typescript
Papa.parse(file, {
  header: true, // Maps columns to object keys
  dynamicTyping: true, // Auto-converts types
  skipEmptyLines: true,
  complete: (results) => {
    // results.data = array of row objects
    // results.errors = array of parsing errors
    // results.meta = metadata (fields, delimiter, etc.)
  }
});
```

### Multi-file Handling

- Can parse files sequentially in loop
- Each parse is independent
- Store results in state array indexed by file
- Max 10 files requirement easily achievable

---

## react-dropzone - File Upload Component

**Source**: https://react-dropzone.js.org/, https://github.com/react-dropzone/react-dropzone

### Key Features

- ✅ **Drag & drop interface** - HTML5 compliant
- ✅ **Multiple file support** - `multiple={true}` prop
- ✅ **File type validation** - `accept` prop with MIME types
- ✅ **Max files limit** - `maxFiles` prop (perfect for 10 file limit)
- ✅ **File size limits** - `minSize` and `maxSize` props
- ✅ **Custom validation** - `validator` function
- ✅ **TypeScript support** - excellent type definitions
- ✅ **Hooks API** - `useDropzone` hook for functional components
- ✅ **Click to select** - fallback for systems without drag/drop

### Implementation Pattern

```typescript
const {getRootProps, getInputProps, acceptedFiles} = useDropzone({
  accept: {'text/csv': ['.csv']},
  maxFiles: 10,
  multiple: true,
  onDrop: (acceptedFiles) => {
    acceptedFiles.forEach(file => {
      Papa.parse(file, { /* options */ });
    });
  }
});
```

### Pros

- Handles 10 file limit with `maxFiles` prop
- Built-in file rejection/validation UI
- Clean React hooks API
- Active maintenance and wide adoption

### Cons/Limitations

- MIME type detection unreliable cross-platform (CSV can be `text/plain` or `text/csv`)
- No built-in progress indicators (must implement custom)
